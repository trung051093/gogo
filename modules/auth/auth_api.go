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

// Register godoc
// @Summary      Register
// @Description  Register
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      authmodel.AuthRegister      true  "register"
// @Success      200   {object}  common.Response{data=bool}  "desc"
// @Failure      400   {object}  common.AppError
// @Router       /api/v1/auth/register [post]
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
		esService := appCtx.GetESService()
		hashService := hasher.NewHashService()
		userService := user.NewUserService(userRepo, esService)
		jwtProvider := jwtauthprovider.NewJWTProvider(appCtx.GetConfig().JWT.Secret)
		authService := NewAuthService(jwtProvider, userService, hashService, appConfig)

		err := authService.Register(ginCtx.Request.Context(), &newData)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(true))
	}
}

// Login godoc
// @Summary      Login
// @Description  Login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      authmodel.AuthLogin                           true  "login"
// @Success      200   {object}  common.Response{data=authmodel.AuthResponse}  "desc"
// @Failure      400   {object}  common.AppError
// @Router       /api/v1/auth/login [post]
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
		esService := appCtx.GetESService()
		userService := user.NewUserService(userRepo, esService)
		hashService := hasher.NewHashService()
		jwtProvider := jwtauthprovider.NewJWTProvider(appConfig.JWT.Secret)
		authService := NewAuthService(jwtProvider, userService, hashService, appConfig)

		token, err := authService.Login(ginCtx.Request.Context(), &loginData)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(token))
	}
}
