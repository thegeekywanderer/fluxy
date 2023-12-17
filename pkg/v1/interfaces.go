// Package v1 contains interface definitions for v1API
package v1

import (
	"time"

	"github.com/thegeekywanderer/fluxy/models"
)

// RepoInterface defines an interface for repository of v1API
type RepoInterface interface{
  RegisterClient(models.Client) (models.Client, error)
  GetClient(name string) (models.Client, error)
  UpdateClient(models.Client) (error)
  DeleteClient(name string) (error)
}

// UseCaseInterface defines the logic interface of v1API
type UseCaseInterface interface{
  RegisterClient(models.Client) (models.Client, error)
  GetClient(name string) (models.Client, error)
  UpdateClient(models.Client) (error)
  DeleteClient(name string) (error)
  VerifyLimit(name string) (*Result, error)
}

// State defines the rate limiter state
type State int64

// Define constants for rate limit state
const (
	Deny  State = 0
	Allow       = 1
)

// Request struct defines a rate limit check request
type Request struct {
	Key      string
	Limit    uint64
	Duration time.Duration
}

// Result struct defines a rate limit check result
type Result struct {
	State         State
	TotalRequests uint64
	ExpiresAt     time.Time
}

// Strategy is an interface for rate limiting algorithms. Any rate limiting algorithm that implements this interface can be used with fluxy
type Strategy interface {
	Run(r *Request) (*Result, error)
}
