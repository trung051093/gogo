package auth

import (
	"net/http"
	"user_management/common"
	"user_management/components/appctx"
	"user_management/components/hasher"
	authmodel "user_management/modules/auth/model"
	jwtauthprovider "user_management/modules/auth_providers/jwt"
	"user_management/modules/user"
	usermodel "user_management/modules/user/model"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	ctx *gin.Context
}

func RegisterUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		var newData authmodel.AuthRegister

		if err := ginCtx.ShouldBind(&newData); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		validate := appCtx.GetValidator()
		if err := validate.Struct(&newData); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		appConfig := appCtx.GetConfig()
		userRepo := user.NewUserRepository(appCtx.GetMainDBConnection())
		hashService := hasher.NewHashService()
		userService := user.NewUserService(userRepo)
		jwtProvider := jwtauthprovider.NewJWTProvider(appCtx.GetConfig().JWT.Secret)
		authService := NewAuthService(jwtProvider, userService, hashService, appConfig)

		err := authService.Register(ginCtx, &newData)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SimpleSuccessResponse(true))
	}
}

func LoginUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		var loginData authmodel.AuthLogin

		if err := ginCtx.ShouldBind(&loginData); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		validate := appCtx.GetValidator()
		if err := validate.Struct(&loginData); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		appConfig := appCtx.GetConfig()
		userRepo := user.NewUserRepository(appCtx.GetMainDBConnection())
		userService := user.NewUserService(userRepo)
		hashService := hasher.NewHashService()
		jwtProvider := jwtauthprovider.NewJWTProvider(appConfig.JWT.Secret)
		authService := NewAuthService(jwtProvider, userService, hashService, appConfig)

		token, err := authService.Login(ginCtx, &loginData)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SimpleSuccessResponse(token))
	}
}
