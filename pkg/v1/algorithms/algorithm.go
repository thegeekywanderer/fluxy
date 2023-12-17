package algorithm

import (
	"errors"
	"time"

	"github.com/go-redis/redis"
	interfaces "github.com/thegeekywanderer/fluxy/pkg/v1"
)

// New function instantiates a new usecase
func New(algorithm string, client *redis.Client, now func() time.Time) (interfaces.Strategy, error) {
	if algorithm == "fixed-window" {
		return &FixedWindow{
			client: client,
			now:    now,
		}, nil
	} else if algorithm == "rolling-window" {
		return &RollingWindow{
			client: client,
			now: 	now,
		}, nil
	} 
	return nil, errors.New("Algorithm not implemented")
}