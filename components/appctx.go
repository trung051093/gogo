package component

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	GetValidator() *validator.Validate
	GetConfig() *Config
}

type appContext struct {
	config   *Config
	db       *gorm.DB
	validate *validator.Validate
}

func NewAppContext(
	db *gorm.DB,
	validate *validator.Validate,
	config *Config,
) *appContext {
	return &appContext{db: db, validate: validate, config: config}
}

func (appCtx *appContext) GetMainDBConnection() *gorm.DB {
	return appCtx.db
}

func (appCtx *appContext) GetValidator() *validator.Validate {
	return appCtx.validate
}

func (appCtx *appContext) GetConfig() *Config {
	return appCtx.config
}
