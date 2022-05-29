package main

import (
	"context"
	"fmt"
	"log"
	"user_management/components/appctx"
	usermodel "user_management/modules/user/model"

	"user_management/modules/user"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	config := appctx.GetConfig()
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

	log.Println(db)

	db = db.Debug()

	db.AutoMigrate(&usermodel.User{})

	repository := user.NewUserRepository(db)

	log.Println(repository)

	c := &QueryConfig{
		5000, "", "", "",
	}
	users, err := Generate(c)

	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := context.Background()

	for _, user := range users {
		fmt.Println(user.Name.Title + " : " + user.Name.First + " " + user.Name.Last)
		repository.Create(ctx, &usermodel.UserCreate{
			FirstName:   user.Name.First,
			LastName:    user.Name.Last,
			Email:       user.Email,
			PhoneNumber: user.Phone,
			Gender:      user.Gender,
			Address: fmt.Sprintf("%d %s %s %s",
				user.Location.Street.Number,
				user.Location.Street.Name,
				user.Location.City,
				user.Location.State),
		})
	}
}
