package pubsub

import (
	"github.com/go-redis/redis"
)

// Maybe accepting a

// Define a subscriber layer for subscribing to particular topics
// Bonus points if it adheres to idiomatic concurrency, but especially on this one.

// So the output is ideally a channel of the desired message type

type Subscriber struct {
	TopicName, SubscriberName string
	Messages                  chan []byte
	ps *redis.PubSub
}

func (s *Subscriber) Close() {
	_ = s.ps.Close()
	close(s.Messages)
}

func NewSubscriber(tn string) *Subscriber {
	ps := client.Subscribe(tn)
	dataChan := make(chan []byte)
	go func() {
		for msg := range ps.Channel() {
			dataChan <- []byte(msg.Payload)
		}
	}()
	return &Subscriber{
		TopicName:      tn,
		Messages:       dataChan,
	}
	// Create a channel that
}
