package middleware

import (
	"os"

	"github.com/gofiber/fiber/v2"
)

// ValidateAPIKey middleware ensures requests come from authorized clients (SvelteKit)
// Whitelists the health check endpoint (GET /) for monitoring purposes
func ValidateAPIKey() fiber.Handler {
	expectedKey := os.Getenv("API_SECRET_KEY")

	return func(c *fiber.Ctx) error {
		// Whitelist health check endpoint
		if c.Method() == "GET" && c.Path() == "/" {
			return c.Next()
		}

		// Get API key from X-API-Key header
		providedKey := c.Get("X-API-Key")

		// Validate key
		if providedKey == "" || providedKey != expectedKey {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized: Invalid or missing API key",
			})
		}

		return c.Next()
	}
}
