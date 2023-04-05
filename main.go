package main

import (
	"inferno/routes"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	serv := fiber.New()

	routes.Cloudstorage(serv)

	serv.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, Inferno! 🔥")
	})

	log.Fatal(serv.Listen(":3473"))
}
