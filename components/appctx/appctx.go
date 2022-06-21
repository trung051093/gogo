package appctx

import (
	"context"
	"time"
	esprovider "user_management/components/elasticsearch"
	jaegerprovider "user_management/components/jaeger"
	rabbitmqprovider "user_management/components/rabbitmq"
	redisprovider "user_management/components/redis"
	socketprovider "user_management/components/socketio"
	"user_management/components/storage"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type AppContext interface {
	GetMainDBConnection() *gorm.DB
	GetValidator() *validator.Validate
	GetConfig() *Config
	GetESService() esprovider.ElasticSearchSevice
	GetRabbitMQService() rabbitmqprovider.RabbitmqSerivce
	GetRedisService() redisprovider.RedisService
	GetStorageService() storage.StorageService
	GetSocketService() socketprovider.SocketService
}

type appContext struct {
	config          *Config
	db              *gorm.DB
	validate        *validator.Validate
	esService       esprovider.ElasticSearchSevice
	rabbitmqService rabbitmqprovider.RabbitmqSerivce
	redisService    redisprovider.RedisService
	storageService  storage.StorageService
	socketService   socketprovider.SocketService
	jaegerService   jaegerprovider.JaegerService
}

type key string

var AppContextKey key = "AppContextKey"

func NewAppContext(
	db *gorm.DB,
	validate *validator.Validate,
	config *Config,
	esService esprovider.ElasticSearchSevice,
	rabbitmqService rabbitmqprovider.RabbitmqSerivce,
	redisService redisprovider.RedisService,
	storageService storage.StorageService,
	socketService socketprovider.SocketService,
	jaegerService jaegerprovider.JaegerService,
) *appContext {
	return &appContext{
		db:              db,
		validate:        validate,
		config:          config,
		esService:       esService,
		rabbitmqService: rabbitmqService,
		redisService:    redisService,
		storageService:  storageService,
		socketService:   socketService,
		jaegerService:   jaegerService,
	}
}

func WithContext(ctx context.Context, appCtx AppContext) context.Context {
	return context.WithValue(ctx, AppContextKey, appCtx)
}

func FromContext(ctx context.Context) (*appContext, bool) {
	appCtx := ctx.Value(AppContextKey)
	if c, ok := appCtx.(*appContext); ok {
		return c, true
	}
	return nil, false
}

func (appCtx *appContext) GetMainDBConnection() *gorm.DB {
	return appCtx.db.Session(&gorm.Session{
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		NewDB:  true,
		Logger: appCtx.db.Logger.LogMode(logger.Error),
	})
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

func (appCtx *appContext) GetStorageService() storage.StorageService {
	return appCtx.storageService
}

func (appCtx *appContext) GetSocketService() socketprovider.SocketService {
	return appCtx.socketService
}

func (appCtx *appContext) GetJaegerService() jaegerprovider.JaegerService {
	return appCtx.jaegerService
}
