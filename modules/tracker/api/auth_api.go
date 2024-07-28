package api

import (
	"gogo/common"
	"gogo/components/appctx"
	"gogo/modules/tracker/dto"
	"gogo/modules/tracker/entity"
	"gogo/modules/tracker/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register godoc
//	@Summary		Register
//	@Description	Register
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.AuthRegisterDto				true	"register"
//	@Success		200		{object}	common.Response{data=string}	"desc"
//	@Failure		400		{object}	common.AppError
//	@Router			/api/v1/auth/register [post]
func RegisterUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		authRegisterDto := &dto.AuthRegisterDto{}

		err := common.ParseRequest[dto.AuthRegisterDto](appCtx, ginCtx)(authRegisterDto)
		common.PanicIf(err != nil, common.ErrorInvalidRequest(entity.Auth{}.EntityName(), err))

		authService := service.NewAuthServiceWithAppCtx(appCtx)
		tokenProvider, err := authService.Register(ginCtx.Request.Context(), authRegisterDto)
		common.PanicIf(err != nil, err)

		cfg := appCtx.GetConfig()
		ginCtx.SetCookie(cfg.Auth.SessionName, tokenProvider.Session, cfg.Auth.Expire, "/", cfg.Auth.Client, false, true)
		ginCtx.JSON(http.StatusCreated, common.SuccessResponse("OK"))
	}
}

// Login godoc
//	@Summary		Login
//	@Description	Login
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		dto.AuthLoginDto				true	"login"
//	@Success		200		{object}	common.Response{data=string}	"desc"
//	@Failure		400		{object}	common.AppError
//	@Router			/api/v1/auth/login [post]
func LoginUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		authLoginDto := &dto.AuthLoginDto{}
		err := common.ParseRequest[dto.AuthLoginDto](appCtx, ginCtx)(authLoginDto)
		common.PanicIf(err != nil, common.ErrorInvalidRequest(entity.Auth{}.EntityName(), err))

		authService := service.NewAuthServiceWithAppCtx(appCtx)
		tokenProvider, err := authService.Login(ginCtx.Request.Context(), authLoginDto)
		common.PanicIf(err != nil, err)

		cfg := appCtx.GetConfig()
		ginCtx.SetCookie(cfg.Auth.SessionName, tokenProvider.Session, cfg.Auth.Expire, "/", cfg.Auth.Client, false, true)
		ginCtx.JSON(http.StatusOK, common.SuccessResponse("OK"))
	}
}

// Logout godoc
//	@Summary		Logout
//	@Description	Logout
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{object}	common.Response{data=string}	"desc"
//	@Failure		400	{object}	common.AppError
//	@Router			/api/v1/auth/logout [post]
func LogoutUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		tokenPayload := ginCtx.Value(common.CurrentAuth).(*entity.TokenPayload)
		authLogoutDto := &dto.AuthLogoutDto{
			Session: tokenPayload.Session,
		}
		err := common.ParseRequest[dto.AuthLogoutDto](appCtx, ginCtx)(authLogoutDto)
		common.PanicIf(err != nil, common.ErrorInvalidRequest(entity.Auth{}.EntityName(), err))

		authService := service.NewAuthServiceWithAppCtx(appCtx)
		err = authService.Logout(ginCtx.Request.Context(), authLogoutDto)
		common.PanicIf(err != nil, err)

		cfg := appCtx.GetConfig()
		ginCtx.SetCookie(cfg.Auth.SessionName, "", -1, "/", cfg.Auth.Client, false, true)
		ginCtx.JSON(http.StatusOK, common.SuccessResponse("OK"))
	}
}
