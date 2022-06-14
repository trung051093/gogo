package middleware

import (
	"user_management/components/appctx"
	socketprovider "user_management/components/socketio"

	"github.com/gin-gonic/gin"
)

func SetSocketIO(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		ctx := socketprovider.WithContext(ginCtx.Request.Context(), appCtx.GetSocketService())
		ginCtx.Request = ginCtx.Request.WithContext(ctx)
		ginCtx.Next()
	}
}
