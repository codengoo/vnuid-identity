package main

import (
	"fmt"
	"log"
	"os"
	"vnuid-identity/databases"
	"vnuid-identity/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	dba := os.Getenv("DATABASE_URL")
	fmt.Println("Port: ", dba)
	fmt.Println("Setup env successfully")
}

func main() {
	// Get port from environment variable, default to 8080
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	databases.ConnectDB()
	databases.ConnectRD()

	app := fiber.New()

	routes.SetupRoutes(app)

	log.Printf("Server started at http://0.0.0.0:%s", port)
	log.Fatal(app.Listen(":" + port))
}
