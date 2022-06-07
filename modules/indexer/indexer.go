package indexer

import (
	"context"
	"log"
	"user_management/common"
	"user_management/components/appctx"

	"github.com/wagslane/go-rabbitmq"
)

func Handler(appctx appctx.AppContext) {
	defer common.Recovery()
	ctx := context.Background()

	esService := appctx.GetESService()
	rabbitmqService := appctx.GetRabbitMQService()
	esService.LogInfo(ctx)

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
}
