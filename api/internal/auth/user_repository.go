package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/brendenbissett/help-me-budget/api/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// User represents a user in the database
type User struct {
	ID            uuid.UUID
	Email         string
	EmailVerified bool
	Name          string
	AvatarURL     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
	LastLoginAt   *time.Time
}

// UserOAuthProvider represents an OAuth provider linked to a user
type UserOAuthProvider struct {
	ID             uuid.UUID
	UserID         uuid.UUID
	Provider       string
	ProviderUserID string
	CreatedAt      time.Time
	UpdatedAt      time.Time
	LastUsedAt     time.Time
}

// UpsertUserWithOAuth creates or updates a user and links an OAuth provider
// Returns the user record (new or existing)
func UpsertUserWithOAuth(ctx context.Context, email, name, avatarURL, provider, providerUserID string) (*User, error) {
	tx, err := database.DB.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// First, try to find existing user by email
	var user User
	err = tx.QueryRow(ctx, `
		SELECT id, email, email_verified, name, avatar_url, created_at, updated_at, last_login_at
		FROM auth.users
		WHERE email = $1
	`, email).Scan(
		&user.ID, &user.Email, &user.EmailVerified, &user.Name,
		&user.AvatarURL, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)

	if err == pgx.ErrNoRows {
		// User doesn't exist, create new user
		err = tx.QueryRow(ctx, `
			INSERT INTO auth.users (email, email_verified, name, avatar_url, last_login_at)
			VALUES ($1, $2, $3, $4, $5)
			RETURNING id, email, email_verified, name, avatar_url, created_at, updated_at, last_login_at
		`, email, true, name, avatarURL, time.Now()).Scan(
			&user.ID, &user.Email, &user.EmailVerified, &user.Name,
			&user.AvatarURL, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create user: %w", err)
		}
	} else if err != nil {
		return nil, fmt.Errorf("failed to query user: %w", err)
	} else {
		// User exists, update last login and name/avatar if changed
		_, err = tx.Exec(ctx, `
			UPDATE auth.users
			SET last_login_at = $1, name = $2, avatar_url = $3
			WHERE id = $4
		`, time.Now(), name, avatarURL, user.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to update user: %w", err)
		}
	}

	// Now upsert the OAuth provider (no tokens stored - just for login tracking)
	_, err = tx.Exec(ctx, `
		INSERT INTO auth.user_oauth_providers (user_id, provider, provider_user_id, last_used_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (provider, provider_user_id)
		DO UPDATE SET
			last_used_at = EXCLUDED.last_used_at
	`, user.ID, provider, providerUserID, time.Now())
	if err != nil {
		return nil, fmt.Errorf("failed to upsert OAuth provider: %w", err)
	}

	// Commit transaction
	if err = tx.Commit(ctx); err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return &user, nil
}

// GetUserByID retrieves a user by their ID
func GetUserByID(ctx context.Context, userID uuid.UUID) (*User, error) {
	var user User
	err := database.DB.QueryRow(ctx, `
		SELECT id, email, email_verified, name, avatar_url, created_at, updated_at, last_login_at
		FROM auth.users
		WHERE id = $1
	`, userID).Scan(
		&user.ID, &user.Email, &user.EmailVerified, &user.Name,
		&user.AvatarURL, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}

// GetUserByEmail retrieves a user by their email
func GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := database.DB.QueryRow(ctx, `
		SELECT id, email, email_verified, name, avatar_url, created_at, updated_at, last_login_at
		FROM auth.users
		WHERE email = $1
	`, email).Scan(
		&user.ID, &user.Email, &user.EmailVerified, &user.Name,
		&user.AvatarURL, &user.CreatedAt, &user.UpdatedAt, &user.LastLoginAt,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	return &user, nil
}
