package pubsub

import (
	"github.com/Comcast/pulsar-client-go"
)

// Create middleware that converts the request to the pb data type

// Define a publisher layer that publishes messages of a given type to a given topic
// Bonus points if it allows for channel semantics
// Writable streams for sending messages, please?

type Producer struct {
	TopicName string
	producer *pulsar.Producer
}

func Publish(channel string, msg []byte) error {
	return client.Publish(channel, string(msg)).Err()
}
