package budget

import (
	"context"
	"fmt"
	"time"

	"github.com/brendenbissett/help-me-budget/api/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GetTransactionsByUserID retrieves all transactions for a user with optional filters
func GetTransactionsByUserID(ctx context.Context, userID uuid.UUID, accountID *uuid.UUID, categoryID *uuid.UUID, startDate *string, endDate *string) ([]Transaction, error) {
	query := `
		SELECT id, user_id, account_id, category_id, budget_entry_id, amount,
		       transaction_type, description, transaction_date::text, notes,
		       match_confidence, created_at, updated_at
		FROM budget.transactions
		WHERE user_id = $1
	`
	args := []interface{}{userID}
	argIndex := 2

	// Add optional filters
	if accountID != nil {
		query += fmt.Sprintf(" AND account_id = $%d", argIndex)
		args = append(args, *accountID)
		argIndex++
	}

	if categoryID != nil {
		query += fmt.Sprintf(" AND category_id = $%d", argIndex)
		args = append(args, *categoryID)
		argIndex++
	}

	if startDate != nil {
		query += fmt.Sprintf(" AND transaction_date >= $%d", argIndex)
		args = append(args, *startDate)
		argIndex++
	}

	if endDate != nil {
		query += fmt.Sprintf(" AND transaction_date <= $%d", argIndex)
		args = append(args, *endDate)
		argIndex++
	}

	query += " ORDER BY transaction_date DESC, created_at DESC"

	rows, err := database.DB.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query transactions: %w", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.AccountID,
			&t.CategoryID,
			&t.BudgetEntryID,
			&t.Amount,
			&t.TransactionType,
			&t.Description,
			&t.TransactionDate,
			&t.Notes,
			&t.MatchConfidence,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

// GetTransactionByID retrieves a specific transaction by ID
func GetTransactionByID(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID) (*Transaction, error) {
	query := `
		SELECT id, user_id, account_id, category_id, budget_entry_id, amount,
		       transaction_type, description, transaction_date::text, notes,
		       match_confidence, created_at, updated_at
		FROM budget.transactions
		WHERE id = $1 AND user_id = $2
	`

	var t Transaction
	err := database.DB.QueryRow(ctx, query, transactionID, userID).Scan(
		&t.ID,
		&t.UserID,
		&t.AccountID,
		&t.CategoryID,
		&t.BudgetEntryID,
		&t.Amount,
		&t.TransactionType,
		&t.Description,
		&t.TransactionDate,
		&t.Notes,
		&t.MatchConfidence,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("transaction not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return &t, nil
}

// CreateTransaction creates a new transaction
func CreateTransaction(ctx context.Context, userID uuid.UUID, req CreateTransactionRequest) (*Transaction, error) {
	query := `
		INSERT INTO budget.transactions
		(user_id, account_id, category_id, amount, transaction_type, description,
		 transaction_date, notes, match_confidence)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 'unmatched')
		RETURNING id, user_id, account_id, category_id, budget_entry_id, amount,
		          transaction_type, description, transaction_date::text, notes,
		          match_confidence, created_at, updated_at
	`

	var t Transaction
	err := database.DB.QueryRow(
		ctx,
		query,
		userID,
		req.AccountID,
		req.CategoryID,
		req.Amount,
		req.TransactionType,
		req.Description,
		req.TransactionDate,
		req.Notes,
	).Scan(
		&t.ID,
		&t.UserID,
		&t.AccountID,
		&t.CategoryID,
		&t.BudgetEntryID,
		&t.Amount,
		&t.TransactionType,
		&t.Description,
		&t.TransactionDate,
		&t.Notes,
		&t.MatchConfidence,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create transaction: %w", err)
	}

	return &t, nil
}

// UpdateTransaction updates an existing transaction
func UpdateTransaction(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID, req UpdateTransactionRequest) (*Transaction, error) {
	// Verify transaction belongs to user
	_, err := GetTransactionByID(ctx, transactionID, userID)
	if err != nil {
		return nil, err
	}

	// Build dynamic update query
	query := "UPDATE budget.transactions SET updated_at = $1"
	args := []interface{}{time.Now()}
	argIndex := 2

	if req.AccountID != nil {
		query += fmt.Sprintf(", account_id = $%d", argIndex)
		args = append(args, *req.AccountID)
		argIndex++
	}

	if req.CategoryID != nil {
		query += fmt.Sprintf(", category_id = $%d", argIndex)
		args = append(args, *req.CategoryID)
		argIndex++
	}

	if req.BudgetEntryID != nil {
		query += fmt.Sprintf(", budget_entry_id = $%d", argIndex)
		args = append(args, *req.BudgetEntryID)
		argIndex++
	}

	if req.Amount != nil {
		query += fmt.Sprintf(", amount = $%d", argIndex)
		args = append(args, *req.Amount)
		argIndex++
	}

	if req.TransactionType != nil {
		query += fmt.Sprintf(", transaction_type = $%d", argIndex)
		args = append(args, *req.TransactionType)
		argIndex++
	}

	if req.Description != nil {
		query += fmt.Sprintf(", description = $%d", argIndex)
		args = append(args, *req.Description)
		argIndex++
	}

	if req.TransactionDate != nil {
		query += fmt.Sprintf(", transaction_date = $%d", argIndex)
		args = append(args, *req.TransactionDate)
		argIndex++
	}

	if req.Notes != nil {
		query += fmt.Sprintf(", notes = $%d", argIndex)
		args = append(args, *req.Notes)
		argIndex++
	}

	if req.MatchConfidence != nil {
		query += fmt.Sprintf(", match_confidence = $%d", argIndex)
		args = append(args, *req.MatchConfidence)
		argIndex++
	}

	query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d", argIndex, argIndex+1)
	args = append(args, transactionID, userID)
	query += ` RETURNING id, user_id, account_id, category_id, budget_entry_id, amount,
	           transaction_type, description, transaction_date::text, notes,
	           match_confidence, created_at, updated_at`

	var t Transaction
	err = database.DB.QueryRow(ctx, query, args...).Scan(
		&t.ID,
		&t.UserID,
		&t.AccountID,
		&t.CategoryID,
		&t.BudgetEntryID,
		&t.Amount,
		&t.TransactionType,
		&t.Description,
		&t.TransactionDate,
		&t.Notes,
		&t.MatchConfidence,
		&t.CreatedAt,
		&t.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to update transaction: %w", err)
	}

	return &t, nil
}

// DeleteTransaction deletes a transaction (hard delete)
func DeleteTransaction(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID) error {
	query := `DELETE FROM budget.transactions WHERE id = $1 AND user_id = $2`

	result, err := database.DB.Exec(ctx, query, transactionID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete transaction: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}

// GetUnmatchedTransactions retrieves transactions that haven't been matched to budget entries
func GetUnmatchedTransactions(ctx context.Context, userID uuid.UUID) ([]Transaction, error) {
	query := `
		SELECT id, user_id, account_id, category_id, budget_entry_id, amount,
		       transaction_type, description, transaction_date::text, notes,
		       match_confidence, created_at, updated_at
		FROM budget.transactions
		WHERE user_id = $1 AND match_confidence = 'unmatched'
		ORDER BY transaction_date DESC, created_at DESC
	`

	rows, err := database.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query unmatched transactions: %w", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var t Transaction
		err := rows.Scan(
			&t.ID,
			&t.UserID,
			&t.AccountID,
			&t.CategoryID,
			&t.BudgetEntryID,
			&t.Amount,
			&t.TransactionType,
			&t.Description,
			&t.TransactionDate,
			&t.Notes,
			&t.MatchConfidence,
			&t.CreatedAt,
			&t.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan transaction: %w", err)
		}
		transactions = append(transactions, t)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating transactions: %w", err)
	}

	return transactions, nil
}

// CategorizeTransaction assigns a category to a transaction
func CategorizeTransaction(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID, categoryID uuid.UUID) (*Transaction, error) {
	req := UpdateTransactionRequest{
		CategoryID: &categoryID,
	}
	return UpdateTransaction(ctx, transactionID, userID, req)
}

// LinkTransactionToBudgetEntry links a transaction to a budget entry
func LinkTransactionToBudgetEntry(ctx context.Context, transactionID uuid.UUID, userID uuid.UUID, budgetEntryID uuid.UUID, confidence string) (*Transaction, error) {
	req := UpdateTransactionRequest{
		BudgetEntryID:   &budgetEntryID,
		MatchConfidence: &confidence,
	}
	return UpdateTransaction(ctx, transactionID, userID, req)
}
