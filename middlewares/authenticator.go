package middleware

import (
	"strings"
	"user_management/common"
	"user_management/components/appctx"
	jwtauthprovider "user_management/modules/auth_providers/jwt"

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
