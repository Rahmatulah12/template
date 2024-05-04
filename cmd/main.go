package main

import (
	"context"
	"fmt"
	"os"
	"template/conf/apm"
	redisConf "template/conf/databases/redis"
	"template/pkg/circuit_breaker"
	"template/pkg/dotenv"
	"time"

	"github.com/gofiber/contrib/fibernewrelic"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/sirupsen/logrus"
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
		MaxActiveConns:  100,
		ConnMaxIdleTime: 0,
		ConnMaxLifetime: 0,
		MinIdleConns:    5,
		MaxIdleConns:    10,
		DialTimeout:     10,
		ReadTimeout:     10,
		WriteTimeout:    10,
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

	file, err := os.OpenFile(fmt.Sprintf("./logs/%s.log", "http-logs"), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		logrus.Panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer file.Close()

	app.Use(recover.New(), logger.New(logger.Config{
		Format:     "${pid} | ${time} | ${latency} |  ${status} | ${protocol} | ${ip} | ${host} | ${method} | ${path} | reqHeaders: ${reqHeaders} | queryParams: ${queryParams} | req: ${body} | res: ${resBody} | error: ${error}\n\n",
		TimeFormat: "2006-01-02 15:04:05.00",
		TimeZone:   "Asia/Jakarta",
		DisableColors: true,
		Output: file,
	}))

    app.Get("/", func(c *fiber.Ctx) error {
		c.Response().Header.Add("Content-Type", "application/json")
		err, httpCode := circuitBreaker.Execute("cb_open", "cb_count", func() (error, int) {
			if 2 == 2 { return fmt.Errorf("Error bro"), fiber.StatusBadRequest }

			return nil, fiber.StatusOK
		})
		
		if err != nil { return c.Status(httpCode).JSON(fiber.Map{"message": err.Error()}) }

        return c.SendString(`{ "message" : "Service OK" }`)
    })

	Nr := apm.InstanceNewrelic(&apm.NewrelicConfig{
		AppName:    os.Getenv("APP_NAME_NEWRELIC"),
		LicenseKey: os.Getenv("APP_KEY_NEWRELIC"),
	})
	
	cfg := fibernewrelic.Config{
		Application: Nr,
		Enabled: true,
	}

	app.Use(fibernewrelic.New(cfg))
	
    app.Listen(":8080")
}