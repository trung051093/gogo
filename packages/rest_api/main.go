package main

import (
	"context"
	"fmt"
	"log"
	"user_management/common"
	"user_management/components/appctx"
	"user_management/components/dbprovider"
	esprovider "user_management/components/elasticsearch"
	rabbitmqprovider "user_management/components/rabbitmq"
	redisprovider "user_management/components/redis"
	socketprovider "user_management/components/socketio"
	"user_management/components/storage"
	"user_management/middleware"
	"user_management/modules/indexer"
	"user_management/modules/notificator"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func main() {
	config := appctx.GetConfig()
	dbprovider, err := dbprovider.NewDBProvider(
		&dbprovider.DBConfig{
			Host:     config.Database.Host,
			Username: config.Database.Username,
			Password: config.Database.Password,
			Name:     config.Database.Name,
			Port:     config.Database.Port,
			SSLMode:  config.Database.SSLMode,
			TimeZone: config.Database.TimeZone,
		},
		// dbprovider.WithDebug,
		// dbprovider.WithAutoMigration(&usermodel.User{}),
	)

	if err != nil {
		log.Println("Connect Database Error: ", err)
		return
	}

	validate := validator.New()

	configEs := config.GetElasticSearchConfig()
	esService, esErr := esprovider.NewEsService(configEs)
	if esErr != nil {
		return
	}

	configRabbitMQ := config.GetRabbitMQConfig()
	rabbitmqService, rabbitErr := rabbitmqprovider.NewRabbitMQ(configRabbitMQ)
	if rabbitErr != nil {
		return
	}
	configRedis := config.GetRedisConfig()
	redisProvider := redisprovider.NewRedisService(configRedis)
	configStorage := config.GetStorageConfig()

	storageService, storageErr := storage.NewStorage(configStorage)
	if storageErr != nil {
		return
	}
	createBucketErr := storageService.CreateBucket(context.Background(), common.PhotoBucket, common.PhotoBucketRegion)
	if createBucketErr != nil {
		log.Println("Create bucket error: ", createBucketErr)
	}

	socketService := socketprovider.NewSocketProvider(
		socketprovider.WithRedisAdapter(&socketprovider.SocketRedisAdapterConfig{
			Addr:     configRedis.Addr,
			Password: configRedis.Password,
			Prefix:   "socketio",
		}),
		socketprovider.WithWebsocketTransport,
	)
	appCtx := appctx.NewAppContext(
		dbprovider.GetDBConnection(),
		validate,
		config,
		esService,
		rabbitmqService,
		redisProvider,
		storageService,
		socketService,
	)

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))
	router.Use(middleware.ErrorHandler(appCtx))
	router.Use(middleware.SetElasticSearch(appCtx))
	router.Use(middleware.SetRabbitMQ(appCtx))
	router.Use(middleware.SetSocketIO(appCtx))

	// handler background
	go notificator.FileHandler(appCtx)
	go notificator.Handler(appCtx)
	go indexer.Handler(appCtx)

	mainRoutes(appCtx, router)
	socketRoutes(appCtx, router)
	router.Run(fmt.Sprintf(":%d", config.Server.Port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
