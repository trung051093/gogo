package component

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AppContext interface {
	GetMainDBConnection() (db *gorm.DB)
	GetValidator() *validator.Validate
}

type appContext struct {
	db       *gorm.DB
	validate *validator.Validate
}

func NewAppContext(db *gorm.DB, validate *validator.Validate) *appContext {
	return &appContext{db: db, validate: validate}
}

func (appCtx *appContext) GetMainDBConnection() (db *gorm.DB) {
	return appCtx.db
}

func (appCtx *appContext) GetValidator() *validator.Validate {
	return appCtx.validate
}
