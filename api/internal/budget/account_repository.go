package budget

import (
	"context"
	"fmt"

	"github.com/brendenbissett/help-me-budget/api/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GetAccountsByUserID retrieves all accounts for a specific user
func GetAccountsByUserID(ctx context.Context, userID uuid.UUID) ([]Account, error) {
	query := `
		SELECT id, user_id, name, account_type, balance, currency, is_active, created_at, updated_at
		FROM budget.accounts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := database.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query accounts: %w", err)
	}
	defer rows.Close()

	var accounts []Account
	for rows.Next() {
		var account Account
		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Name,
			&account.AccountType,
			&account.Balance,
			&account.Currency,
			&account.IsActive,
			&account.CreatedAt,
			&account.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan account: %w", err)
		}
		accounts = append(accounts, account)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating accounts: %w", err)
	}

	return accounts, nil
}

// GetAccountByID retrieves a specific account by ID
func GetAccountByID(ctx context.Context, accountID uuid.UUID, userID uuid.UUID) (*Account, error) {
	query := `
		SELECT id, user_id, name, account_type, balance, currency, is_active, created_at, updated_at
		FROM budget.accounts
		WHERE id = $1 AND user_id = $2
	`

	var account Account
	err := database.DB.QueryRow(ctx, query, accountID, userID).Scan(
		&account.ID,
		&account.UserID,
		&account.Name,
		&account.AccountType,
		&account.Balance,
		&account.Currency,
		&account.IsActive,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query account: %w", err)
	}

	return &account, nil
}

// CreateAccount creates a new account for a user
func CreateAccount(ctx context.Context, userID uuid.UUID, req CreateAccountRequest) (*Account, error) {
	query := `
		INSERT INTO budget.accounts (user_id, name, account_type, balance, currency)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, name, account_type, balance, currency, is_active, created_at, updated_at
	`

	var account Account
	err := database.DB.QueryRow(
		ctx,
		query,
		userID,
		req.Name,
		req.AccountType,
		req.Balance,
		req.Currency,
	).Scan(
		&account.ID,
		&account.UserID,
		&account.Name,
		&account.AccountType,
		&account.Balance,
		&account.Currency,
		&account.IsActive,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create account: %w", err)
	}

	return &account, nil
}

// UpdateAccount updates an existing account
func UpdateAccount(ctx context.Context, accountID uuid.UUID, userID uuid.UUID, req UpdateAccountRequest) (*Account, error) {
	// Build dynamic update query based on provided fields
	query := `UPDATE budget.accounts SET updated_at = CURRENT_TIMESTAMP`
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		query += fmt.Sprintf(", name = $%d", argIndex)
		args = append(args, *req.Name)
		argIndex++
	}
	if req.AccountType != nil {
		query += fmt.Sprintf(", account_type = $%d", argIndex)
		args = append(args, *req.AccountType)
		argIndex++
	}
	if req.Balance != nil {
		query += fmt.Sprintf(", balance = $%d", argIndex)
		args = append(args, *req.Balance)
		argIndex++
	}
	if req.Currency != nil {
		query += fmt.Sprintf(", currency = $%d", argIndex)
		args = append(args, *req.Currency)
		argIndex++
	}
	if req.IsActive != nil {
		query += fmt.Sprintf(", is_active = $%d", argIndex)
		args = append(args, *req.IsActive)
		argIndex++
	}

	query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d", argIndex, argIndex+1)
	args = append(args, accountID, userID)
	query += ` RETURNING id, user_id, name, account_type, balance, currency, is_active, created_at, updated_at`

	var account Account
	err := database.DB.QueryRow(ctx, query, args...).Scan(
		&account.ID,
		&account.UserID,
		&account.Name,
		&account.AccountType,
		&account.Balance,
		&account.Currency,
		&account.IsActive,
		&account.CreatedAt,
		&account.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("account not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update account: %w", err)
	}

	return &account, nil
}

// DeleteAccount deletes an account (soft delete by setting is_active = false)
func DeleteAccount(ctx context.Context, accountID uuid.UUID, userID uuid.UUID) error {
	query := `
		UPDATE budget.accounts
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2
	`

	result, err := database.DB.Exec(ctx, query, accountID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("account not found")
	}

	return nil
}

// GetTotalBalance calculates the total balance across all active accounts for a user
func GetTotalBalance(ctx context.Context, userID uuid.UUID) (float64, error) {
	query := `
		SELECT COALESCE(SUM(balance), 0)
		FROM budget.accounts
		WHERE user_id = $1 AND is_active = true
	`

	var totalBalance float64
	err := database.DB.QueryRow(ctx, query, userID).Scan(&totalBalance)
	if err != nil {
		return 0, fmt.Errorf("failed to calculate total balance: %w", err)
	}

	return totalBalance, nil
}
