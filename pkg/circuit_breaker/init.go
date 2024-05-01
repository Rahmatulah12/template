package circuit_breaker

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func NewCircuitBreaker(circuitBreaker *CircuitBreaker) *CircuitBreaker {
	context.WithTimeout(ctx, time.Duration(10) * time.Second)
	return circuitBreaker
}

func (cb *CircuitBreaker) Execute(operation func() error) error {
	if cb.IsOpen && time.Since(*cb.LastFailure) > cb.ResetTimeout {
		cb.IsOpen = false
		if err := cb.RedisClient.Set(cb.Ctx, cb.CbOpenKey, "false", 0).Err(); err != nil {
			msg := fmt.Sprintf("Error updating circuit breaker status in Redis: %v", err.Error())
			logrus.Error(msg)
			return fmt.Errorf(msg)
		}
	}

	if cb.IsOpen { // Circuit breaker is open
		return fmt.Errorf("Currently, our system is busy. Please try again in a few moments.")
	}

	err := operation()
	if err != nil {
		if err := cb.RedisClient.Incr(cb.Ctx, cb.CbFailureCountKey).Err(); err != nil {
			msg := fmt.Sprintf("Error updating failure count in Redis: %v", err.Error())
			logrus.Error(msg)
			return fmt.Errorf(msg)
		}

		now := time.Now()
		cb.LastFailure = &now

		failureCount, err := cb.RedisClient.Get(cb.Ctx, cb.CbFailureCountKey).Int()
		if err != nil && err != redis.Nil {
			msg := fmt.Sprintf("Error getting failure count from Redis: %v", err)
			logrus.Error(msg)
			return fmt.Errorf(msg)
		}

		if failureCount >= cb.Threshold {
			cb.IsOpen = true
			if err := cb.RedisClient.Set(cb.Ctx, cb.CbOpenKey, "true", 0).Err(); err != nil {
				msg := fmt.Sprintf("Error updating circuit breaker status in Redis: %v", err)
				logrus.Error(msg)
				return fmt.Errorf(msg)
			}
			logrus.Println("Circuit breaker is open")
			go func() {
				time.Sleep(cb.ResetTimeout)
				cb.IsOpen = false
				cb.RedisClient.Del(cb.Ctx, cb.CbOpenKey)
				cb.RedisClient.Del(cb.Ctx, cb.CbFailureCountKey)
			}()
		}
	} else {
		if err := cb.RedisClient.Set(cb.Ctx, cb.CbFailureCountKey, "0", 0).Err(); err != nil {
			msg := fmt.Sprintf("Error resetting failure count in Redis: %v", err)
			logrus.Error(msg)
			return fmt.Errorf(msg)
		}
	}

	return err
}