package auth

import (
	"github.com/gofiber/fiber/v2"
)

// SetupAuthRoutes configures all authentication routes
// Note: Authentication is handled by Supabase. This API only syncs user data to local PostgreSQL.
func SetupAuthRoutes(app *fiber.App) {
	// IMPORTANT: Register specific routes BEFORE wildcard routes

	// ROUTE: get user roles -> /auth/roles (must be before :provider wildcard)
	app.Get("/auth/roles", HandleGetUserRoles)

	// ROUTE: get user roles by email -> /auth/roles/by-email (must be before :provider wildcard)
	app.Get("/auth/roles/by-email", HandleGetUserRolesByEmail)

	// ROUTE: sync Supabase user to local PostgreSQL -> /auth/sync
	app.Post("/auth/sync", HandleSupabaseUserSync)
}
