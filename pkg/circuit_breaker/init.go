package circuit_breaker

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func NewCircuitBreaker(circuitBreaker *CircuitBreaker) CircuitBreakerInterface {
	return circuitBreaker
}

func (cb *CircuitBreaker) Execute(cbOpenKey string, cbFailureCountKey string, operation func() (error, int)) (error, int) {
	context.WithTimeout(cb.Ctx, time.Duration(10) * time.Second)
	if cb.IsOpen && time.Since(*cb.LastFailure) > cb.ResetTimeout {
		cb.IsOpen = false
		if err := cb.RedisClient.Set(cb.Ctx, cbOpenKey, "false", 0).Err(); err != nil {
			msg := fmt.Sprintf("Error updating circuit breaker status in Redis: %v", err.Error())
			logrus.Error(msg)
			return fmt.Errorf(msg), fiber.StatusInternalServerError
		}
	}

	if cb.IsOpen { // Circuit breaker is open
		msg := fmt.Sprintf("Currently, our system is busy. Please try again in %.0f seconds.", cb.ResetTimeout.Seconds())
		return fmt.Errorf(msg), fiber.StatusServiceUnavailable
	}

	err, code := operation()
	if err != nil {
		if err := cb.RedisClient.Incr(cb.Ctx, cbFailureCountKey).Err(); err != nil {
			msg := fmt.Sprintf("Error updating failure count in Redis: %v", err.Error())
			logrus.Error(msg)
			return fmt.Errorf(msg), fiber.StatusInternalServerError
		}

		now := time.Now()
		cb.LastFailure = &now

		failureCount, err := cb.RedisClient.Get(cb.Ctx, cbFailureCountKey).Int()
		if err != nil && err != redis.Nil {
			msg := fmt.Sprintf("Error getting failure count from Redis: %v", err)
			logrus.Error(msg)
			return fmt.Errorf(msg), fiber.StatusInternalServerError
		}

		if failureCount >= cb.Threshold {
			cb.IsOpen = true
			if err := cb.RedisClient.Set(cb.Ctx, cbOpenKey, "true", 0).Err(); err != nil {
				msg := fmt.Sprintf("Error updating circuit breaker status in Redis: %v", err)
				logrus.Error(msg)
				return fmt.Errorf(msg), fiber.StatusInternalServerError
			}
			logrus.Println("Circuit breaker is open")
			go func() {
				time.Sleep(cb.ResetTimeout)
				cb.IsOpen = false
				cb.RedisClient.Del(cb.Ctx, cbOpenKey)
				cb.RedisClient.Del(cb.Ctx, cbFailureCountKey)
			}()
		}
	} else {
		if err := cb.RedisClient.Set(cb.Ctx, cbFailureCountKey, "0", 0).Err(); err != nil {
			msg := fmt.Sprintf("Error resetting failure count in Redis: %v", err)
			logrus.Error(msg)
			return fmt.Errorf(msg), fiber.StatusInternalServerError
		}
	}

	return err, code
}