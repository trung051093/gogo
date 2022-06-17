package auth

import (
	"context"
	"errors"
	"user_management/common"
	"user_management/components/appctx"
	"user_management/components/hasher"
	authmodel "user_management/modules/auth/model"
	authprovider "user_management/modules/auth_providers"
	jwtauthprovider "user_management/modules/auth_providers/jwt"
	usermodel "user_management/modules/user/model"

	"user_management/modules/user"
)

type authService struct {
	jwtProvider jwtauthprovider.JWTProvider
	userService user.UserService
	hashService hasher.HashService
	config      *appctx.Config
}

func NewAuthService(
	jwtProvider jwtauthprovider.JWTProvider,
	userService user.UserService,
	hashService hasher.HashService,
	config *appctx.Config,
) *authService {
	return &authService{jwtProvider: jwtProvider, userService: userService, hashService: hashService, config: config}
}

func (s *authService) Register(ctx context.Context, payload *authmodel.AuthRegister) error {
	condFindWithEmail := map[string]interface{}{"email": payload.Email}
	user, searchErr := s.userService.SearchUser(ctx, condFindWithEmail)
	if user != nil && searchErr == nil {
		return common.ErrorEntityExisted(usermodel.EntityName, errors.New("user already exists"))
	}

	passwordSalt := s.hashService.GenerateRandomString(s.config.JWT.PasswordSaltLength)
	hashPassword := s.hashService.GenerateSHA256(payload.Password, passwordSalt)
	newUser := &usermodel.UserCreate{
		Email:        payload.Email,
		Password:     hashPassword,
		PasswordSalt: passwordSalt,
	}

	_, createErr := s.userService.CreateUser(ctx, newUser)

	if createErr != nil {
		return createErr
	}

	return nil
}

func (s *authService) Login(ctx context.Context, payload *authmodel.AuthLogin) (*authprovider.TokenProvider, error) {
	condFindWithEmail := map[string]interface{}{"email": payload.Email}
	user, searchErr := s.userService.SearchUser(ctx, condFindWithEmail)
	if searchErr != nil {
		return nil, common.ErrorCannotFoundEntity(usermodel.EntityName, searchErr)
	}

	hashPassword := s.hashService.GenerateSHA256(payload.Password, user.PasswordSalt)
	if user.Password != hashPassword {
		return nil, common.ErrorInvalidRequest(usermodel.EntityName, errors.New("email or password wrong"))
	}

	token, gererateTokenErr := s.jwtProvider.Generate(authprovider.TokenPayload{
		UserId: *user.Id,
		Email:  user.Email,
		Role:   user.Role,
	}, uint(s.config.JWT.ExpireDays))

	if gererateTokenErr != nil {
		return nil, gererateTokenErr
	}

	return token, nil
}
