package middleware

import (
	"user_management/components/appctx"
	esprovider "user_management/components/elasticsearch"

	"github.com/gin-gonic/gin"
)

func SetElasticSearch(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := esprovider.WithContext(ginCtx.Request.Context(), appCtx.GetESService())
		ginCtx.Request = ginCtx.Request.WithContext(ctx)
		ginCtx.Next()
	}
}
