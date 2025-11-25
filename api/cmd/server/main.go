package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Config struct {
	GOOGLE_KEY    string
	GOOGLE_SECRET string
}

func main() {

	API_VERSION := "0.0.1"

	// Load .env file, oly if in non-production environment
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Set up config from environemnt variables

	GOOGLE_KEY := os.Getenv("GOOGLE_KEY")
	if GOOGLE_KEY == "" {
		log.Fatal("Error loading the GOOGLE_KEY from environment variables")
	}

	GOOGLE_SECRET := os.Getenv("GOOGLE_SECRET")
	if GOOGLE_SECRET == "" {
		log.Fatal("Error loading the GOOGLE_SECRET from environment variables.")
	}

	_ = Config{
		GOOGLE_KEY:    GOOGLE_KEY,
		GOOGLE_SECRET: GOOGLE_SECRET,
	}

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Welcome to Help-Me-Budget API (V %s)", API_VERSION)
		return c.SendString(msg)
	})

	log.Fatal(app.Listen(":3000"))
}
