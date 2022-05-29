package main

import (
	"log"
	"user_management/common"
	"user_management/components/appctx"
	"user_management/components/elasticsearch"
	"user_management/components/rabbitmq"

	es "github.com/elastic/go-elasticsearch/v7"
)

func main() {
	config := &appctx.Config{}
	appctx.GetConfig(config)

	configEs := &es.Config{
		Addresses: []string{config.ElasticSearch.Host},
		Username:  config.ElasticSearch.Username,
		Password:  config.ElasticSearch.Password,
	}
	esService, esErr := elasticsearch.NewEsService(*configEs)
	if esErr != nil {
		return
	}
	esService.LogInfo()

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

	queue, queueErr := rabbitmqService.GetQueue(common.IndexingQueue)
	if queueErr != nil {
		log.Println("Get Queue Failled: ", queueErr)
		return
	}
	msgs, consumErr := rabbitmqService.Consume(queue)
	if consumErr != nil {
		log.Println("Consume Queue Failled: ", consumErr)
		return
	}
	forever := make(chan bool)

	go func() {
		defer common.Recovery()
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
