package common

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AppContext interface {
	GetValidator() *validator.Validate
}

func ParseRequest[D any](appCtx AppContext, ginCtx *gin.Context) func(*D) error {
	return func(dto *D) error {
		if err := ginCtx.ShouldBindUri(dto); err != nil {
			return err
		}
		if err := ginCtx.ShouldBind(dto); err != nil {
			return err
		}
		if err := appCtx.GetValidator().Struct(dto); err != nil {
			return err
		}
		return nil
	}
}
