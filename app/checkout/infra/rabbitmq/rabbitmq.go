package rabbitmq

import (
	"context"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wifi32767/TikTokMall/app/checkout/conf"
)

var (
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      *amqp.Queue
	ctx        context.Context
)

func Init() {
	conn, err := amqp.Dial(conf.GetConf().AsyncMq.Address)
	if err != nil {
		panic("Failed to connect to RabbitMQ: " + err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		panic("Failed to open a channel" + err.Error())
	}

	q, err := ch.QueueDeclare(
		"async", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		panic("Failed to declare a queue" + err.Error())
	}
	connection = conn
	channel = ch
	queue = &q
}

func SendMessage(msg []byte) error {
	var err error
	if ctx == nil {
		err = channel.Publish(
			"",         // exchange
			queue.Name, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        msg,
			},
		)
	} else {
		err = channel.PublishWithContext(
			ctx,
			"",         // exchange
			queue.Name, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        msg,
			},
		)
	}
	return err
}
