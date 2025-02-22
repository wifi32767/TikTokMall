package logger

import (
	"io"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Receiver struct {
	channel *amqp.Channel
	queue   *amqp.Queue
	output  io.Writer
}

func NewReceiver(channel *amqp.Channel, queue *amqp.Queue) *Receiver {
	return &Receiver{
		channel: channel,
		queue:   queue,
		output:  os.Stdout,
	}
}

func (r *Receiver) SetOutput(output io.Writer) {
	r.output = output
}

// 这个函数不会自旋，需要加一个循环保持运行
// 这样可以在开着接收器的情况下干别的
func (r *Receiver) Run() {
	msgs, err := r.channel.Consume(
		r.queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic("Reciever: " + err.Error())
	}
	go func() {
		for d := range msgs {
			r.output.Write(d.Body)
		}
	}()
}
