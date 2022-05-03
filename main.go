package main

import (
	"fmt"
	"log"
	"user_management/components/appctx"
	"user_management/middleware"
	"user_management/modules/auth"
	"user_management/modules/user"
	usermodel "user_management/modules/user/model"

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

	appCtx := appctx.NewAppContext(db, validate, config)

	router := gin.Default()
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	router.Use(cors.New(corsConfig))
	router.Use(middleware.ErrorHandler(appCtx))
	v1 := router.Group("/v1")
	{
		v1.POST("/user", user.CreateUserHandler(appCtx))
		v1.PATCH("/user/:id", user.UpdateUserHandler(appCtx))
		v1.DELETE("/user/:id", user.DeleteUserHandler(appCtx))
		v1.GET("/user/:id", user.GetUserHandler(appCtx))
		v1.GET("/users", middleware.JWTRequireHandler(appCtx), user.ListUserHandler(appCtx))

		v1.POST("/auth/register", auth.RegisterUserHandler(appCtx))
		v1.POST("/auth/login", auth.LoginUserHandler(appCtx))
	}
	router.Run(fmt.Sprintf(":%d", config.Server.Port)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
