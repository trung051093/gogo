package middleware

import (
	"context"
	"user_management/common"
	"user_management/components/appctx"
	"user_management/components/rabbitmq"

	"github.com/gin-gonic/gin"
)

func SetRabbitMQ(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		rabbitmqService := appCtx.GetRabbitMQService()
		// create topic if not exist
		go func() {
			defer common.Recovery()
			rabbitmqService.GetQueue(context.Background(), common.IndexingQueue)
		}()
		ctx := rabbitmq.WithContext(ginCtx.Request.Context(), rabbitmqService)
		ginCtx.Request = ginCtx.Request.WithContext(ctx)
		ginCtx.Next()
	}
}
