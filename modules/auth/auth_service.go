package auth

import (
	"context"
	"errors"
	"gogo/common"
	"gogo/components/appctx"
	cacheprovider "gogo/components/cache"
	"gogo/components/hasher"
	"gogo/components/mailer"
	authmodel "gogo/modules/auth/model"
	authprovider "gogo/modules/auth_providers"
	jwtauthprovider "gogo/modules/auth_providers/jwt"
	usermodel "gogo/modules/user/model"
	"strconv"

	"gogo/modules/user"
)

type authService struct {
	jwtProvider  jwtauthprovider.JWTProvider
	userService  user.UserService
	hashService  hasher.HashService
	cacheService cacheprovider.CacheService
	mailService  mailer.MailService
	config       *appctx.Config
}

func NewAuthService(
	config *appctx.Config,
	jwtProvider jwtauthprovider.JWTProvider,
	userService user.UserService,
	hashService hasher.HashService,
	cacheService cacheprovider.CacheService,
	mailService mailer.MailService,
) *authService {
	return &authService{
		config:       config,
		jwtProvider:  jwtProvider,
		userService:  userService,
		hashService:  hashService,
		cacheService: cacheService,
		mailService:  mailService,
	}
}

func NewAuthServiceFromContext(appCtx appctx.AppContext) *authService {
	appConfig := appCtx.GetConfig()
	userRepo := user.NewUserRepository(appCtx.GetMainDBConnection())
	esService := appCtx.GetESService()
	cacheService := appCtx.GetCacheService()
	mailService := appCtx.GetMailService()
	userService := user.NewUserService(userRepo, esService)
	hashService := hasher.NewHashService()
	jwtProvider := jwtauthprovider.NewJWTProvider(appConfig.JWT.Secret)
	authService := NewAuthService(appConfig, jwtProvider, userService, hashService, cacheService, mailService)
	return authService
}

func (s *authService) Register(ctx context.Context, payload *authmodel.AuthRegister) (int, error) {
	condFindWithEmail := map[string]interface{}{"email": payload.Email}
	user, searchErr := s.userService.SearchUserTrace(ctx, condFindWithEmail)
	if user != nil && searchErr == nil {
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

func (s *authService) Login(ctx context.Context, payload *authmodel.AuthLogin) (*authprovider.TokenProvider, error) {
	condFindWithEmail := map[string]interface{}{"email": payload.Email}
	user, searchErr := s.userService.SearchUserTrace(ctx, condFindWithEmail)
	if searchErr != nil || user == nil {
		return nil, common.ErrorCannotFoundEntity(usermodel.EntityName, searchErr)
	}

	hashPassword := s.hashService.GenerateSHA256(payload.Password, user.PasswordSalt)
	if user.Password != hashPassword {
		return nil, common.ErrorInvalidRequest(usermodel.EntityName, errors.New("email or password wrong"))
	}

	token, gererateTokenErr := s.jwtProvider.Generate(authprovider.TokenPayload{
		UserId: user.Id,
		Email:  user.Email,
		Role:   user.Role,
	}, uint(s.config.JWT.ExpireDays))

	if gererateTokenErr != nil {
		return nil, gererateTokenErr
	}

	tokenKey := s.getKeyToken(ctx, AuthenticationToken, strconv.Itoa(user.Id))
	s.setToken(ctx, tokenKey, token.Token, s.config.JWT.ExpireDays)

	return token, nil
}

func (s *authService) Logout(ctx context.Context, user *usermodel.User) (int, error) {
	// delete login token
	tokenKey := s.getKeyToken(ctx, AuthenticationToken, strconv.Itoa(user.Id))
	deleteTokenErr := s.deleteToken(ctx, tokenKey)
	if deleteTokenErr != nil {
		return -1, common.ErrorNotFound("token", deleteTokenErr)
	}
	return user.Id, nil
}

func (s *authService) ForgotPassword(ctx context.Context, payload *authmodel.AuthForgotPassword) (int, error) {
	condFindWithEmail := map[string]interface{}{"email": payload.Email}
	user, searchErr := s.userService.SearchUserTrace(ctx, condFindWithEmail)
	if searchErr != nil || user == nil {
		return -1, common.ErrorCannotFoundEntity(usermodel.EntityName, searchErr)
	}

	tokenKey := s.getKeyToken(ctx, ForgotPasswordToken, strconv.Itoa(user.Id))
	resetPasswordToken := s.hashService.GenerateRandomString(15)
	s.setToken(ctx, tokenKey, resetPasswordToken, s.config.JWT.ExpireDays)

	go s.mailService.SendMail(mailer.Mail{
		Sender:  s.config.Mail.Sender,
		To:      []string{user.Email},
		Subject: "Forgot Password",
		Body:    resetPasswordToken,
	})

	return user.Id, nil
}

func (s *authService) ResetPassword(ctx context.Context, payload *authmodel.AuthResetPassword) (int, error) {
	condFindWithEmail := map[string]interface{}{"email": payload.Email}
	user, searchErr := s.userService.SearchUserTrace(ctx, condFindWithEmail)
	if searchErr != nil || user == nil {
		return -1, common.ErrorCannotFoundEntity(usermodel.EntityName, searchErr)
	}

	// validate token
	tokenKey := s.getKeyToken(ctx, ForgotPasswordToken, strconv.Itoa(user.Id))
	resetPasswordToken, tokenErr := s.getToken(ctx, tokenKey)
	if tokenErr != nil {
		return -1, common.ErrorNotFound("forgot password token", tokenErr)
	}

	if resetPasswordToken != payload.Token {
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

	s.deleteToken(ctx, tokenKey)

	return user.Id, nil
}
