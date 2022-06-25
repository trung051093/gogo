package appctx

import (
	"context"
	"time"
	cacheprovider "user_management/components/cache"
	esprovider "user_management/components/elasticsearch"
	jaegerprovider "user_management/components/jaeger"
	"user_management/components/mailer"
	rabbitmqprovider "user_management/components/rabbitmq"
	redisprovider "user_management/components/redis"
	socketprovider "user_management/components/socketio"
	storageprovider "user_management/components/storage"

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
	GetStorageService() storageprovider.StorageService
	GetSocketService() socketprovider.SocketService
	GetCacheService() cacheprovider.CacheService
	GetMailService() mailer.MailService
}

type appContext struct {
	config          *Config
	db              *gorm.DB
	validate        *validator.Validate
	esService       esprovider.ElasticSearchSevice
	rabbitmqService rabbitmqprovider.RabbitmqSerivce
	redisService    redisprovider.RedisService
	storageService  storageprovider.StorageService
	socketService   socketprovider.SocketService
	jaegerService   jaegerprovider.JaegerService
	cacheService    cacheprovider.CacheService
	mailService     mailer.MailService
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
	storageService storageprovider.StorageService,
	socketService socketprovider.SocketService,
	jaegerService jaegerprovider.JaegerService,
	cacheService cacheprovider.CacheService,
	mailService mailer.MailService,
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
		cacheService:    cacheService,
		mailService:     mailService,
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

func (appCtx *appContext) GetStorageService() storageprovider.StorageService {
	return appCtx.storageService
}

func (appCtx *appContext) GetSocketService() socketprovider.SocketService {
	return appCtx.socketService
}

func (appCtx *appContext) GetJaegerService() jaegerprovider.JaegerService {
	return appCtx.jaegerService
}

func (appCtx *appContext) GetCacheService() cacheprovider.CacheService {
	return appCtx.cacheService
}

func (appCtx *appContext) GetMailService() mailer.MailService {
	return appCtx.mailService
}
