package middleware

import (
	"errors"
	"gogo/common"
	"gogo/components/appctx"
	jwtauthprovider "gogo/modules/auth_provider/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTRequireHandler(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		auth := ginCtx.Request.Header.Get("Authorization")
		if auth == "" {
			panic(common.ErrorUnauthorized(errors.New("cannot found authentication header")))
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		if token == auth {
			panic(common.ErrorUnauthorized(errors.New("cannot found token")))
		}
		jwtProvider := jwtauthprovider.NewJWTProvider(appCtx.GetConfig().JWT.Secret)
		tokenPayload, err := jwtProvider.Validate(token)
		if err != nil {
			panic(common.ErrorUnauthorized(errors.New("token invalid")))
		}
		ginCtx.Set(common.CurrentUser, tokenPayload)
		ginCtx.Next()
	}
}
