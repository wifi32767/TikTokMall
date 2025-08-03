package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wifi32767/TikTokMall/app/order/biz/dal"
	"github.com/wifi32767/TikTokMall/app/order/conf"
	"github.com/wifi32767/TikTokMall/common/logger"
	"github.com/wifi32767/TikTokMall/common/mtl"
	order "github.com/wifi32767/TikTokMall/rpc/kitex_gen/order/orderservice"
)

func main() {
	// log
	conn, ch, cancel := loggerInit()
	defer conn.Close()
	defer ch.Close()
	defer cancel()
	klog.SetLevel(conf.LogLevel())
	// dal
	dal.MysqlInit()
	// prometheus
	mtl.InitMetric(conf.GetConf().Kitex.Service, conf.GetConf().Kitex.PrometheusPort, conf.GetConf().Kitex.ConsulAddress)
	// kitex
	opts := kitexInit()
	svr := order.NewServer(new(OrderServiceImpl), opts...)

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}

func loggerInit() (*amqp.Connection, *amqp.Channel, context.CancelFunc) {
	conn, err := amqp.Dial(conf.GetConf().Log.RabbitmqAddress)
	if err != nil {
		panic("Logger: Failed to connect to RabbitMQ: " + err.Error())
	}

	ch, err := conn.Channel()
	if err != nil {
		panic("Logger: Failed to open a channel" + err.Error())
	}

	q, err := ch.QueueDeclare(
		"log", // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		panic("Logger: Failed to declare a queue" + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	klog.SetLogger(logger.NewLogger(
		ch,
		&q,
		"order service",
		ctx,
	))
	return conn, ch, cancel
}

func kitexInit() (opts []server.Option) {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1"+conf.GetConf().Kitex.Address)
	opts = append(opts, server.WithServiceAddr(addr))
	opts = append(opts, server.WithMetaHandler(transmeta.ServerTTHeaderHandler))

	// consul
	r, err := consul.NewConsulRegister(conf.GetConf().Kitex.ConsulAddress)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithRegistry(r))
	opts = append(opts, server.WithRegistryInfo(&registry.Info{
		ServiceName: conf.GetConf().Kitex.Service,
		Weight:      1,
	}))
	return
}
