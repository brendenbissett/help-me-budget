package auth

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

// SupabaseUserSyncRequest represents the request payload for syncing a Supabase user
type SupabaseUserSyncRequest struct {
	Email          string `json:"email"`
	Name           string `json:"name"`
	AvatarURL      string `json:"avatar_url"`
	Provider       string `json:"provider"`
	ProviderUserID string `json:"provider_user_id"`
}

// HandleSupabaseUserSync creates or updates a user in local PostgreSQL from Supabase auth data
func HandleSupabaseUserSync(c *fiber.Ctx) error {
	var req SupabaseUserSyncRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println("failed to parse request body:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate required fields
	if req.Email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}
	if req.Provider == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Provider is required",
		})
	}
	if req.ProviderUserID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Provider user ID is required",
		})
	}

	// Use existing UpsertUserWithOAuth function to sync the user
	dbUser, err := UpsertUserWithOAuth(
		c.Context(),
		req.Email,
		req.Name,
		req.AvatarURL,
		req.Provider,
		req.ProviderUserID,
	)
	if err != nil {
		log.Println("failed to sync user to local PostgreSQL:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to sync user",
		})
	}

	// Note: Session management is now handled by Supabase, not Redis
	// Redis sessions have been removed in favor of Supabase as the auth master

	// Return the synced user data
	return c.JSON(fiber.Map{
		"success": true,
		"user": fiber.Map{
			"id":          dbUser.ID.String(),
			"email":       dbUser.Email,
			"name":        dbUser.Name,
			"avatar_url":  dbUser.AvatarURL,
			"is_active":   dbUser.IsActive,
			"created_at":  dbUser.CreatedAt,
			"updated_at":  dbUser.UpdatedAt,
			"last_login":  dbUser.LastLoginAt,
		},
	})
}
