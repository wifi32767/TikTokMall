package main

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/transport"
	consul "github.com/kitex-contrib/registry-consul"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/wifi32767/TikTokMall/app/async/conf"
	"github.com/wifi32767/TikTokMall/app/async/message"
	"github.com/wifi32767/TikTokMall/common/logger"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/cart/cartservice"
	"github.com/wifi32767/TikTokMall/rpc/kitex_gen/order/orderservice"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var (
	connection *amqp.Connection
	channel    *amqp.Channel
	queue      amqp.Queue
	MethodMap  map[string]reflect.Method = make(map[string]reflect.Method)
	ServiceMap map[string]any            = make(map[string]any)
)

func main() {
	// log
	conn, ch, cancel := loggerInit()
	defer conn.Close()
	defer ch.Close()
	defer cancel()
	klog.SetLevel(conf.LogLevel())
	// async mq
	initAsync()
	initService()
	defer connection.Close()
	defer channel.Close()
	Run()
}

func loggerInit() (*amqp.Connection, *amqp.Channel, context.CancelFunc) {
	conn, err := amqp.Dial(conf.GetConf().Log.MqAddress)
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
		"checkout service",
		ctx,
	))
	return conn, ch, cancel
}

func initAsync() {
	conn, err := amqp.Dial(conf.GetConf().MqAddress)
	if err != nil {
		klog.Fatal("Failed to connect to RabbitMQ: " + err.Error())
	}
	defer conn.Close()

	channel, err = conn.Channel()
	if err != nil {
		klog.Fatal("Failed to open a channel: " + err.Error())
	}
	defer channel.Close()

	queue, err = channel.QueueDeclare(
		"async", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		klog.Fatal("Failed to declare a queue: " + err.Error())
	}
}

func initService() {
	r, err := consul.NewConsulResolver(conf.GetConf().ConsulAddress)
	if err != nil {
		klog.Fatal(err)
	}
	addService("cart", cartservice.MustNewClient(
		"cart",
		client.WithResolver(r),
		client.WithTransportProtocol(transport.GRPC),
	))
	addService("order", orderservice.MustNewClient(
		"order",
		client.WithResolver(r),
		client.WithTransportProtocol(transport.GRPC),
	))
}

func addService(name string, service any) {
	ServiceMap[name] = service
	typ := reflect.ValueOf(service).Type()
	klog.Debug("add service: ", name)
	for i := 0; i < typ.NumMethod(); i++ {
		methodName := typ.Method(i).Name
		klog.Debugf("add method: %s.%s", name, methodName)
		MethodMap[name+"."+methodName] = typ.Method(i)
	}
}

func Run() {
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
		if err != nil {
			klog.Fatalf("Failed to consume messages: %v", err)
		}

		for s := range msgs {
			// 解码
			msg, err := message.Decode(s.Body)
			if err != nil {
				klog.Infof("decode failed: %v", err)
				continue
			}
			// 避免 nil ctx 引起 panic
			if msg.Ctx == nil {
				msg.Ctx = context.Background()
			}
			client, ok := ServiceMap[msg.ServiceName]
			if !ok {
				klog.Infof("service not found: %s", msg.ServiceName)
				continue
			}
			method, ok := MethodMap[msg.ServiceName+"."+msg.MethodName]
			if !ok {
				klog.Infof("method not found")
				continue
			}
			// 反射调用
			reqType := method.Type.In(2)
			body, ok := msg.Body.(map[string]any)
			if !ok {
				klog.Infof("body type error")
				continue
			}
			req, err := getReq(reqType, body)
			if err != nil {
				klog.Infof("get req failed: %v", err)
				continue
			}
			klog.Infof("call %s.%s with req: %v", msg.ServiceName, msg.MethodName, req)
			method.Func.Call([]reflect.Value{
				reflect.ValueOf(client),
				reflect.ValueOf(msg.Ctx),
				reflect.ValueOf(req)})
		}
	}()

	<-forever
}

func getReq(reqType reflect.Type, body map[string]any) (any, error) {
	var msg proto.Message
	if reqType.Kind() == reflect.Ptr {
		msg = reflect.New(reqType.Elem()).Interface().(proto.Message)
	} else {
		msg = reflect.New(reqType).Interface().(proto.Message)
	}
	// mapstructure无法识别protobuf标签，因此选择将map转为json再转换
	jsonData, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	err = protojson.Unmarshal(jsonData, msg)
	return msg, err
}
