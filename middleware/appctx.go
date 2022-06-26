package middleware

import (
	"gogo/components/appctx"

	"github.com/gin-gonic/gin"
)

func SetAppContextIntoRequest(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := appctx.WithContext(ginCtx.Request.Context(), appCtx)
		ginCtx.Request = ginCtx.Request.WithContext(ctx)
		ginCtx.Next()
	}
}
