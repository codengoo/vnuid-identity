package main

import (
	"log"
	"vnuid-identity/databases"
	"vnuid-identity/routes"

	"github.com/gofiber/fiber/v2"
)

func main() {
	databases.ConnectDB()

	app := fiber.New()

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
