package middleware

import (
	"net/http"
	"strings"
	"user_management/common"
	component "user_management/components"
	jwtauthprovider "user_management/modules/auth_providers/jwt"

	"github.com/gin-gonic/gin"
)

func JWTRequireHandler(appCtx component.AppContext) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		auth := ginCtx.Request.Header.Get("Authorization")

		if auth == "" {
			ginCtx.String(http.StatusForbidden, "No Authorization header provided")
			ginCtx.Abort()
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			ginCtx.String(http.StatusForbidden, "Could not find bearer token in Authorization header")
			ginCtx.Abort()
			return
		}
		jwtProvider := jwtauthprovider.NewJWTProvider(appCtx.GetConfig().JWT.Secret)
		tokenPayload, err := jwtProvider.Validate(token)
		if err != nil {
			ginCtx.AbortWithError(http.StatusInternalServerError, err)
			return
		}
		ginCtx.Set(common.CurrentUser, tokenPayload)
		ginCtx.Next()
	}
}
