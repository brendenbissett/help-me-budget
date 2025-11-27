package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/brendenbissett/help-me-budget/api/internal/database"
	"github.com/google/uuid"
)

// Role represents a system role
type Role struct {
	ID          uuid.UUID
	Name        string
	Description string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// UserRole represents a user's role assignment
type UserRole struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	RoleID    uuid.UUID
	GrantedBy *uuid.UUID
	GrantedAt time.Time
}

// AuditLog represents an audit log entry
type AuditLog struct {
	ID           uuid.UUID
	ActorID      *uuid.UUID
	Action       string
	ResourceType string
	ResourceID   *uuid.UUID
	Details      map[string]interface{}
	IPAddress    string
	UserAgent    string
	CreatedAt    time.Time
}

// GetRoleByName retrieves a role by its name
func GetRoleByName(ctx context.Context, name string) (*Role, error) {
	var role Role
	err := database.DB.QueryRow(ctx, `
		SELECT id, name, description, created_at, updated_at
		FROM auth.roles
		WHERE name = $1
	`, name).Scan(
		&role.ID, &role.Name, &role.Description,
		&role.CreatedAt, &role.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get role: %w", err)
	}
	return &role, nil
}

// GetUserRoles retrieves all roles for a given user
func GetUserRoles(ctx context.Context, userID uuid.UUID) ([]Role, error) {
	rows, err := database.DB.Query(ctx, `
		SELECT r.id, r.name, r.description, r.created_at, r.updated_at
		FROM auth.roles r
		INNER JOIN auth.user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = $1
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query user roles: %w", err)
	}
	defer rows.Close()

	var roles []Role
	for rows.Next() {
		var role Role
		if err := rows.Scan(
			&role.ID, &role.Name, &role.Description,
			&role.CreatedAt, &role.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan role: %w", err)
		}
		roles = append(roles, role)
	}

	return roles, nil
}

// HasRole checks if a user has a specific role
func HasRole(ctx context.Context, userID uuid.UUID, roleName string) (bool, error) {
	var exists bool
	err := database.DB.QueryRow(ctx, `
		SELECT EXISTS(
			SELECT 1
			FROM auth.user_roles ur
			INNER JOIN auth.roles r ON r.id = ur.role_id
			WHERE ur.user_id = $1 AND r.name = $2
		)
	`, userID, roleName).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check role: %w", err)
	}
	return exists, nil
}

// GrantRole assigns a role to a user
func GrantRole(ctx context.Context, userID, roleID, grantedBy uuid.UUID) error {
	_, err := database.DB.Exec(ctx, `
		INSERT INTO auth.user_roles (user_id, role_id, granted_by)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, role_id) DO NOTHING
	`, userID, roleID, grantedBy)
	if err != nil {
		return fmt.Errorf("failed to grant role: %w", err)
	}
	return nil
}

// RevokeRole removes a role from a user
func RevokeRole(ctx context.Context, userID, roleID uuid.UUID) error {
	_, err := database.DB.Exec(ctx, `
		DELETE FROM auth.user_roles
		WHERE user_id = $1 AND role_id = $2
	`, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to revoke role: %w", err)
	}
	return nil
}

// CreateAuditLog creates a new audit log entry
func CreateAuditLog(ctx context.Context, log *AuditLog) error {
	_, err := database.DB.Exec(ctx, `
		INSERT INTO auth.audit_logs (actor_id, action, resource_type, resource_id, details, ip_address, user_agent)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, log.ActorID, log.Action, log.ResourceType, log.ResourceID, log.Details, log.IPAddress, log.UserAgent)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}
	return nil
}

// GetAuditLogs retrieves audit logs with pagination
func GetAuditLogs(ctx context.Context, limit, offset int) ([]AuditLog, error) {
	rows, err := database.DB.Query(ctx, `
		SELECT id, actor_id, action, resource_type, resource_id, details, ip_address, user_agent, created_at
		FROM auth.audit_logs
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query audit logs: %w", err)
	}
	defer rows.Close()

	var logs []AuditLog
	for rows.Next() {
		var log AuditLog
		if err := rows.Scan(
			&log.ID, &log.ActorID, &log.Action, &log.ResourceType,
			&log.ResourceID, &log.Details, &log.IPAddress, &log.UserAgent, &log.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan audit log: %w", err)
		}
		logs = append(logs, log)
	}

	return logs, nil
}

// DeactivateUser marks a user as inactive
func DeactivateUser(ctx context.Context, userID, deactivatedBy uuid.UUID, reason string) error {
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update user
	_, err = tx.Exec(ctx, `
		UPDATE auth.users
		SET is_active = false, deactivated_at = $1, deactivated_by = $2
		WHERE id = $3
	`, time.Now(), deactivatedBy, userID)
	if err != nil {
		return fmt.Errorf("failed to deactivate user: %w", err)
	}

	// Create audit log
	_, err = tx.Exec(ctx, `
		INSERT INTO auth.audit_logs (actor_id, action, resource_type, resource_id, details)
		VALUES ($1, 'user.deactivate', 'user', $2, $3)
	`, deactivatedBy, userID, map[string]interface{}{"reason": reason})
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return tx.Commit(ctx)
}

// ReactivateUser marks a user as active again
func ReactivateUser(ctx context.Context, userID, reactivatedBy uuid.UUID) error {
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Update user
	_, err = tx.Exec(ctx, `
		UPDATE auth.users
		SET is_active = true, deactivated_at = NULL, deactivated_by = NULL
		WHERE id = $1
	`, userID)
	if err != nil {
		return fmt.Errorf("failed to reactivate user: %w", err)
	}

	// Create audit log
	_, err = tx.Exec(ctx, `
		INSERT INTO auth.audit_logs (actor_id, action, resource_type, resource_id)
		VALUES ($1, 'user.reactivate', 'user', $2)
	`, reactivatedBy, userID)
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	return tx.Commit(ctx)
}

// DeleteUser permanently deletes a user
func DeleteUser(ctx context.Context, userID, deletedBy uuid.UUID, reason string) error {
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Create audit log before deletion
	_, err = tx.Exec(ctx, `
		INSERT INTO auth.audit_logs (actor_id, action, resource_type, resource_id, details)
		VALUES ($1, 'user.delete', 'user', $2, $3)
	`, deletedBy, userID, map[string]interface{}{"reason": reason})
	if err != nil {
		return fmt.Errorf("failed to create audit log: %w", err)
	}

	// Delete user (cascade will handle related records)
	_, err = tx.Exec(ctx, `
		DELETE FROM auth.users
		WHERE id = $1
	`, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return tx.Commit(ctx)
}

// GetAllUsers retrieves all users with pagination
func GetAllUsers(ctx context.Context, limit, offset int) ([]User, error) {
	rows, err := database.DB.Query(ctx, `
		SELECT id, email, email_verified, name, avatar_url, is_active, created_at, updated_at, last_login_at
		FROM auth.users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to query users: %w", err)
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(
			&user.ID, &user.Email, &user.EmailVerified, &user.Name,
			&user.AvatarURL, &user.IsActive, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	return users, nil
}
