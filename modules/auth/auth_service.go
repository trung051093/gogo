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
	googleauthprovider "gogo/modules/auth/providers/google"
	jwtauthprovider "gogo/modules/auth/providers/jwt"
	usermodel "gogo/modules/user/model"
	"net/url"

	"gogo/modules/user"

	"github.com/google/uuid"
)

type AuthService interface {
	common.Service[authmodel.Auth, AuthRepository]

	Register(ctx context.Context, payload *authmodel.AuthRegisterDto) (*authmodel.Auth, error)
	Login(ctx context.Context, payload *authmodel.AuthLoginDto) (*jwtauthprovider.TokenProvider, error)
	Logout(ctx context.Context, payload *jwtauthprovider.TokenPayload) (string, error)
	ForgotPassword(ctx context.Context, payload *authmodel.AuthForgotPasswordDto) (string, error)
	ResetPassword(ctx context.Context, payload *authmodel.AuthResetPasswordDto) (string, error)

	GoogleLogin(ctx context.Context, state string) string
	GoogleCallback(ctx context.Context, code, state string) (string, error)
}

type authService struct {
	common.Service[authmodel.Auth, AuthRepository]
	authRepo         AuthRepository
	authProviderRepo AuthProviderRepository
	config           *appctx.Config
	jwtProvider      jwtauthprovider.JWTProvider
	googleProvider   googleauthprovider.GoogleAuthProvider
	userService      user.UserService
	hashService      hasher.HashService
	cacheService     cacheprovider.CacheService
	mailService      mailer.MailService
}

func NewAuthService(
	authRepo AuthRepository,
	authProviderRepo AuthProviderRepository,
	config *appctx.Config,
	jwtProvider jwtauthprovider.JWTProvider,
	googleProvider googleauthprovider.GoogleAuthProvider,
	userService user.UserService,
	hashService hasher.HashService,
	cacheService cacheprovider.CacheService,
	mailService mailer.MailService,
) AuthService {
	service := common.NewService[authmodel.Auth, AuthRepository](authRepo)
	return &authService{
		Service:          service,
		authRepo:         authRepo,
		authProviderRepo: authProviderRepo,
		config:           config,
		jwtProvider:      jwtProvider,
		googleProvider:   googleProvider,
		userService:      userService,
		hashService:      hashService,
		cacheService:     cacheService,
		mailService:      mailService,
	}
}

func NewAuthServiceWithAppCtx(appCtx appctx.AppContext) AuthService {
	appConfig := appCtx.GetConfig()
	authRepo := NewAuthRepository(appCtx.GetMainDBConnection())
	authProviderRepo := NewAuthProviderRepository(appCtx.GetMainDBConnection())
	userRepo := user.NewUserRepository(appCtx.GetMainDBConnection())
	esService := appCtx.GetESService()
	cacheService := appCtx.GetCacheService()
	mailService := appCtx.GetMailService()
	hashService := hasher.NewHashService()
	userService := user.NewUserService(userRepo, appConfig, esService)
	jwtProvider := jwtauthprovider.NewJWTProvider(appConfig.JWT.Secret)
	googleProvider := googleauthprovider.NewGoogleAuthProvider((*googleauthprovider.GoogleAuthConfig)(&appConfig.GoogleOauth2))

	authService := NewAuthService(
		authRepo,
		authProviderRepo,
		appConfig,
		jwtProvider,
		googleProvider,
		userService,
		hashService,
		cacheService,
		mailService,
	)
	return authService
}

func (s *authService) JWTGenerate(ctx context.Context, auth *authmodel.Auth) (*jwtauthprovider.TokenProvider, error) {
	token, err := s.jwtProvider.Generate(jwtauthprovider.TokenPayload{
		AuthId: auth.Id.String(),
		UserId: auth.UserId,
		Email:  auth.Email,
		Role:   auth.User.Role,
	}, uint(s.config.JWT.ExpireDays))

	if err != nil || token == nil {
		return nil, err
	}

	// store to cache
	key := s.getCacheKey(ctx, AuthenticationToken, auth.Id.String())
	s.setSession(ctx, key, token.Token, s.config.JWT.ExpireDays)

	return token, nil
}

