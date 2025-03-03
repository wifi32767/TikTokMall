package main

import (
	"context"
	"net"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/server"
	consul "github.com/kitex-contrib/registry-consul"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wifi32767/TikTokMall/app/auth/conf"
	"github.com/wifi32767/TikTokMall/app/auth/dal"
	"github.com/wifi32767/TikTokMall/common/logger"
	auth "github.com/wifi32767/TikTokMall/rpc/kitex_gen/auth/authservice"
)

func main() {
	// log
	conn, ch, cancel := loggerInit()
	defer conn.Close()
	defer ch.Close()
	defer cancel()
	klog.SetLevel(conf.LogLevel())
	// redis
	dal.RedisInit()
	// kitex
	opts := kitexInit()
	svr := auth.NewServer(new(AuthServiceImpl), opts...)

	err := svr.Run()

	if err != nil {
		klog.Error(err.Error())
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
		"auth service",
		ctx,
	))
	return conn, ch, cancel
}

func kitexInit() (opts []server.Option) {
	addr, _ := net.ResolveTCPAddr("tcp", "127.0.0.1"+conf.GetConf().Kitex.Address)
	opts = append(opts, server.WithServiceAddr(addr))
	opts = append(opts, server.WithMetaHandler(transmeta.ServerTTHeaderHandler))

	// consul
	r, err := consul.NewConsulRegister(conf.GetConf().Kitex.Consul_address)
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
