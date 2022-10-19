package main

import (
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"log"
	"tag-qr-service-example/routes"
)

var (
	port = flag.String("port", ":4004", "Port to listen on")
)

// This sample application has been written to easily explain how to integrate Alternatıf SuperApp Tag-QR.
// Explains sample requests and hash validations with examples.
func main() {
	loadEnv()

	app := fiber.New()

	setupRoutes(app)

	app.Use(cors.New())

	app.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusNotFound).
			JSON(fiber.Map{
				"message": "Not found",
			})
	})

	log.Fatal(app.Listen(*port))
}

// You can name the routes as you wish.
func setupRoutes(router *fiber.App) {
	router.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Example Tag QR Service")
	})

	router.Post("/get-price", routes.GetPrice)

	// Open Door is just a name given as an example.
	router.Post("/open-gate", routes.OpenGate)
}

func loadEnv() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}
}
