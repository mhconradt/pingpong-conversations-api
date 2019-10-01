package pubsub

import (
	"github.com/apache/pulsar/pulsar-client-go/pulsar"
	"os"
	"runtime"
)

func GetPulsarEndpoint() string {
	if e, found := os.LookupEnv("PULSAR_URL"); found {
		return e
	}
	return "pulsar://localhost:6650"
}

func NewClient() (pulsar.Client, error) {
	e := GetPulsarEndpoint()
	return pulsar.NewClient(pulsar.ClientOptions{
		URL: e,
		OperationTimeoutSeconds: 5,
		MessageListenerThreads: runtime.NumCPU(),
	})
}
