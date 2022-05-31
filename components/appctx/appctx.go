package appctx

import (
	"user_management/components/elasticsearch"
	rabbitmqprovider "user_management/components/rabbitmq"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	GetValidator() *validator.Validate
	GetConfig() *Config
	GetESService() elasticsearch.ElasticSearchSevice
	GetRabbitMQService() rabbitmqprovider.RabbitmqSerivce
}

type appContext struct {
	config          *Config
	db              *gorm.DB
	validate        *validator.Validate
	esService       elasticsearch.ElasticSearchSevice
	rabbitmqService rabbitmqprovider.RabbitmqSerivce
}

func NewAppContext(
	db *gorm.DB,
	validate *validator.Validate,
	config *Config,
	esService elasticsearch.ElasticSearchSevice,
	rabbitmqService rabbitmqprovider.RabbitmqSerivce,
) *appContext {
	return &appContext{
		db:              db,
		validate:        validate,
		config:          config,
		esService:       esService,
		rabbitmqService: rabbitmqService,
	}
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

func (appCtx *appContext) GetESService() elasticsearch.ElasticSearchSevice {
	return appCtx.esService
}

func (appCtx *appContext) GetRabbitMQService() rabbitmqprovider.RabbitmqSerivce {
	return appCtx.rabbitmqService
}
