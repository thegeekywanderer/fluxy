package database

import (
	"github.com/go-redis/redis"
	"github.com/thegeekywanderer/fluxy/config"
)

// RedisConnection returns a redis connection
func RedisConnection(config *config.RedisConfiguration) (*redis.Client) {
	client := redis.NewClient(&redis.Options{
		Addr: config.Host,
		Password: config.Password,
		DB: 0,
	})
	return client
}