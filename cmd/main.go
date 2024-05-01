package main

import "github.com/gofiber/fiber/v2"

func main() {
    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
		c.Response().Header.Add("Content-Type", "application/json")
        return c.SendString(`{ "message" : "Service OK" }`)
    })

    app.Listen(":8080")
}