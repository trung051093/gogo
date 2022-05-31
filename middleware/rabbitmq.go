package middleware

import (
	"user_management/components/appctx"
	rabbitmqprovider "user_management/components/rabbitmq"

	"github.com/gin-gonic/gin"
)

func SetRabbitMQ(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		rabbitmqService := appCtx.GetRabbitMQService()
		ctx := rabbitmqprovider.WithContext(ginCtx.Request.Context(), rabbitmqService)
		ginCtx.Request = ginCtx.Request.WithContext(ctx)
		ginCtx.Next()
	}
}
