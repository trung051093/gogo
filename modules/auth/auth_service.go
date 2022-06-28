package auth

import (
	"context"
	"errors"
	"fmt"
	"gogo/common"
	"gogo/components/appctx"
	cacheprovider "gogo/components/cache"
	"gogo/components/hasher"
	"gogo/components/mailer"
	authmodel "gogo/modules/auth/model"
	authprovider "gogo/modules/auth_provider"
	googleauthprovider "gogo/modules/auth_provider/google"
	jwtauthprovider "gogo/modules/auth_provider/jwt"
	authprovidermodel "gogo/modules/auth_provider/model"
	usermodel "gogo/modules/user/model"
	"strconv"

	"gogo/modules/user"

	"github.com/google/uuid"
)

type authService struct {
	config              *appctx.Config
	jwtProvider         jwtauthprovider.JWTProvider
	googleProvider      googleauthprovider.GoogleAuthProvider
	userService         user.UserService
	authProviderService authprovider.AuthProviderService
	hashService         hasher.HashService
	cacheService        cacheprovider.CacheService
	mailService         mailer.MailService
}

func NewAuthService(
	config *appctx.Config,
	jwtProvider jwtauthprovider.JWTProvider,
	googleProvider googleauthprovider.GoogleAuthProvider,
	userService user.UserService,
	authProviderService authprovider.AuthProviderService,
	hashService hasher.HashService,
	cacheService cacheprovider.CacheService,
	mailService mailer.MailService,
) *authService {
	return &authService{
		config:              config,
		jwtProvider:         jwtProvider,
		googleProvider:      googleProvider,
		userService:         userService,
		authProviderService: authProviderService,
		hashService:         hashService,
		cacheService:        cacheService,
		mailService:         mailService,
	}
}

func NewAuthServiceFromContext(appCtx appctx.AppContext) *authService {
	appConfig := appCtx.GetConfig()
	userRepo := user.NewUserRepository(appCtx.GetMainDBConnection())
	authProviderRepository := authprovider.NewAuthProviderRepository(appCtx.GetMainDBConnection())
	esService := appCtx.GetESService()
	cacheService := appCtx.GetCacheService()
	mailService := appCtx.GetMailService()
	userService := user.NewUserService(userRepo, esService)
	authProviderService := authprovider.NewAuthProviderService(authProviderRepository)
	hashService := hasher.NewHashService()
	jwtProvider := jwtauthprovider.NewJWTProvider(appConfig.JWT.Secret)
	googleProvider := googleauthprovider.NewGoogleAuthProvider((*googleauthprovider.GoogleAuthConfig)(&appConfig.GoogleOauth2))

	authService := NewAuthService(
		appConfig,
		jwtProvider,
		googleProvider,
		userService,
		authProviderService,
		hashService,
		cacheService,
		mailService,
	)
	return authService
}

func (s *authService) JWTUserGenerate(ctx context.Context, user *usermodel.User) (*authprovidermodel.TokenProvider, error) {
	token, err := s.jwtProvider.Generate(authprovidermodel.TokenPayload{
		UserId: user.Id,
		Email:  user.Email,
		Role:   user.Role,
	}, uint(s.config.JWT.ExpireDays))

	if err != nil || token == nil {
		return nil, err
	}

	// store to cache
	key := s.getCacheKey(ctx, AuthenticationToken, strconv.Itoa(user.Id))
	s.setSession(ctx, key, token.Token, s.config.JWT.ExpireDays)

	return token, nil
}

func (s *authService) JWTUserExpire(ctx context.Context, user *usermodel.User) (int, error) {
	key := s.getCacheKey(ctx, AuthenticationToken, strconv.Itoa(user.Id))
	s.deleteKey(ctx, key)
	return user.Id, nil
}

func (s *authService) Register(ctx context.Context, payload *authmodel.AuthRegister) (int, error) {
	condFindWithEmail := map[string]interface{}{"email": payload.Email}
	user, err := s.userService.SearchUserTrace(ctx, condFindWithEmail)
	if user != nil && err == nil {
		return -1, common.ErrorEntityExisted(usermodel.EntityName, errors.New("user already exists"))
	}

	passwordSalt := s.hashService.GenerateRandomString(s.config.JWT.PasswordSaltLength)
	hashPassword := s.hashService.GenerateSHA256(payload.Password, passwordSalt)
	newUser := &usermodel.UserCreate{
		Email:        payload.Email,
		Password:     hashPassword,
		PasswordSalt: passwordSalt,
	}

	userId, createErr := s.userService.CreateUserTrace(ctx, newUser)
	if createErr != nil {
		return -1, createErr
	}

	return userId, nil
}

func (s *authService) Login(ctx context.Context, payload *authmodel.AuthLogin) (*authprovidermodel.TokenProvider, error) {
	condFindWithEmail := map[string]interface{}{"email": payload.Email}
	user, err := s.userService.SearchUserTrace(ctx, condFindWithEmail)
	if err != nil || user == nil {
		return nil, common.ErrorCannotFoundEntity(usermodel.EntityName, err)
	}

	hashPassword := s.hashService.GenerateSHA256(payload.Password, user.PasswordSalt)
	if user.Password != hashPassword {
		return nil, common.ErrorInvalidRequest(usermodel.EntityName, errors.New("email or password wrong"))
	}

	return s.JWTUserGenerate(ctx, user)
}

