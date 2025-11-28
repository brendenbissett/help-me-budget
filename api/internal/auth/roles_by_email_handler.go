package auth

import (
	"github.com/gofiber/fiber/v2"
)

// HandleGetUserRolesByEmail returns the roles for a user identified by email
// This is used when the frontend only has the Supabase user data (email)
// and needs to look up roles in the local PostgreSQL database
func HandleGetUserRolesByEmail(c *fiber.Ctx) error {
	// Get email from query parameter
	email := c.Query("email")
	if email == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is required",
		})
	}

	// Get user by email from local PostgreSQL
	user, err := GetUserByEmail(c.Context(), email)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	// Get all roles for the user
	roles, err := GetUserRoles(c.Context(), user.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch user roles",
		})
	}

	// Check specific roles
	isAdmin, _ := HasRole(c.Context(), user.ID, "admin")
	isModerator, _ := HasRole(c.Context(), user.ID, "moderator")

	return c.JSON(fiber.Map{
		"user_id":      user.ID,
		"email":        user.Email,
		"roles":        roles,
		"is_admin":     isAdmin,
		"is_moderator": isModerator,
	})
}
