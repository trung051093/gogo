package middleware

import (
	"errors"
	"strings"
	"user_management/common"
	"user_management/components/appctx"
	jwtauthprovider "user_management/modules/auth_providers/jwt"

	"github.com/gin-gonic/gin"
)

func JWTRequireHandler(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		auth := ginCtx.Request.Header.Get("Authorization")
		unauthorizedError := common.NewUnauthorized(
			errors.New("Unauthorized !!!"),
			"Unauthorized",
			"Unauthorized",
		)
		if auth == "" {
			ginCtx.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
			panic(unauthorizedError)
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			ginCtx.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
			panic(unauthorizedError)
		}
		jwtProvider := jwtauthprovider.NewJWTProvider(appCtx.GetConfig().JWT.Secret)
		tokenPayload, err := jwtProvider.Validate(token)
		if err != nil {
			ginCtx.AbortWithStatusJSON(unauthorizedError.StatusCode, unauthorizedError)
			panic(unauthorizedError)
		}
		ginCtx.Set(common.CurrentUser, tokenPayload)
		ginCtx.Next()
	}
}
