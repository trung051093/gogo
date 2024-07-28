package middleware

import (
	"errors"
	"gogo/common"
	"gogo/components/appctx"
	"gogo/modules/tracker/service"

	"github.com/gin-gonic/gin"
)

func Authenticaticator(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		cfg := appCtx.GetConfig()
		cookieValue, err := ginCtx.Cookie(cfg.Auth.SessionName)
		common.PanicIf(err != nil || cookieValue == "", common.ErrorUnauthorized(errors.New("unauthorized")))

		authService := service.NewAuthServiceWithAppCtx(appCtx)
		tokenPayload, err := authService.JWTValidate(ginCtx.Request.Context(), cookieValue)
		common.PanicIf(err != nil || tokenPayload == nil, common.ErrorUnauthorized(errors.New("unauthorized")))

		ginCtx.Set(common.CurrentAuth, tokenPayload)
		ginCtx.Next()
	}
}
