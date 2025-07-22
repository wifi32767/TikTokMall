package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/transport"
	consul "github.com/kitex-contrib/registry-consul"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wifi32767/TikTokMall/app/async/conf"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/cart"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/cart/cartservice"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/order/orderservice"

	"github.com/wifi32767/TikTokMall/app/async/message"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/order"
)

var (
	OrderClient orderservice.Client
	CartClient  cartservice.Client
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://admin:123456@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"async", // name
		true,    // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	initService()
	Run(channel, &queue)
}

func initService() {
	r, err := consul.NewConsulResolver(conf.GetConf().ConsulAddress)
	if err != nil {
		panic(err)
	}
	CartClient = cartservice.MustNewClient(
		"cart",
		client.WithResolver(r),
		client.WithTransportProtocol(transport.GRPC),
	)
	OrderClient = orderservice.MustNewClient(
		"order",
		client.WithResolver(r),
		client.WithTransportProtocol(transport.GRPC),
	)
}

func Run(channel *amqp.Channel, queue *amqp.Queue) {
	var forever chan struct{}

	go func() {
		msgs, err := channel.Consume(
			queue.Name,
			"",
			true,
			false,
			false,
			false,
			nil,
		)
		failOnError(err, "Failed to consume")

		for s := range msgs {
			var msg message.AsyncMessage
			err = json.Unmarshal(s.Body, &msg)
			if err != nil {
				fmt.Println(err)
			}
			if msg.Head == "order.UpdateOrderState" {
				var body order.UpdateOrderStateReq
				if err = json.Unmarshal(msg.Body, &body); err != nil {
					fmt.Println(err)
				} else {
					OrderClient.UpdateOrderState(msg.Ctx, &body)
				}
			} else if msg.Head == "cart.EmptyCart" {
				var body cart.EmptyCartReq
				if err = json.Unmarshal(msg.Body, &body); err != nil {
					fmt.Println(err)
				} else {
					CartClient.EmptyCart(msg.Ctx, &body)
				}
			}
		}
	}()

	<-forever
}
