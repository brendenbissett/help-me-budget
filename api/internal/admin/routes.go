package admin

import (
	"github.com/gofiber/fiber/v2"
)

// SetupAdminRoutes registers all admin routes
func SetupAdminRoutes(app *fiber.App) {
	// Admin routes group
	admin := app.Group("/admin")

	// Apply admin authentication middleware to all admin routes
	admin.Use(RequireAdmin())

	// User management
	admin.Get("/users", GetUsers)
	admin.Post("/users/:id/deactivate", DeactivateUserHandler)
	admin.Post("/users/:id/reactivate", ReactivateUserHandler)
	admin.Delete("/users/:id", DeleteUserHandler)
	admin.Post("/users/:id/roles/grant", GrantRoleHandler)
	admin.Post("/users/:id/roles/revoke", RevokeRoleHandler)

	// Session management
	admin.Get("/sessions", GetActiveSessions)
	admin.Delete("/sessions/:key", KillSession)

	// Audit logs (admin or moderator can view)
	auditLogs := admin.Group("/audit-logs")
	auditLogs.Use(RequireAdminOrModerator())
	auditLogs.Get("/", GetAuditLogsHandler)
}
