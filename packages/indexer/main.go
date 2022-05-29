package main

import (
	"context"
	"log"
	"sync"
	"user_management/common"
	"user_management/components/appctx"
	"user_management/components/elasticsearch"
	"user_management/components/rabbitmq"

	"github.com/streadway/amqp"
)

func main() {
	config := appctx.GetConfig()
	ctx := context.Background()

	configEs := config.GetElasticSearchConfig()
	esService, esErr := elasticsearch.NewEsService(*configEs)
	if esErr != nil {
		return
	}
	esService.LogInfo(ctx)

	configRabbitMQ := config.GetRabbitMQConfig()
	rabbitmqService, rabbitErr := rabbitmq.NewRabbitMQ(*configRabbitMQ)
	if rabbitErr != nil {
		return
	}
	defer rabbitmqService.Close()

	queue, queueErr := rabbitmqService.GetQueue(ctx, common.IndexingQueue)
	if queueErr != nil {
		log.Println("Get Queue Failled: ", queueErr)
		return
	}
	// remove all message in queue
	rabbitmqService.QueuePurge(ctx, common.IndexingQueue)

	msgs, consumErr := rabbitmqService.Consume(ctx, queue)
	if consumErr != nil {
		log.Println("Consume Queue Failled: ", consumErr)
		return
	}
	forever := make(chan bool)
	var wg sync.WaitGroup

	for msg := range msgs {
		wg.Add(1)
		go func(message amqp.Delivery) {
			defer common.Recovery()
			defer wg.Done()
			defer message.Ack(true)
			log.Printf("Received a message: %s", message.Body)
			dataIndex, dataByte, err := common.MessageToDataIndex(message.Body)
			if err != nil {
				log.Println("Error message: ", err)
				return
			}
			switch dataIndex.Action {
			case common.Create:
				esService.Index(ctx, dataIndex.Index, dataIndex.Id, dataByte)
			case common.Update:
				esService.Index(ctx, dataIndex.Index, dataIndex.Id, dataByte)
			case common.Delete:
				esService.Delete(ctx, dataIndex.Index, dataIndex.Id)
			}
		}(msg)
	}
	wg.Wait()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
