package main

import (
	"context"
	"fmt"
	"gogo/components/appctx"
	dbprovider "gogo/components/dbprovider"
	rabbitmqprovider "gogo/components/rabbitmq"
	randomuserapi "gogo/components/randomuserapi"
	usermodel "gogo/modules/user/model"
	"log"
	"sync"
	"time"

	"gogo/modules/user"
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
		dbprovider.WithDebug,
		dbprovider.WithAutoMigration(&usermodel.User{}),
	)

	if err != nil {
		log.Println("Connect Database Error: ", err)
		return
	}

	configRabbitMQ := config.GetRabbitMQConfig()
	rabbitmqService, rabbitErr := rabbitmqprovider.NewRabbitMQ(configRabbitMQ)
	if rabbitErr != nil {
		return
	}
	defer rabbitmqService.Close()

	repository := user.NewUserRepository(dbprovider.GetDBConnection())

	log.Println(repository)

	c := &randomuserapi.QueryConfig{
		MaxResults: 5000,
		Gender:     "",
		Password:   "",
		Seed:       "",
	}
	users, err := randomuserapi.Generate(c)
	log.Println("Number random", len(users))

	if err != nil {
		fmt.Println(err)
		return
	}

	ctx := rabbitmqprovider.WithContext(context.Background(), rabbitmqService)
	createCh := make(chan usermodel.UserCreateDto, 10)
	done := make(chan int)
	var wg sync.WaitGroup

	go func() {
		for _, user := range users {
			wg.Add(1)
			fmt.Println(user.Name.Title + " : " + user.Name.First + " " + user.Name.Last)
			createCh <- usermodel.UserCreateDto{
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
		case userCreateDto := <-createCh:
			go func(u usermodel.UserCreateDto) {
				repository.Create(ctx, u.ToEntity())
				wg.Done()
			}(userCreateDto)
		case <-done:
			wg.Wait()
			println("Done !!!")
			return
		}
	}
}
