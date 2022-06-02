package rabbitmqprovider

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/wagslane/go-rabbitmq"
)

type RabbitmqSerivce interface {
	Consuming(ctx context.Context, handler rabbitmq.Handler, queue string, routingKeys []string, optionFuncs ...func(*rabbitmq.ConsumeOptions)) error
	Publish(ctx context.Context, data []byte, routingKeys []string, optionFuncs ...func(*rabbitmq.PublishOptions)) error
	Close()
}

type RabbitmqConfig struct {
	Host string
	Port int
	User string
	Pass string
}

type rabbitmqSerivce struct {
	consumer  rabbitmq.Consumer
	publisher *rabbitmq.Publisher
}

type key string

var RabbitMQServiceKey key = "RabbitMQService"
var once sync.Once
var instance *rabbitmqSerivce
var instanceErr error

func NewRabbitMQ(config *RabbitmqConfig) (*rabbitmqSerivce, error) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/", config.User, config.Pass, config.Host, config.Port)
	consumer, err := rabbitmq.NewConsumer(
		connStr,
		rabbitmq.Config{},
		rabbitmq.WithConsumerOptionsLogging,
		rabbitmq.WithConsumerOptionsReconnectInterval(time.Second),
	)
	if err != nil {
		log.Println("Cannot create consumer: ", err)
		return nil, err
	}

	publisher, err := rabbitmq.NewPublisher(
		connStr,
		rabbitmq.Config{},
		rabbitmq.WithPublisherOptionsLogging,
	)
	if err != nil {
		log.Println("Cannot create publisher: ", err)
		return nil, err
	}

	returns := publisher.NotifyReturn()
	go func() {
		for r := range returns {
			log.Printf("message returned from server: %s", string(r.Body))
		}
	}()

	confirmations := publisher.NotifyPublish()
	go func() {
		for c := range confirmations {
			log.Printf("message confirmed from server. tag: %v, ack: %v", c.DeliveryTag, c.Ack)
		}
	}()

	return &rabbitmqSerivce{
		consumer:  consumer,
		publisher: publisher,
	}, nil
}

// singleton
func GetIntance(config *RabbitmqConfig) (*rabbitmqSerivce, error) {
	once.Do(func() {
		service, instanceErr := NewRabbitMQ(config)
		if instanceErr != nil {
			return
		}
		instanceErr = nil
		instance = service
	})
	return instance, instanceErr
}

func WithContext(ctx context.Context, rabbitmq RabbitmqSerivce) context.Context {
	return context.WithValue(ctx, RabbitMQServiceKey, rabbitmq)
}

func FromContext(ctx context.Context) (*rabbitmqSerivce, bool) {
	rabbitmqService := ctx.Value(RabbitMQServiceKey)
	if es, ok := rabbitmqService.(*rabbitmqSerivce); ok {
		return es, true
	}
	return nil, false
}

func (r *rabbitmqSerivce) Consuming(ctx context.Context, handler rabbitmq.Handler, queue string, routingKeys []string, optionFuncs ...func(*rabbitmq.ConsumeOptions)) error {
	return r.consumer.StartConsuming(
		handler,
		queue,
		routingKeys,
		rabbitmq.WithConsumeOptionsConcurrency(10),
	)
}

func (r *rabbitmqSerivce) Publish(ctx context.Context, data []byte, routingKeys []string, optionFuncs ...func(*rabbitmq.PublishOptions)) error {
	return r.publisher.Publish(
		data,
		routingKeys,
		optionFuncs...,
	)
}

func (r *rabbitmqSerivce) Close() {
	r.consumer.Close()
	r.publisher.Close()
}
