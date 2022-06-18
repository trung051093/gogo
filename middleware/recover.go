package middleware

import (
	"user_management/common"
	"user_management/components/appctx"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func ErrorHandler(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				ginCtx.Header("Content-Type", "application/json")

				if appErr, ok := err.(*common.AppError); ok {
					ginCtx.AbortWithStatusJSON(appErr.StatusCode, appErr)
					log.Error(appErr)
					panic(appErr)
				}

				appErr := common.ErrorInternal(err.(error))
				ginCtx.AbortWithStatusJSON(appErr.StatusCode, appErr)
				log.Error(appErr)
				panic(appErr)
			}
		}()

		ginCtx.Next()
	}
}
