package main

import (
	"fmt"
	"log"
	component "user_management/components"
	"user_management/config"
	"user_management/modules/user"
	usermodel "user_management/modules/user/model"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		config.DB_HOST,
		config.DB_USER,
		config.DB_PASSWORD,
		config.DB_DATABASE,
		config.DB_PORT,
		config.SSL_MODE,
		config.TIME_ZONE)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Println("Connect Database Error: ", err)
		return
	}

	log.Println(db)

	db = db.Debug()
	validate := validator.New()

	db.AutoMigrate(&usermodel.User{})

	appCtx := component.NewAppContext(db, validate)

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		v1.POST("/user", user.CreateUserHandler(appCtx))
		v1.PATCH("/user/:id", user.UpdateUserHandler(appCtx))
		v1.DELETE("/user/:id", user.DeleteUserHandler(appCtx))
		v1.GET("/user/:id", user.GetUserHandler(appCtx))
		v1.GET("/users", user.ListUserHandler(appCtx))
	}
	router.Run(fmt.Sprintf(":%d", config.API_PORT)) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
