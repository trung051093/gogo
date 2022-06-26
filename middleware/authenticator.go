package middleware

import (
	"gogo/common"
	"gogo/components/appctx"
	jwtauthprovider "gogo/modules/auth_providers/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTRequireHandler(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		auth := ginCtx.Request.Header.Get("Authorization")
		unauthorizedError := common.ErrorUnauthorized()
		if auth == "" {
			panic(unauthorizedError)
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			panic(unauthorizedError)
		}
		jwtProvider := jwtauthprovider.NewJWTProvider(appCtx.GetConfig().JWT.Secret)
		tokenPayload, err := jwtProvider.Validate(token)
		if err != nil {
			panic(unauthorizedError)
		}
		ginCtx.Set(common.CurrentUser, tokenPayload)
		ginCtx.Next()
	}
}
