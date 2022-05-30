package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type RabbitmqSerivce interface {
	GetQueue(ctx context.Context, topic string) (amqp.Queue, error)
	PublishWithTopic(ctx context.Context, topic string, data interface{}) error
	Consume(ctx context.Context, q amqp.Queue) (<-chan amqp.Delivery, error)
	QueuePurge(ctx context.Context, topic string) (int, error)
	Close()
}

type RabbitmqConfig struct {
	Host string
	Port int
	User string
	Pass string
}

type rabbitmqSerivce struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type key string

var RabbitMQServiceKey key = "RabbitMQService"
var once sync.Once
var instance *rabbitmqSerivce
var instanceErr error

func NewRabbitMQ(config RabbitmqConfig) (*rabbitmqSerivce, error) {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/", config.User, config.Pass, config.Host, config.Port)
	conn, connErr := amqp.Dial(connStr)
	if connErr != nil {
		log.Println("Failed to connect to RabbitMQ: ", connErr)
		return nil, connErr
	}

	channel, channelErr := conn.Channel()
	if channelErr != nil {
		log.Println("Failed to open the channel RabbitMQ: ", channelErr)
		return nil, channelErr
	}

	return &rabbitmqSerivce{
		conn:    conn,
		channel: channel,
	}, nil
}

// singleton
func GetIntance(config RabbitmqConfig) (*rabbitmqSerivce, error) {
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

func (r *rabbitmqSerivce) GetQueue(ctx context.Context, topic string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		topic, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func (r *rabbitmqSerivce) publish(ctx context.Context, queue amqp.Queue, body []byte) error {
	return r.channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		})
}

func (r *rabbitmqSerivce) PublishWithTopic(ctx context.Context, topic string, data interface{}) error {
	queue, queueErr := r.GetQueue(ctx, topic)
	if queueErr != nil {
		return queueErr
	}
	databyte, dataErr := json.Marshal(data)
	if dataErr != nil {
		return dataErr
	}
	r.publish(ctx, queue, databyte)
	return nil
}

func (r *rabbitmqSerivce) Consume(ctx context.Context, q amqp.Queue) (<-chan amqp.Delivery, error) {
	return r.channel.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}

func (r *rabbitmqSerivce) QueuePurge(ctx context.Context, topic string) (int, error) {
	return r.channel.QueuePurge(topic, false)
}

func (r *rabbitmqSerivce) Close() {
	r.channel.Close()
	r.conn.Close()
}
