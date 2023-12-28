package common

import (
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type AppContext interface {
	GetValidator() *validator.Validate
}

func ParseRequest[D any](appCtx AppContext, ginCtx *gin.Context) func(*D) error {
	return func(dto *D) (err error) {
		if err = ginCtx.ShouldBind(dto); err != nil {
			return err
		}

		validator := appCtx.GetValidator()
		if err = validator.Struct(dto); err != nil {
			return err
		}
		return nil
	}
}
