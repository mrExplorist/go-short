package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/mrExplorist/shorten-url-fiber-redis/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// 1. Load the .env file
// 2. Create a new fiber app
// 3. Setup the routes
// 4. Listen on the port

func setupRoutes(app *fiber.App) {
	app.Get("/:url", routes.ResolveURL)
	app.Post("/api/v1", routes.ShortenURL) 

}

func main(){
	err := godotenv.Load() 

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := fiber.New() // create a new fiber app 

	app.Use(logger.New()) // use the logger middleware

	setupRoutes(app) // setup the routes 

	log.Fatal(app.Listen(os.Getenv("APP_PORT"))) // listen on the port specified in the .env file 

}