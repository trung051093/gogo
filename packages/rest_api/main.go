package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"user_management/common"
	"user_management/components/appctx"
	cacheprovider "user_management/components/cache"
	dbprovider "user_management/components/dbprovider"
	esprovider "user_management/components/elasticsearch"
	jaegerprovider "user_management/components/jaeger"
	graylog "user_management/components/log"
	"user_management/components/mailer"
	rabbitmqprovider "user_management/components/rabbitmq"
	redisprovider "user_management/components/redis"
	socketprovider "user_management/components/socketio"
	storageprovider "user_management/components/storage"
	"user_management/middleware"
	"user_management/modules/indexer"
	"user_management/modules/notificator"
	usermodel "user_management/modules/user/model"
	"user_management/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/trace"
)

func main() {
	config := appctx.GetConfig()

	// graylog
	graylog.Integrate(config.GetGraylogConfig())

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
		dbprovider.WithAutoMigration(&usermodel.User{}),
	)

	if err != nil {
		log.Fatalln("Connect Database Error: ", err)
	}

	validate := validator.New()

	configEs := config.GetElasticSearchConfig()
	esService, esErr := esprovider.NewEsService(configEs)
	if esErr != nil {
		log.Fatalln("Connect Elastic Search Error: ", esErr)
	}

	configRabbitMQ := config.GetRabbitMQConfig()
	rabbitmqService, rabbitErr := rabbitmqprovider.NewRabbitMQ(configRabbitMQ)
	if rabbitErr != nil {
		log.Fatalln("Connect RabbitMQ Error: ", rabbitErr)
	}
	configRedis := config.GetRedisConfig()
	redisService := redisprovider.NewRedisService(configRedis)

	configStorage := config.GetStorageConfig()
	storageService, storageErr := storageprovider.NewStorage(configStorage)
	if storageErr != nil {
		log.Fatalln("Connect Minio Error: ", storageErr)
	}
	createBucketErr := storageService.CreateBucket(context.Background(), common.ImageBucket, common.ImageBucketRegion)
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
	jaegerService := jaegerprovider.NewExporter(config.GetJaegerConfig())
	cacheService := cacheprovider.NewCacheService(redisService.GetClient())
	mailService := mailer.NewMailer(config.GetMailConfig())

	appCtx := appctx.NewAppContext(
		dbprovider.GetDBConnection(),
		validate,
		config,
		esService,
		rabbitmqService,
		redisService,
		storageService,
		socketService,
		jaegerService,
		cacheService,
		mailService,
	)

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))
	router.Use(middleware.ErrorHandler(appCtx))
	router.Use(middleware.SetAppContextIntoRequest(appCtx))

	// handler background
	// go notificator.FileHandler(appCtx)
	go notificator.Handler(appCtx)
	go indexer.Handler(appCtx)

	routes.MainRoutes(appCtx, router)
	routes.SocketRoutes(appCtx, router)
	routes.SwaggerRoutes(appCtx, router)

	// And now finally register it as a Trace Exporter
	trace.RegisterExporter(jaegerService.GetExporter())
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.ProbabilitySampler(1)})

	// listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	http.ListenAndServe(fmt.Sprintf(":%d", config.Server.Port), &ochttp.Handler{Handler: router})
}
