// Package algorithm contains all rate limiting algorithms that define the Strategy interface
package algorithm

import (
	interfaces "github.com/thegeekywanderer/fluxy/pkg/v1"
)

// RollingWindow strategy
type RollingWindow struct{}

// Run function implements the fixed window rate limiter logic
func (rw *RollingWindow) Run(req *interfaces.Request) (*interfaces.Result, error) {
	return nil, nil	
}