func (s *authService) Logout(ctx context.Context, user *usermodel.User) (int, error) {
	return s.JWTUserExpire(ctx, user)
}

func (s *authService) ForgotPassword(ctx context.Context, payload *authmodel.AuthForgotPassword) (int, error) {
	condFindWithEmail := map[string]interface{}{"email": payload.Email}
	user, err := s.userService.SearchUserTrace(ctx, condFindWithEmail)
	if err != nil || user == nil {
		return -1, common.ErrorCannotFoundEntity(usermodel.EntityName, err)
	}

	key := s.getCacheKey(ctx, ForgotPasswordToken, strconv.Itoa(user.Id))
	resetPasswordToken := s.hashService.GenerateRandomString(15)
	s.setSession(ctx, key, resetPasswordToken, s.config.JWT.ExpireDays)

	go s.mailService.SendMail(mailer.Mail{
		Sender:  s.config.Mail.Sender,
		To:      []string{user.Email},
		Subject: "Forgot Password",
		Body:    fmt.Sprintf("%s?token=%s", payload.ForgotPasswordUri, resetPasswordToken),
	})

	return user.Id, nil
}

func (s *authService) ResetPassword(ctx context.Context, payload *authmodel.AuthResetPassword) (int, error) {
	condFindWithEmail := map[string]interface{}{"email": payload.Email}
	user, err := s.userService.SearchUserTrace(ctx, condFindWithEmail)
	if err != nil || user == nil {
		return -1, common.ErrorCannotFoundEntity(usermodel.EntityName, err)
	}

	// validate token
	key := s.getCacheKey(ctx, ForgotPasswordToken, strconv.Itoa(user.Id))
	resetPasswordToken, tokenErr := s.getSession(ctx, key)
	if tokenErr != nil {
		return -1, common.ErrorNotFound("forgot password token", tokenErr)
	}

	if resetPasswordToken.(string) != payload.Token {
		return -1, common.ErrorInvalidRequest("token is invalid", errors.New("forgot password token is invalid"))
	}

	passwordSalt := s.hashService.GenerateRandomString(s.config.JWT.PasswordSaltLength)
	hashPassword := s.hashService.GenerateSHA256(payload.Password, passwordSalt)
	userUpdate := &usermodel.UserUpdate{
		Password:     hashPassword,
		PasswordSalt: passwordSalt,
	}

	_, updateErr := s.userService.UpdateUserTrace(ctx, user.Id, userUpdate)
	if updateErr != nil {
		return -1, updateErr
	}

	s.deleteKey(ctx, key)

	return user.Id, nil
}

func (s *authService) GoogleLogin(ctx context.Context, redirect string) string {
	state := uuid.New().String()
	key := s.getCacheKey(ctx, authprovidermodel.GoogleAuthProvider, state)
	s.setSession(ctx, key, map[string]interface{}{"redirect": redirect}, 1)
	return s.googleProvider.GetAuthUri(ctx, state)
}

func (s *authService) GoogleValidate(ctx context.Context, code string) (*authprovidermodel.TokenProvider, error) {
	googleUser, err := s.googleProvider.GetUserTrace(ctx, code)
	if err != nil || googleUser == nil {
		return nil, err
	}
	condFindWithEmail := map[string]interface{}{"email": googleUser.Email}
	user, err := s.userService.SearchUserTrace(ctx, condFindWithEmail)

	// if user is not exist, we should create the user
	if err != nil || user == nil {
		// random password
		password := s.hashService.GenerateRandomString(s.config.JWT.PasswordSaltLength)
		passwordSalt, hashPassword := s.hashService.HashPassword(password, passwordSalt)
		userId, err := s.userService.CreateUserTrace(ctx, &usermodel.UserCreate{
			Email:        googleUser.Email,
			FirstName:    googleUser.GivenName,
			LastName:     googleUser.FamilyName,
			Password:     hashPassword,
			PasswordSalt: passwordSalt,
		})
		if err != nil {
			return nil, err
		}

		user, err = s.userService.GetUserTrace(ctx, userId)
		if err != nil {
			return nil, err
		}
	}

	condFindProviderWithUserId := map[string]interface{}{"user_id": user.Id, "provider_name": authprovidermodel.GoogleAuthProvider}
	provider, _ := s.authProviderService.SearchOneTrace(ctx, condFindProviderWithUserId)

	if provider == nil {
		_, err = s.authProviderService.CreateTrace(ctx, &authprovidermodel.AuthProviderCreate{
			UserId:       user.Id,
			ProviderName: authprovidermodel.GoogleAuthProvider,
			ProviderId:   googleUser.Id,
		})
		if err != nil {
			return nil, err
		}
	}

	return s.JWTUserGenerate(ctx, user)
}
