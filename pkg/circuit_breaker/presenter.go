package circuit_breaker

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()
)

type CircuitBreaker struct {
	Threshold     int
	ResetTimeout  time.Duration
	LastFailure  *time.Time
	IsOpen       bool
	RedisClient  *redis.Client
	CbOpenKey	 string
	CbFailureCountKey	string
	Ctx          context.Context
}