package pubsub

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var client *redis.ClusterClient

func NewClient() *redis.ClusterClient {
	opts := GetRedisOpts()
	return redis.NewClusterClient(opts)
}

func GetRedisOpts() *redis.ClusterOptions {
	return &redis.ClusterOptions{
		Addrs:       ClusterAddresses(6),
		PoolSize:    100,
		MaxConnAge:  15 * time.Second,
		IdleTimeout: 5 * time.Second,
	}
}

func ClusterAddresses(count int) []string {
	addresses := make([]string, count, count)
	for i := 0; i < count; i++ {
		f := "redis-cluster-%v.redis-cluster.default.svc.cluster.local:6379"
		addresses[i] = fmt.Sprintf(f, i)
	}
	return addresses
}

func init() {
	client = NewClient()
}
