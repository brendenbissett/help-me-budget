package admin

import (
	"strconv"

	"github.com/brendenbissett/help-me-budget/api/internal/auth"
	"github.com/brendenbissett/help-me-budget/api/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// GetUsers returns a paginated list of all users
func GetUsers(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100 // Cap at 100
	}

	// Get total count
	totalCount, err := auth.CountAllUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to count users",
		})
	}

	users, err := auth.GetAllUsers(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch users",
		})
	}

	// Get roles for each user
	type UserWithRoles struct {
		User  auth.User   `json:"user"`
		Roles []auth.Role `json:"roles"`
	}

	usersWithRoles := make([]UserWithRoles, 0, len(users))
	for _, user := range users {
		roles, err := auth.GetUserRoles(c.Context(), user.ID)
		if err != nil {
			roles = []auth.Role{} // Continue with empty roles on error
		}
		usersWithRoles = append(usersWithRoles, UserWithRoles{
			User:  user,
			Roles: roles,
		})
	}

	return c.JSON(fiber.Map{
		"users": usersWithRoles,
		"total": totalCount,
		"limit": limit,
		"offset": offset,
	})
}

// DeactivateUserHandler deactivates a user account
func DeactivateUserHandler(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	adminID := getUserIDFromContext(c)
	if adminID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var body struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := auth.DeactivateUser(c.Context(), userID, adminID, body.Reason); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to deactivate user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deactivated successfully",
	})
}

// ReactivateUserHandler reactivates a user account
func ReactivateUserHandler(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	adminID := getUserIDFromContext(c)
	if adminID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	if err := auth.ReactivateUser(c.Context(), userID, adminID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to reactivate user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User reactivated successfully",
	})
}

// DeleteUserHandler permanently deletes a user
func DeleteUserHandler(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	adminID := getUserIDFromContext(c)
	if adminID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var body struct {
		Reason string `json:"reason"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := auth.DeleteUser(c.Context(), userID, adminID, body.Reason); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete user",
		})
	}

	return c.JSON(fiber.Map{
		"message": "User deleted successfully",
	})
}

// GetAuditLogsHandler returns paginated audit logs
func GetAuditLogsHandler(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	if limit > 100 {
		limit = 100 // Cap at 100
	}

	logs, err := auth.GetAuditLogs(c.Context(), limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch audit logs",
		})
	}

	return c.JSON(fiber.Map{
		"logs":   logs,
		"limit":  limit,
		"offset": offset,
	})
}

// GetActiveSessions returns active Redis sessions
func GetActiveSessions(c *fiber.Ctx) error {
	ctx := c.Context()

	// Scan for all session keys in Redis
	var cursor uint64
	var sessions []map[string]interface{}

	for {
		var keys []string
		var err error

		keys, cursor, err = database.RedisClient.Scan(ctx, cursor, "session:*", 100).Result()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to scan sessions",
			})
		}

		// Get session data for each key
		for _, key := range keys {
			val, err := database.RedisClient.Get(ctx, key).Result()
			if err != nil {
				continue // Skip if error reading session
			}

			sessions = append(sessions, map[string]interface{}{
				"key":   key,
				"value": val,
			})
		}

		if cursor == 0 {
			break
		}
	}

	return c.JSON(fiber.Map{
		"sessions": sessions,
		"count":    len(sessions),
	})
}

// KillSession deletes a specific Redis session
func KillSession(c *fiber.Ctx) error {
	sessionKey := c.Params("key")
	if sessionKey == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Session key is required",
		})
	}

	adminID := getUserIDFromContext(c)
	ctx := c.Context()

	// Delete the session from Redis
	err := database.RedisClient.Del(ctx, sessionKey).Err()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete session",
		})
	}

	// Log the action
	_ = auth.CreateAuditLog(ctx, &auth.AuditLog{
		ActorID:      &adminID,
		Action:       "session.kill",
		ResourceType: "session",
		Details: map[string]interface{}{
			"session_key": sessionKey,
		},
		IPAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
	})

	return c.JSON(fiber.Map{
		"message": "Session killed successfully",
	})
}

// GrantRoleHandler grants a role to a user
func GrantRoleHandler(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	adminID := getUserIDFromContext(c)
	if adminID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var body struct {
		RoleName string `json:"role_name"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get role by name
	role, err := auth.GetRoleByName(c.Context(), body.RoleName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Role not found",
		})
	}

	if err := auth.GrantRole(c.Context(), userID, role.ID, adminID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to grant role",
		})
	}

	// Log the action
	_ = auth.CreateAuditLog(c.Context(), &auth.AuditLog{
		ActorID:      &adminID,
		Action:       "role.grant",
		ResourceType: "user",
		ResourceID:   &userID,
		Details: map[string]interface{}{
			"role_name": body.RoleName,
		},
		IPAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
	})

	return c.JSON(fiber.Map{
		"message": "Role granted successfully",
	})
}

// RevokeRoleHandler revokes a role from a user
func RevokeRoleHandler(c *fiber.Ctx) error {
	userID, err := uuid.Parse(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid user ID",
		})
	}

	adminID := getUserIDFromContext(c)
	if adminID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	var body struct {
		RoleName string `json:"role_name"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Get role by name
	role, err := auth.GetRoleByName(c.Context(), body.RoleName)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Role not found",
		})
	}

	if err := auth.RevokeRole(c.Context(), userID, role.ID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to revoke role",
		})
	}

	// Log the action
	_ = auth.CreateAuditLog(c.Context(), &auth.AuditLog{
		ActorID:      &adminID,
		Action:       "role.revoke",
		ResourceType: "user",
		ResourceID:   &userID,
		Details: map[string]interface{}{
			"role_name": body.RoleName,
		},
		IPAddress: c.IP(),
		UserAgent: c.Get("User-Agent"),
	})

	return c.JSON(fiber.Map{
		"message": "Role revoked successfully",
	})
}
