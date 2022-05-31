package main

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"
	"user_management/components/appctx"
	rabbitmqprovider "user_management/components/rabbitmq"
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

	configRabbitMQ := config.GetRabbitMQConfig()
	rabbitmqService, rabbitErr := rabbitmqprovider.NewRabbitMQ(*configRabbitMQ)
	if rabbitErr != nil {
		return
	}
	defer rabbitmqService.Close()

	repository := user.NewUserRepository(db)

	log.Println(repository)

	c := &QueryConfig{
		5000, "", "", "",
	}
	users, err := Generate(c)
	log.Println("Number random", len(users))

	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := rabbitmqprovider.WithContext(context.Background(), rabbitmqService)
	createCh := make(chan usermodel.UserCreate, 10)
	done := make(chan int)
	var wg sync.WaitGroup

	go func() {
		for _, user := range users {
			wg.Add(1)
			fmt.Println(user.Name.Title + " : " + user.Name.First + " " + user.Name.Last)
			createCh <- usermodel.UserCreate{
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
			}
			time.Sleep(time.Millisecond * 100)
		}
		done <- 0
	}()

	for {
		select {
		case userCreate := <-createCh:
			go func(u usermodel.UserCreate) {
				repository.Create(ctx, &u)
				wg.Done()
			}(userCreate)
		case <-done:
			wg.Wait()
			println("Done !!!")
			return
		}
	}
}
