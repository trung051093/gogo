package main

import (
	"fmt"
	"log"
	"user_management/components/appctx"
	"user_management/components/elasticsearch"
	"user_management/components/rabbitmq"
	"user_management/middleware"
	"user_management/modules/auth"
	"user_management/modules/user"
	usermodel "user_management/modules/user/model"

	esv7 "github.com/elastic/go-elasticsearch/v7"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	var config = &appctx.Config{}
	appctx.GetConfig(config)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.Database.Host,
		config.Database.Username,
		config.Database.Password,
		config.Database.Name,
		config.Database.Port,
		config.Database.SSLMode,
		config.Database.TimeZone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Println("Connect Database Error: ", err)
		return
	}

	db = db.Debug()

	validate := validator.New()

	db.AutoMigrate(&usermodel.User{})

	configEs := &esv7.Config{
		Addresses: []string{config.ElasticSearch.Host},
		Username:  config.ElasticSearch.Username,
		Password:  config.ElasticSearch.Password,
	}
	esService, esErr := elasticsearch.NewEsService(*configEs)
	if esErr != nil {
		return
	}

	configRabbitMQ := &rabbitmq.RabbitmqConfig{
		Host: config.RabbitMQ.Host,
		Port: config.RabbitMQ.Port,
		User: config.RabbitMQ.Username,
		Pass: config.RabbitMQ.Password,
	}
	rabbitmqService, rabbitErr := rabbitmq.NewRabbitMQ(*configRabbitMQ)
	if rabbitErr != nil {
		return
	}
	defer rabbitmqService.Close()

	appCtx := appctx.NewAppContext(db, validate, config, esService, rabbitmqService)

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

		// authentication
		v1.POST("/auth/register", auth.RegisterUserHandler(appCtx))
		v1.POST("/auth/login", auth.LoginUserHandler(appCtx))
	}
	router.Run(fmt.Sprintf(":%d", config.Server.Port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
