package main

import (
	"fmt"
	"log"
	"os"

	"github.com/brendenbissett/help-me-budget/api/internal/admin"
	"github.com/brendenbissett/help-me-budget/api/internal/auth"
	"github.com/brendenbissett/help-me-budget/api/internal/database"
	"github.com/brendenbissett/help-me-budget/api/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
)

func main() {

	API_VERSION := "0.0.1"

	// Load .env file, only if in non-production environment
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	// Initialize database connection
	if err := database.InitDatabase(); err != nil {
		log.Fatal("Error initializing database:", err)
	}
	defer database.Close()

	// Initialize Redis connection (used for caching and future admin dashboard metrics)
	if err := database.InitRedis(); err != nil {
		log.Fatal("Error initializing Redis:", err)
	}
	defer database.CloseRedis()

	// Validate API secret key is set
	apiSecret := os.Getenv("API_SECRET_KEY")
	if apiSecret == "" {
		log.Fatal("API_SECRET_KEY environment variable is required")
	}

	app := fiber.New(fiber.Config{
		// Increase header size limit to handle large browser cookies/headers
		ReadBufferSize: 16384, // 16KB (default is 4KB)
	})

	// Set up middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:5173, http://127.0.0.1:5173",
		AllowHeaders:     "Origin, Content-Type, Accept, X-User-ID, X-API-Key",
		AllowCredentials: true,
		AllowMethods:     "GET, POST, DELETE, PUT, PATCH, OPTIONS",
	}))

	// API key authentication (validates requests come from SvelteKit)
	app.Use(middleware.ValidateAPIKey())

	// Add user context middleware (extracts user ID from X-User-ID header)
	app.Use(admin.SetUserContext())

	// Health check endpoint
	app.Get("/", func(c *fiber.Ctx) error {
		msg := fmt.Sprintf("Welcome to Help-Me-Budget API (V %s)", API_VERSION)
		return c.SendString(msg)
	})

	// Setup authentication routes
	auth.SetupAuthRoutes(app)

	// Setup admin routes
	admin.SetupAdminRoutes(app)

	log.Fatal(app.Listen(":3000"))
}
