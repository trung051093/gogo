package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/streadway/amqp"
)

type RabbitmqConfig struct {
	Host string
	Port int
	User string
	Pass string
}

type RabbitmqSerivce struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

type key string

var RabbitMQServiceKey key = "RabbitMQService"
var once sync.Once
var instance *RabbitmqSerivce
var connErr error
var channelErr error

func NewRabbitMQ(config RabbitmqConfig) (*RabbitmqSerivce, error) {
	once.Do(func() {
		connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/", config.User, config.Pass, config.Host, config.Port)
		conn, connErr := amqp.Dial(connStr)
		if connErr != nil {
			log.Println("Failed to connect to RabbitMQ: ", connErr)
		}
		channel, channelErr := conn.Channel()
		if channelErr != nil {
			log.Println("Failed to open the channel RabbitMQ: ", channelErr)
		}
		instance = &RabbitmqSerivce{
			conn:    conn,
			channel: channel,
		}
	})
	if connErr != nil {
		return nil, connErr
	}
	if channelErr != nil {
		return nil, channelErr
	}
	return instance, nil
}

func WithContext(ctx context.Context, rabbitmq *RabbitmqSerivce) context.Context {
	return context.WithValue(ctx, RabbitMQServiceKey, rabbitmq)
}

func FromContext(ctx context.Context) (*RabbitmqSerivce, bool) {
	rabbitmqService := ctx.Value(RabbitMQServiceKey)
	if es, ok := rabbitmqService.(*RabbitmqSerivce); ok {
		return es, true
	}
	return nil, false
}

func (r *RabbitmqSerivce) GetQueue(topic string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		topic, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func (r *RabbitmqSerivce) Publish(queue amqp.Queue, message string) error {
	return r.channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (r *RabbitmqSerivce) Consume(q amqp.Queue) (<-chan amqp.Delivery, error) {
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

func (r *RabbitmqSerivce) Close() {
	r.channel.Close()
	r.conn.Close()
}