func (s *authService) JWTExpire(ctx context.Context, payload *jwtauthprovider.TokenPayload) (string, error) {
	key := s.getCacheKey(ctx, AuthenticationToken, payload.AuthId)
	s.deleteKey(ctx, key)
	return payload.AuthId, nil
}

func (s *authService) Register(ctx context.Context, registerDto *authmodel.AuthRegisterDto) (*authmodel.Auth, error) {
	existAuth, err := s.GetRepository().FindByUserName(ctx, registerDto.Email)
	if err != nil {
		return nil, common.ErrorCannotFoundEntity(authmodel.Auth{}.EntityName(), err)
	}
	if existAuth != nil {
		return nil, common.ErrorEntityExisted(authmodel.Auth{}.EntityName(), errors.New("user already exists"))
	}

	// create new user
	newUser, err := s.userService.Create(ctx, registerDto.ToCreateUserDto())
	if err != nil {
		return nil, err
	}

	// create new auth
	passwordSalt := s.hashService.GenerateRandomString(s.config.JWT.PasswordSaltLength)
	hashPassword := s.hashService.GenerateSHA256(registerDto.Password, passwordSalt)
	newAuth := &authmodel.Auth{
		Email:        registerDto.Email,
		UserId:       newUser.Id.String(),
		Password:     hashPassword,
		PasswordSalt: passwordSalt,
	}
	_, err = s.GetRepository().Create(ctx, newAuth)
	if err != nil {
		return nil, err
	}

	return newAuth, nil
}

func (s *authService) Login(ctx context.Context, loginDto *authmodel.AuthLoginDto) (*jwtauthprovider.TokenProvider, error) {
	existAuth, err := s.GetRepository().FindByUserName(ctx, loginDto.UserName)
	if err != nil || existAuth == nil {
		return nil, common.ErrorCannotFoundEntity(authmodel.Auth{}.EntityName(), err)
	}

	hashPassword := s.hashService.GenerateSHA256(loginDto.Password, existAuth.PasswordSalt)
	if existAuth.Password != hashPassword {
		return nil, common.ErrorInvalidRequest(authmodel.Auth{}.EntityName(), errors.New("email or password wrong"))
	}

	return s.JWTGenerate(ctx, existAuth)
}

func (s *authService) Logout(ctx context.Context, payload *jwtauthprovider.TokenPayload) (string, error) {
	return s.JWTExpire(ctx, payload)
}

func (s *authService) ForgotPassword(ctx context.Context, forgotPasswordDto *authmodel.AuthForgotPasswordDto) (string, error) {
	existAuth, err := s.GetRepository().FindByEmail(ctx, forgotPasswordDto.Email)
	if err != nil || existAuth == nil {
		return "", common.ErrorCannotFoundEntity(authmodel.Auth{}.EntityName(), err)
	}

	key := s.getCacheKey(ctx, ForgotPasswordToken, existAuth.UserId, existAuth.Email)
	resetPasswordToken := s.hashService.GenerateRandomString(15)
	s.setSession(ctx, key, map[string]string{
		"token": resetPasswordToken,
		"email": existAuth.Email,
	}, s.config.JWT.ExpireDays)

	resetPasswordUri, resetPasswordUriErr := url.Parse(forgotPasswordDto.ForgotPasswordUri)
	if resetPasswordUriErr != nil {
		panic(common.ErrorInvalidRequest("google auth invalid redirect uri", err))
	}

	q := resetPasswordUri.Query()
	q.Set("token", resetPasswordToken)
	q.Set("email", existAuth.Email)
	resetPasswordUri.RawQuery = q.Encode()

	go s.mailService.SendMail(mailer.Mail{
		Sender:  s.config.Mail.Sender,
		To:      []string{existAuth.Email},
		Subject: "Forgot Password",
		Body:    resetPasswordUri.String(),
	})

	return existAuth.Email, nil
}

