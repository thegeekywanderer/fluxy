// Package algorithm contains all rate limiting algorithms that define the Strategy interface
package algorithm

import (
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/google/uuid"
	interfaces "github.com/thegeekywanderer/fluxy/pkg/v1"
)

// RollingWindow strategy
type RollingWindow struct{
	client 	*redis.Client
	now 	func() time.Time
}

// Run function implements the fixed window rate limiter logic
func (rw *RollingWindow) Run(req *interfaces.Request) (*interfaces.Result, error) {
	now := rw.now()
	// every request needs an UUID
	item := uuid.New()

	minimum := now.Add(-req.Duration)
	p := rw.client.Pipeline()
	// we then remove all requests that have already expired on this set
	removeByScore := p.ZRemRangeByScore(req.Key, "0", strconv.FormatInt(minimum.UnixMilli(), 10))

	// we add the current request
	add := p.ZAdd(req.Key, redis.Z{
		Score:  float64(now.UnixMilli()),
		Member: item.String(),
	})

	// count how many non-expired requests we have on the sorted set
	count := p.ZCount(req.Key, "-inf", "+inf")

	if _, err := p.Exec(); err != nil {
		return nil, fmt.Errorf("failed to execute sorted set pipeline for key: %v", req.Key)
	}

	if err := removeByScore.Err(); err != nil {
		return nil, fmt.Errorf("failed to remove items from key %v", req.Key)
	}

	if err := add.Err(); err != nil {
		return nil, fmt.Errorf("failed to add item to key %v", req.Key)
	}

	totalRequests, err := count.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to count items for key %v", req.Key)
	}

	expiresAt := now.Add(req.Duration)
	requests := uint64(totalRequests)

	if requests > req.Limit {
		return &interfaces.Result{
			State:         interfaces.Deny,
			TotalRequests: requests,
			ExpiresAt:     expiresAt,
		}, nil
	}

	return &interfaces.Result{
		State:         interfaces.Allow,
		TotalRequests: requests,
		ExpiresAt:     expiresAt,
	}, nil
}
