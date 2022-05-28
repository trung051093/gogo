package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

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

func NewRabbitMQ(config RabbitmqConfig) *rabbitmqSerivce {
	connStr := fmt.Sprintf("amqp://%s:%s@%s:%d/", config.User, config.Pass, config.Host, config.Port)
	conn, connErr := amqp.Dial(connStr)
	if connErr != nil {
		log.Println("Failed to connect to RabbitMQ: ", connErr)
	}
	channel, channelErr := conn.Channel()
	if channelErr != nil {
		log.Println("Failed to open the channel RabbitMQ: ", channelErr)
	}
	return &rabbitmqSerivce{
		conn:    conn,
		channel: channel,
	}
}

func (r *rabbitmqSerivce) GetQueue(topic string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		topic, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

func (r *rabbitmqSerivce) Publish(queue amqp.Queue, message string) error {
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

func (r *rabbitmqSerivce) Consume(q amqp.Queue) (<-chan amqp.Delivery, error) {
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

func (r *rabbitmqSerivce) Close() {
	r.channel.Close()
	r.conn.Close()
}
