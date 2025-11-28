package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// HandleGetUserRoles returns the roles for the authenticated user
func HandleGetUserRoles(c *fiber.Ctx) error {
	// Get user ID from context (set by SetUserContext middleware)
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized - user ID not found",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	// Get all roles for the user
	roles, err := GetUserRoles(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user roles",
		})
	}

	// Check specific roles
	isAdmin, _ := HasRole(c.Context(), userID, "admin")
	isModerator, _ := HasRole(c.Context(), userID, "moderator")

	return c.JSON(fiber.Map{
		"user_id":     userID,
		"roles":       roles,
		"is_admin":    isAdmin,
		"is_moderator": isModerator,
	})
}
