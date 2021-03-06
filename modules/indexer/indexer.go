package indexer

import (
	"context"
	"gogo/common"
	"gogo/components/appctx"
	"log"

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
			var dataIndex common.DataIndex
			if err := dataIndex.FromByte(d.Body); err != nil {
				log.Println("Error message: ", err)
				return rabbitmq.NackRequeue
			}
			dataByte, err := dataIndex.GetByte()
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
