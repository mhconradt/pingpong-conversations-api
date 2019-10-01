package pubsub

import (
	"fmt"
	"github.com/apache/pulsar/pulsar-client-go/pulsar"
)

// Maybe accepting a

// Define a subscriber layer for subscribing to particular topics
// Bonus points if it adheres to idiomatic concurrency, but especially on this one.

// So the output is ideally a channel of the desired message type

type Subscriber struct {
	TopicName, SubscriberName string
	Messages                  chan []byte
	consumer                  pulsar.Consumer
}

func (s *Subscriber) Close() {
	fmt.Println("closing")
	if err := s.consumer.Close(); err != nil {
		fmt.Println("error closing consumer:", err)
	}
	close(s.Messages)
}

func NewSubscriber(tn, sn string) (*Subscriber, error) {
	client, err := NewClient()
	if err != nil {
		return new(Subscriber), err
	}
	mc := make(chan pulsar.ConsumerMessage)
	consumerOpts := pulsar.ConsumerOptions{
		Topic:            tn,
		SubscriptionName: sn,
		Type:             pulsar.Exclusive,
		MessageChannel:   mc,
	}
	c, err := client.Subscribe(consumerOpts)
	if err != nil {
		return new(Subscriber), err
	}
	dataChan := make(chan []byte)
	go func() {
		for msg := range mc {
			dataChan <- msg.Payload()
			if err := c.Ack(msg.Message); err != nil {
				fmt.Println("error acknowledging message: ", err)
			}
		}
	}()
	return &Subscriber{
		TopicName:      tn,
		SubscriberName: sn,
		Messages:       dataChan,
		consumer: c,
	}, nil
	// Create a channel that
}
