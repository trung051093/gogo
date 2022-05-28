package main

import (
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
	esService := elasticsearch.NewEsService(*configEs)
	esService.LogInfo()

	configRabbitMQ := &rabbitmq.RabbitmqConfig{
		Host: config.RabbitMQ.Host,
		Port: config.RabbitMQ.Port,
		User: config.RabbitMQ.Username,
		Pass: config.RabbitMQ.Password,
	}
	rabbitmqService := rabbitmq.NewRabbitMQ(*configRabbitMQ)
	defer rabbitmqService.Close()
}
