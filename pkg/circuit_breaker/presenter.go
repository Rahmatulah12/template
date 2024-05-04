package circuit_breaker

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type CircuitBreaker struct {
	Threshold     int
	ResetTimeout  time.Duration
	LastFailure  *time.Time
	IsOpen       bool
	RedisClient  *redis.Client
	Ctx          context.Context
}

type CircuitBreakerInterface interface {
	Execute(cbOpenKey string, cbFailureCountKey string, operation func() (error, int)) (error, int)
}