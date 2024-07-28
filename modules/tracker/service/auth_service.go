package service

import (
	"context"
	"errors"
	"gogo/common"
	"gogo/components/appctx"
	cacheprovider "gogo/components/cache"
	"gogo/components/hasher"
	"gogo/modules/tracker/dto"
	"gogo/modules/tracker/entity"
	"gogo/modules/tracker/repository"
)

type AuthService interface {
	common.Service[entity.Auth, repository.AuthRepository]
	Register(ctx context.Context, payload *dto.AuthRegisterDto) (*entity.TokenProvider, error)
	Login(ctx context.Context, payload *dto.AuthLoginDto) (*entity.TokenProvider, error)
	Logout(ctx context.Context, payload *dto.AuthLogoutDto) error
	JWTValidate(ctx context.Context, session string) (*entity.TokenPayload, error)
}

type authService struct {
	common.Service[entity.Auth, repository.AuthRepository]
	userRepo     repository.UserRepository
	config       *appctx.Config
	jwtProvider  JWTProvider
	hashService  hasher.HashService
	cacheService cacheprovider.CacheService
}

func NewAuthService(
	authRepo repository.AuthRepository,
	userRepo repository.UserRepository,
	config *appctx.Config,
	jwtProvider JWTProvider,
	hashService hasher.HashService,
	cacheService cacheprovider.CacheService,
) AuthService {
	service := common.NewService(authRepo)
	return &authService{
		Service:      service,
		userRepo:     userRepo,
		config:       config,
		jwtProvider:  jwtProvider,
		hashService:  hashService,
		cacheService: cacheService,
	}
}

func NewAuthServiceWithAppCtx(appCtx appctx.AppContext) AuthService {
	appConfig := appCtx.GetConfig()
	authRepo := repository.NewAuthRepository(appCtx.GetMainDBConnection())
	userRepo := repository.NewUserRepository(appCtx.GetMainDBConnection())
	cacheService := appCtx.GetCacheService()
	hashService := hasher.NewHashService()
	jwtProvider := NewJWTProvider(appConfig.JWT.Secret)

	authService := NewAuthService(
		authRepo,
		userRepo,
		appConfig,
		jwtProvider,
		hashService,
		cacheService,
	)
	return authService
}

func (s *authService) JWTGenerate(ctx context.Context, auth *entity.Auth, user *entity.User) (*entity.TokenProvider, error) {
	session := s.hashService.GenerateRandomString(32)
	token, err := s.jwtProvider.Generate(entity.TokenPayload{
		Session:  session,
		AuthId:   auth.Id.String(),
		UserId:   auth.UserId.String(),
		Username: auth.Username,
		Role:     user.Role.String(),
	}, uint(s.config.JWT.ExpireDays))
	if err != nil || token == nil {
		return nil, err
	}

	key := s.getAuthenicationKey(session)
	s.setJwtSession(ctx, key, token.Token)

	return token, nil
}

func (s *authService) JWTExpire(ctx context.Context, session string) error {
	key := s.getAuthenicationKey(session)
	s.deleteKey(ctx, key)
	return nil
}

func (s *authService) JWTValidate(ctx context.Context, session string) (*entity.TokenPayload, error) {
	var jwt string
	key := s.getAuthenicationKey(session)
	jwt, err := s.getJwtSession(ctx, key)
	common.PanicIf(err != nil || jwt == "", common.ErrorUnauthorized(err))

	return s.jwtProvider.Validate(jwt)
}

func (s *authService) Register(ctx context.Context, registerDto *dto.AuthRegisterDto) (*entity.TokenProvider, error) {
	existAuth, err := s.GetRepository().FindByUserName(ctx, registerDto.Username)
	common.PanicIf(err != nil, common.ErrorCannotFoundEntity(entity.Auth{}.EntityName(), err))
	common.PanicIf(existAuth != nil, common.ErrorEntityExisted(entity.Auth{}.EntityName(), errors.New("user already exists")))

	newUser, err := s.userRepo.Create(ctx, &entity.User{
		Name:  registerDto.Username,
		Phone: registerDto.Phone,
		Role:  entity.AdminRole,
	})
	common.PanicIf(err != nil || newUser == nil, common.ErrorCannotCreateEntity(entity.Auth{}.EntityName(), err))

	passwordSalt := s.hashService.GenerateRandomString(s.config.JWT.PasswordSaltLength)
	hashPassword := s.hashService.GenerateSHA256(registerDto.Password, passwordSalt)
	newAuth, err := s.GetRepository().Create(ctx, &entity.Auth{
		Username:     registerDto.Username,
		UserId:       newUser.Id,
		Password:     hashPassword,
		PasswordSalt: passwordSalt,
	})
	common.PanicIf(err != nil || newAuth == nil, common.ErrorCannotCreateEntity(entity.Auth{}.EntityName(), err))

	return s.JWTGenerate(ctx, newAuth, newUser)
}

func (s *authService) Login(ctx context.Context, loginDto *dto.AuthLoginDto) (*entity.TokenProvider, error) {
	existAuth, err := s.GetRepository().FindByUserName(ctx, loginDto.Username)
	common.PanicIf(err != nil || existAuth == nil, common.ErrorCannotFoundEntity(entity.Auth{}.EntityName(), err))

	user, err := s.userRepo.FindById(ctx, existAuth.UserId.String())
	common.PanicIf(err != nil || user == nil, common.ErrorCannotFoundEntity(entity.User{}.EntityName(), err))

	hashPassword := s.hashService.GenerateSHA256(loginDto.Password, existAuth.PasswordSalt)
	common.PanicIf(existAuth.Password != hashPassword, common.ErrorInvalidRequest(entity.Auth{}.EntityName(), errors.New("email or password wrong")))

	return s.JWTGenerate(ctx, existAuth, user)
}

func (s *authService) Logout(ctx context.Context, payload *dto.AuthLogoutDto) error {
	return s.JWTExpire(ctx, payload.Session)
}
