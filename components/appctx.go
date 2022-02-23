package component

import "gorm.io/gorm"

type AppContext interface {
	GetMainDBConnection() (db *gorm.DB)
}

type appContext struct {
	db *gorm.DB
}

func NewAppContext(db *gorm.DB) *appContext {
	return &appContext{db: db}
}

func (appCtx *appContext) GetMainDBConnection() (db *gorm.DB) {
	return appCtx.db
}
