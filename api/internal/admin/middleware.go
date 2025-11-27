package admin

import (
	"github.com/brendenbissett/help-me-budget/api/internal/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// RequireAdmin middleware ensures the user has admin role
func RequireAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := getUserIDFromContext(c)
		if userID == uuid.Nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized - please log in",
			})
		}

		// Check if user has admin role
		isAdmin, err := auth.HasRole(c.Context(), userID, "admin")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check permissions",
			})
		}

		if !isAdmin {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden - admin access required",
			})
		}

		return c.Next()
	}
}

// RequireAdminOrModerator middleware ensures the user has admin or moderator role
func RequireAdminOrModerator() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userID := getUserIDFromContext(c)
		if userID == uuid.Nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Unauthorized - please log in",
			})
		}

		// Check if user has admin or moderator role
		isAdmin, err := auth.HasRole(c.Context(), userID, "admin")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check permissions",
			})
		}

		if isAdmin {
			return c.Next()
		}

		isModerator, err := auth.HasRole(c.Context(), userID, "moderator")
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to check permissions",
			})
		}

		if !isModerator {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Forbidden - admin or moderator access required",
			})
		}

		return c.Next()
	}
}

// getUserIDFromContext extracts the user ID from the request context
// This assumes you have a middleware that sets the user ID in context after authentication
func getUserIDFromContext(c *fiber.Ctx) uuid.UUID {
	// Try to get user ID from locals (set by auth middleware)
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return uuid.Nil
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return uuid.Nil
	}

	return userID
}

// SetUserContext middleware extracts user info from request header and sets in context
// The SvelteKit backend will send the user ID via X-User-ID header after validating the session
func SetUserContext() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get user ID from header (set by SvelteKit proxy)
		userIDStr := c.Get("X-User-ID")
		if userIDStr != "" {
			// Validate it's a proper UUID
			if _, err := uuid.Parse(userIDStr); err == nil {
				c.Locals("user_id", userIDStr)
			}
		}

		return c.Next()
	}
}
