// Package algorithm contains all rate limiting algorithms that define the Strategy interface
package algorithm

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	interfaces "github.com/thegeekywanderer/fluxy/pkg/v1"
)

const (
	keyWithoutExpire = -1
)

// FixedWindow strategy
type FixedWindow struct{
	client 	*redis.Client
	now 	func() time.Time
}

// Run function implements the fixed window rate limiter logic
func (fw *FixedWindow) Run(req *interfaces.Request) (*interfaces.Result, error) {
	p := fw.client.Pipeline()
	incrResult := p.Incr(req.Key)
	ttlResult := p.TTL( req.Key)
	if _, err := p.Exec(); err != nil {
		return nil, fmt.Errorf("failed to execute increment to key %v", req.Key)
	}	
	totalRequests, err := incrResult.Result()
	if err != nil {
		return nil, fmt.Errorf("failed to increment key %v", req.Key)
	}

	var ttlDuration time.Duration
	if d, err := ttlResult.Result(); err != nil || d.Seconds() == keyWithoutExpire {
		ttlDuration = req.Duration
		if err := fw.client.Expire(req.Key, req.Duration).Err(); err != nil {
			return nil, fmt.Errorf("failed to set an expiration to key %v", req.Key)
		}
	} else {
		ttlDuration = d
	}

	expiresAt := fw.now().Add(ttlDuration)

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

