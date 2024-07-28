package middleware

import (
	"errors"
	"gogo/common"
	"gogo/components/appctx"
	"gogo/modules/tracker/entity"

	"github.com/gin-gonic/gin"
)

func AdminRequired(appCtx appctx.AppContext) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		tokenPayload := ginCtx.MustGet(common.CurrentAuth).(*entity.TokenPayload)
		common.PanicIf(tokenPayload == nil || tokenPayload.Role != entity.AdminRole.String(), common.ErrorUnauthorized(errors.New("unauthorized")))
		ginCtx.Next()
	}
}
