package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"user_management/common"
	"user_management/components/appctx"
	"user_management/components/dbprovider"
	esprovider "user_management/components/elasticsearch"
	rabbitmqprovider "user_management/components/rabbitmq"
	"user_management/components/redisprovider"
	"user_management/components/storage"
	cachedecorator "user_management/decorators/cache"
	"user_management/middleware"
	"user_management/modules/auth"
	"user_management/modules/file"
	"user_management/modules/user"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/minio/minio-go/v7/pkg/notification"
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
	go storageService.ListenNotification(
		context.Background(),
		common.PhotoBucket,
		"",
		"",
		[]string{
			"s3:ObjectCreated:*",
			"s3:ObjectRemoved:*",
		},
		func(noti *notification.Info) {
			log.Println("minio notification info:", noti)
		},
	)
	createBucketErr := storageService.CreateBucket(context.Background(), common.PhotoBucket, common.PhotoBucketRegion)
	if createBucketErr != nil {
		log.Println("Create bucket error: ", createBucketErr)
	}
	appCtx := appctx.NewAppContext(
		dbprovider.GetDBConnection(),
		validate,
		config,
		esService,
		rabbitmqService,
		redisProvider,
		storageService,
	)

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))
	router.Use(middleware.ErrorHandler(appCtx))
	router.Use(middleware.SetElasticSearch(appCtx))
	router.Use(middleware.SetRabbitMQ(appCtx))
	v1 := router.Group("/api/v1")
	{
		// user
		v1.POST("/user", user.CreateUserHandler(appCtx))
		v1.PATCH("/user/:id", user.UpdateUserHandler(appCtx))
		v1.DELETE("/user/:id", user.DeleteUserHandler(appCtx))
		v1.GET("/user/:id", user.GetUserHandler(appCtx))
		v1.GET("/users", user.ListUserHandler(appCtx))
		v1.GET("/user/search", user.SearchUserHandler(appCtx))
		// cache request
		v1.GET("/users-cache", cachedecorator.CacheRequest(appCtx, "user", 15*time.Minute, user.ListUserHandler))

		// authentication
		v1.POST("/auth/register", auth.RegisterUserHandler(appCtx))
		v1.POST("/auth/login", auth.LoginUserHandler(appCtx))

		// photo
		v1.GET("/file/presign-url", file.GetUploadPresignedUrl(appCtx))
	}
	router.Run(fmt.Sprintf(":%d", config.Server.Port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
