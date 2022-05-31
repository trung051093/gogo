package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"user_management/common"
	"user_management/components/appctx"
	"user_management/components/elasticsearch"
	rabbitmqprovider "user_management/components/rabbitmq"

	"github.com/wagslane/go-rabbitmq"
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
	rabbitmqService, rabbitErr := rabbitmqprovider.NewRabbitMQ(*configRabbitMQ)
	if rabbitErr != nil {
		return
	}

	consumeErr := rabbitmqService.Consuming(
		ctx,
		func(d rabbitmq.Delivery) rabbitmq.Action {
			log.Printf("Received a message: %v", string(d.Body))
			dataIndex, dataByte, err := common.MessageToDataIndex(d.Body)
			if err != nil {
				log.Println("Error message: ", err)
				return rabbitmq.NackRequeue
			}
			switch dataIndex.Action {
			case common.Create:
				esService.Index(ctx, dataIndex.Index, dataIndex.Id, dataByte)
			case common.Update:
				esService.Index(ctx, dataIndex.Index, dataIndex.Id, dataByte)
			case common.Delete:
				esService.Delete(ctx, dataIndex.Index, dataIndex.Id)
			}
			// rabbitmq.Ack, rabbitmq.NackDiscard, rabbitmq.NackRequeue
			return rabbitmq.Ack
		},
		common.IndexingQueue,
		[]string{},
	)
	if consumeErr != nil {
		log.Println("Consuming error: ", consumeErr)
	}
	// block main thread - wait for shutdown signal
	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		fmt.Println()
		fmt.Println(sig)
		done <- true
	}()

	fmt.Println("awaiting signal")
	<-done
	fmt.Println("stopping consumer")
}
