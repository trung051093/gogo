package auth

import (
	"gogo/common"
	"gogo/components/appctx"
	authmodel "gogo/modules/auth/model"
	jwtauthprovider "gogo/modules/auth/providers/jwt"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Register godoc
//	@Summary		Register
//	@Description	Register
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		authmodel.AuthRegisterDto		true	"register"
//	@Success		200		{object}	common.Response{data=string}	"desc"
//	@Failure		400		{object}	common.AppError
//	@Router			/api/v1/auth/register [post]
func RegisterUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		authRegisterDto := &authmodel.AuthRegisterDto{}
		if err := common.ParseRequest[authmodel.AuthRegisterDto](appCtx, ginCtx)(authRegisterDto); err != nil {
			panic(common.ErrorInvalidRequest(authmodel.Auth{}.EntityName(), err))
		}

		authService := NewAuthServiceWithAppCtx(appCtx)
		auth, err := authService.Register(ginCtx.Request.Context(), authRegisterDto)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusCreated, common.SuccessResponse(auth.UserId))
	}
}

// Login godoc
//	@Summary		Login
//	@Description	Login
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		authmodel.AuthLoginDto							true	"login"
//	@Success		200		{object}	common.Response{data=authmodel.AuthResponseDto}	"desc"
//	@Failure		400		{object}	common.AppError
//	@Router			/api/v1/auth/login [post]
func LoginUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		authLoginDto := &authmodel.AuthLoginDto{}
		if err := common.ParseRequest[authmodel.AuthLoginDto](appCtx, ginCtx)(authLoginDto); err != nil {
			panic(common.ErrorInvalidRequest(authmodel.Auth{}.EntityName(), err))
		}

		authService := NewAuthServiceWithAppCtx(appCtx)
		token, err := authService.Login(ginCtx.Request.Context(), authLoginDto)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(token))
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
		tokenPayload := ginCtx.Value(common.CurrentAuth).(jwtauthprovider.TokenPayload)

		authService := NewAuthServiceWithAppCtx(appCtx)

		_, err := authService.Logout(ginCtx.Request.Context(), &tokenPayload)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse("OK"))
	}
}

// ForgotPassword godoc
//	@Summary		ForgotPassword
//	@Description	ForgotPassword
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		authmodel.AuthForgotPasswordDto	true	"login"
//	@Success		200		{object}	common.Response{data=string}	"desc"
//	@Failure		400		{object}	common.AppError
//	@Router			/api/v1/auth/forgot-password [post]
func ForgotPasswordUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		authForgotPasswordDto := &authmodel.AuthForgotPasswordDto{}
		if err := common.ParseRequest[authmodel.AuthForgotPasswordDto](appCtx, ginCtx)(authForgotPasswordDto); err != nil {
			panic(common.ErrorInvalidRequest(authmodel.Auth{}.EntityName(), err))
		}

		authService := NewAuthServiceWithAppCtx(appCtx)

		_, err := authService.ForgotPassword(ginCtx.Request.Context(), authForgotPasswordDto)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse("OK"))
	}
}

// ResetPassword godoc
//	@Summary		ResetPassword
//	@Description	ResetPassword
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			user	body		authmodel.AuthResetPasswordDto	true	"login"
//	@Success		200		{object}	common.Response{data=string}	"desc"
//	@Failure		400		{object}	common.AppError
//	@Router			/api/v1/auth/reset-password [post]
func ResetPasswordUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		authResetPasswordDto := &authmodel.AuthResetPasswordDto{}
		if err := common.ParseRequest[authmodel.AuthResetPasswordDto](appCtx, ginCtx)(authResetPasswordDto); err != nil {
			panic(common.ErrorInvalidRequest(authmodel.Auth{}.EntityName(), err))
		}

		authService := NewAuthServiceWithAppCtx(appCtx)

		_, err := authService.ResetPassword(ginCtx.Request.Context(), authResetPasswordDto)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse("OK"))
	}
}

// GoogleLogin godoc
//	@Summary		Google login
//	@Description	Google login
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			redirectUri	query	string	false	"redirect"
//	@Success		307
//	@Failure		400	{object}	common.AppError
//	@Router			/api/v1/auth/google/login [get]
func GoogleLoginUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		redirectUri := ginCtx.Query("redirectUri")
		authService := NewAuthServiceWithAppCtx(appCtx)
		url := authService.GoogleLogin(ginCtx.Request.Context(), redirectUri)
		ginCtx.Redirect(http.StatusTemporaryRedirect, url)
	}
}

// GoogleCallback godoc
//	@Summary		Google callback
//	@Description	Google callback
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	common.Response{data=authmodel.AuthResponseDto}	"desc"
//	@Failure		400	{object}	common.AppError
//	@Router			/api/v1/auth/google/callback [get]
func GoogleCallbackUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		state := ginCtx.Query("state")
		code := ginCtx.Query("code")
		ctx := ginCtx.Request.Context()
		authService := NewAuthServiceWithAppCtx(appCtx)
		redirectUri, err := authService.GoogleCallback(ctx, state, code)
		if err != nil {
			panic(err)
		}

		ginCtx.Redirect(http.StatusTemporaryRedirect, redirectUri)
	}
}
