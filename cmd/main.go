package main

import (
	"context"
	"fmt"
	redisConf "template/conf/databases/redis"
	"template/pkg/circuit_breaker"
	"template/pkg/dotenv"
	"time"

	"github.com/gofiber/fiber/v2"
)

func main() {
    app := fiber.New()
	redisConn := redisConf.NewConn(&redisConf.Conf{
		ClientName:      "Circuit_Breaker",
		Addr:            dotenv.GetString("REDIS_URL", ""),
		Username:        dotenv.GetString("REDIS_USERNAME", ""),
		Pass:            dotenv.GetString("REDIS_PASS", ""),
		DB:              dotenv.GetInt("REDIS_DB", 0),
		MaxRetries:      5,
		MaxActiveConns:  10,
		ConnMaxIdleTime: 0,
		ConnMaxLifetime: 0,
		MinIdleConns:    0,
		MaxIdleConns:    0,
		DialTimeout:     0,
		ReadTimeout:     0,
		WriteTimeout:    0,
		PoolTimeout:     0,
		PoolSize:        0,
		IsUseTls:        false,
		IsUseHooks:      false,
	})
	redisClient := redisConn.InitClient()
	
	ctx := context.Background()

	circuitBreaker := circuit_breaker.NewCircuitBreaker(&circuit_breaker.CircuitBreaker{
		Threshold:    3,
		ResetTimeout: 5 * time.Second,
		LastFailure:  &time.Time{},
		IsOpen:       false,
		RedisClient:  redisClient,
		Ctx:          ctx,
	})

    app.Get("/", func(c *fiber.Ctx) error {
		c.Response().Header.Add("Content-Type", "application/json")
		err := circuitBreaker.Execute("cb_open", "cb_count", func() error {
			if 2 == 2 { return fmt.Errorf("Error bro") }

			return nil
		})
		
		if err != nil { return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": err.Error()}) }

        return c.SendString(`{ "message" : "Service OK" }`)
    })
	
    app.Listen(":8080")
}