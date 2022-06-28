package auth

import (
	"fmt"
	"gogo/common"
	"gogo/components/appctx"
	authmodel "gogo/modules/auth/model"
	authprovidermodel "gogo/modules/auth_provider/model"
	usermodel "gogo/modules/user/model"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

// Register godoc
// @Summary      Register
// @Description  Register
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      authmodel.AuthRegister     true  "register"
// @Success      200   {object}  common.Response{data=int}  "desc"
// @Failure      400  {object}  common.AppError
// @Router       /api/v1/auth/register [post]
func RegisterUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		var newData authmodel.AuthRegister

		if err := ginCtx.ShouldBind(&newData); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		validator := appCtx.GetValidator()
		if err := validator.Struct(&newData); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		authService := NewAuthServiceFromContext(appCtx)

		userId, err := authService.Register(ginCtx.Request.Context(), &newData)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(userId))
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

		validator := appCtx.GetValidator()
		if err := validator.Struct(&loginData); err != nil {
			panic(common.ErrorInvalidRequest(usermodel.EntityName, err))
		}

		authService := NewAuthServiceFromContext(appCtx)

		token, err := authService.Login(ginCtx.Request.Context(), &loginData)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(token))
	}
}

// Logout godoc
// @Security     ApiKeyAuth
// @Summary      Logout
// @Description  Logout
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200  {object}  common.Response{data=bool}  "desc"
// @Failure      400   {object}  common.AppError
// @Router       /api/v1/auth/logout [post]
func LogoutUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		currentUser := ginCtx.Value(common.CurrentUser).(usermodel.User)

		authService := NewAuthServiceFromContext(appCtx)

		_, err := authService.Logout(ginCtx.Request.Context(), &currentUser)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(true))
	}
}

// ForgotPassword godoc
// @Summary      ForgotPassword
// @Description  ForgotPassword
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      authmodel.AuthForgotPassword  true  "login"
// @Success      200   {object}  common.Response{data=bool}    "desc"
// @Failure      400   {object}  common.AppError
// @Router       /api/v1/auth/forgot-password [post]
func ForgotPasswordUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		var forgotPasswordData authmodel.AuthForgotPassword

		if err := ginCtx.ShouldBind(&forgotPasswordData); err != nil {
			panic(common.ErrorInvalidRequest(authmodel.EntityName, err))
		}

		validator := appCtx.GetValidator()
		if err := validator.Struct(&forgotPasswordData); err != nil {
			panic(common.ErrorInvalidRequest(authmodel.EntityName, err))
		}

		authService := NewAuthServiceFromContext(appCtx)

		_, err := authService.ForgotPassword(ginCtx.Request.Context(), &forgotPasswordData)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(true))
	}
}

// ResetPassword godoc
// @Summary      ResetPassword
// @Description  ResetPassword
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        user  body      authmodel.AuthResetPassword  true  "login"
// @Success      200   {object}  common.Response{data=bool}   "desc"
// @Failure      400   {object}  common.AppError
// @Router       /api/v1/auth/reset-password [post]
func ResetPasswordUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		var resetPasswordData authmodel.AuthResetPassword

		if err := ginCtx.ShouldBind(&resetPasswordData); err != nil {
			panic(common.ErrorInvalidRequest(authmodel.EntityName, err))
		}

		validator := appCtx.GetValidator()
		if err := validator.Struct(&resetPasswordData); err != nil {
			panic(common.ErrorInvalidRequest(authmodel.EntityName, err))
		}

		authService := NewAuthServiceFromContext(appCtx)

		_, err := authService.ResetPassword(ginCtx.Request.Context(), &resetPasswordData)

		if err != nil {
			panic(err)
		}

		ginCtx.JSON(http.StatusOK, common.SuccessResponse(true))
	}
}

// GoogleLogin godoc
// @Summary      Google login
// @Description  Google login
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        redirect      query     string                                                                                                        false   "redirect"
// @Success      307
// @Failure      400   {object}  common.AppError
// @Router       /api/v1/auth/google/login [get]
func GoogleLoginUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		var redirect string
		if err := ginCtx.ShouldBind(&redirect); err != nil {
			panic(common.ErrorInvalidRequest("redirect", err))
		}
		authService := NewAuthServiceFromContext(appCtx)
		url := authService.GoogleLogin(ginCtx.Request.Context(), redirect)
		ginCtx.Redirect(http.StatusTemporaryRedirect, url)
	}
}

// GoogleCallback godoc
// @Summary      Google callback
// @Description  Google callback
// @Tags         auth
// @Accept       json
// @Produce      json
// @Success      200   {object}  common.Response{data=authmodel.AuthResponse}  "desc"
// @Failure      400   {object}  common.AppError
// @Router       /api/v1/auth/google/callback [get]
func GoogleCallbackUserHandler(appCtx appctx.AppContext) func(*gin.Context) {
	return func(ginCtx *gin.Context) {
		var state string = ginCtx.Query("state")
		var code string = ginCtx.Query("code")
		ctx := ginCtx.Request.Context()

		authService := NewAuthServiceFromContext(appCtx)
		key := authService.getCacheKey(ctx, authprovidermodel.GoogleAuthProvider, state)
		session, err := authService.getSession(ctx, key)
		if err != nil {
			panic(common.ErrorInvalidRequest("google auth invalid session", err))
		}

		redirect := session.(map[string]interface{})["redirect"].(string)
		redirectUri, err := url.Parse(redirect)
		if err != nil {
			panic(common.ErrorInvalidRequest("google auth invalid redirect uri", err))
		}

		tokenProvider, err := authService.GoogleValidate(ctx, code)
		if err != nil {
			panic(common.ErrorInvalidRequest("google auth code", err))
		}

		q := redirectUri.Query()
		q.Set("token", tokenProvider.Token)
		q.Set("expiry", fmt.Sprintf("%d", tokenProvider.Expiry))
		redirectUri.RawQuery = q.Encode()

		authService.deleteKey(ctx, key)

		ginCtx.JSON(http.StatusOK, redirectUri.RequestURI())
	}
}
