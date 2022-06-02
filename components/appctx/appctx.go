package appctx

import (
	esprovider "user_management/components/elasticsearch"
	rabbitmqprovider "user_management/components/rabbitmq"
	"user_management/components/redisprovider"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	GetValidator() *validator.Validate
	GetConfig() *Config
	GetESService() esprovider.ElasticSearchSevice
	GetRabbitMQService() rabbitmqprovider.RabbitmqSerivce
	GetRedisService() redisprovider.RedisService
}

type appContext struct {
	config          *Config
	db              *gorm.DB
	validate        *validator.Validate
	esService       esprovider.ElasticSearchSevice
	rabbitmqService rabbitmqprovider.RabbitmqSerivce
	redisService    redisprovider.RedisService
}

func NewAppContext(
	db *gorm.DB,
	validate *validator.Validate,
	config *Config,
	esService esprovider.ElasticSearchSevice,
	rabbitmqService rabbitmqprovider.RabbitmqSerivce,
	redisService redisprovider.RedisService,
) *appContext {
	return &appContext{
		db:              db,
		validate:        validate,
		config:          config,
		esService:       esService,
		rabbitmqService: rabbitmqService,
		redisService:    redisService,
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

func (appCtx *appContext) GetESService() esprovider.ElasticSearchSevice {
	return appCtx.esService
}

func (appCtx *appContext) GetRabbitMQService() rabbitmqprovider.RabbitmqSerivce {
	return appCtx.rabbitmqService
}

func (appCtx *appContext) GetRedisService() redisprovider.RedisService {
	return appCtx.redisService
}
