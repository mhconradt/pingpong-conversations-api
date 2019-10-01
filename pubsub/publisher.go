package pubsub

import (
	"context"
	"fmt"
	"github.com/apache/pulsar/pulsar-client-go/pulsar"
	_ "github.com/apache/pulsar/pulsar-client-go/pulsar"
)

// Create middleware that converts the request to the pb data type

// Define a publisher layer that publishes messages of a given type to a given topic
// Bonus points if it allows for channel semantics
// Writable streams for sending messages, please?

type Producer struct {
	TopicName string
	producer pulsar.Producer
}

func (p *Producer) Close()  {
	if err := p.producer.Close(); err != nil {
		fmt.Println("error closing producer: ", err)
	}
}

func (p *Producer) Send(msg []byte) error {
	return p.producer.Send(context.Background(), pulsar.ProducerMessage{
		Payload: msg,
	})
}

func NewProducer(t string) (*Producer, error) {
	client, err := NewClient()
	if err != nil {
		return new(Producer), err
	}
	po := pulsar.ProducerOptions{
		Topic: t,
	}
	p, err := client.CreateProducer(po)
	if err != nil {
		return new(Producer), err
	}
	return &Producer{
		TopicName: t,
		producer: p,
	}, nil
}

func SendOnce(t string, msg []byte) error {
	p, err := NewProducer(t)
	if err != nil {
		return err
	}
	return p.Send(msg)
}