func (s *authService) ResetPassword(ctx context.Context, resetPasswordDto *authmodel.AuthResetPasswordDto) (string, error) {
	existAuth, err := s.GetRepository().FindByEmail(ctx, resetPasswordDto.Email)
	if err != nil || existAuth == nil {
		return "", common.ErrorCannotFoundEntity(authmodel.Auth{}.EntityName(), err)
	}

	// validate token, alway get key from payload email to avoid hack.
	key := s.getCacheKey(ctx, ForgotPasswordToken, existAuth.UserId, existAuth.Email)
	session, tokenErr := s.getSession(ctx, key)
	if tokenErr != nil {
		return "", common.ErrorNotFound("forgot password token", tokenErr)
	}

	token := session.(map[string]interface{})["token"].(string)
	email := session.(map[string]interface{})["email"].(string)

	if token != resetPasswordDto.Token {
		return "", common.ErrorInvalidRequest("token is invalid", errors.New("forgot password token is invalid"))
	}

	if existAuth.Email != email {
		return "", common.ErrorInvalidRequest("email is invalid", errors.New("email is invalid"))
	}

	passwordSalt := s.hashService.GenerateRandomString(s.config.JWT.PasswordSaltLength)
	hashPassword := s.hashService.GenerateSHA256(resetPasswordDto.Password, passwordSalt)
	existAuth.Password = hashPassword
	existAuth.PasswordSalt = passwordSalt

	if _, err := s.GetRepository().Updates(ctx, existAuth, map[string]interface{}{
		"password":      hashPassword,
		"password_salt": passwordSalt,
	}); err != nil {
		return "", err
	}

	// delete reset password key
	s.deleteKey(ctx, key)

	return existAuth.Id.String(), nil
}

func (s *authService) GoogleLogin(ctx context.Context, redirect string) string {
	state := uuid.New().String()
	key := s.getCacheKey(ctx, string(authmodel.GoogleAuthProvider), state)
	s.setSession(ctx, key, map[string]interface{}{"redirect": redirect}, 1)
	return s.googleProvider.GetAuthUri(ctx, state)
}

func (s *authService) GoogleCallback(ctx context.Context, code, state string) (string, error) {
	key := s.getCacheKey(ctx, string(authmodel.GoogleAuthProvider), state)
	session, err := s.getSession(ctx, key)
	if err != nil {
		return "", common.ErrorInvalidRequest("google auth invalid session", err)
	}

	redirect := session.(map[string]interface{})["redirect"].(string)
	redirectUri, err := url.Parse(redirect)
	if err != nil {
		return "", common.ErrorInvalidRequest("google auth invalid redirect uri", err)
	}

	googleUser, err := s.googleProvider.GetUserTrace(ctx, code)
	if err != nil || googleUser == nil {
		return "", err
	}
	existAuth, err := s.GetRepository().FindByEmail(ctx, googleUser.Email)
	if err != nil {
		return "", common.ErrorCannotFoundEntity(authmodel.Auth{}.EntityName(), err)
	}

	// if user is not exist, we should create the user
	if existAuth == nil {
		newUser, err := s.userService.Create(ctx, &usermodel.UserCreateDto{
			Email:     googleUser.Email,
			FirstName: googleUser.GivenName,
			LastName:  googleUser.FamilyName,
		})
		if err != nil {
			return "", err
		}

		_, err = s.GetRepository().Create(ctx, &authmodel.Auth{
			Email:  googleUser.Email,
			UserId: newUser.Id.String(),
		})
		if err != nil {
			return "", err
		}
	}

	auth, err := s.GetRepository().FindByUserId(ctx, existAuth.UserId)
	if err != nil {
		return "", common.ErrorCannotFoundEntity(authmodel.Auth{}.EntityName(), err)
	}

	provider, _ := s.authProviderRepo.FindById(ctx, auth.Id)
	if err != nil {
		return "", common.ErrorCannotFoundEntity(authmodel.AuthProvider{}.EntityName(), err)
	}

	if provider == nil {
		_, err = s.authProviderRepo.Create(ctx, &authmodel.AuthProvider{
			ProviderName: authmodel.GoogleAuthProvider,
			ProviderId:   googleUser.Id,
		})
		if err != nil {
			return "", common.ErrorCannotCreateEntity(authmodel.AuthProvider{}.EntityName(), err)
		}
	}

	tokenProvider, err := s.JWTGenerate(ctx, existAuth)
	if err != nil {
		return "", common.ErrorCannotFoundEntity(authmodel.AuthProvider{}.EntityName(), err)
	}

	q := redirectUri.Query()
	q.Set("token", tokenProvider.Token)
	q.Set("expiry", fmt.Sprintf("%d", tokenProvider.Expiry))
	redirectUri.RawQuery = q.Encode()

	s.deleteKey(ctx, key)

	return redirectUri.String(), nil
}
