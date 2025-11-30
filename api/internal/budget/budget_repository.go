package budget

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/brendenbissett/help-me-budget/api/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GetBudgetsByUserID retrieves all budgets for a specific user
func GetBudgetsByUserID(ctx context.Context, userID uuid.UUID) ([]Budget, error) {
	query := `
		SELECT id, user_id, name, description, is_active, created_at, updated_at
		FROM budget.budgets
		WHERE user_id = $1
		ORDER BY is_active DESC, created_at DESC
	`

	rows, err := database.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query budgets: %w", err)
	}
	defer rows.Close()

	var budgets []Budget
	for rows.Next() {
		var budget Budget
		err := rows.Scan(
			&budget.ID,
			&budget.UserID,
			&budget.Name,
			&budget.Description,
			&budget.IsActive,
			&budget.CreatedAt,
			&budget.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget: %w", err)
		}
		budgets = append(budgets, budget)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating budgets: %w", err)
	}

	return budgets, nil
}

// GetBudgetByID retrieves a specific budget by ID
func GetBudgetByID(ctx context.Context, budgetID uuid.UUID, userID uuid.UUID) (*Budget, error) {
	query := `
		SELECT id, user_id, name, description, is_active, created_at, updated_at
		FROM budget.budgets
		WHERE id = $1 AND user_id = $2
	`

	var budget Budget
	err := database.DB.QueryRow(ctx, query, budgetID, userID).Scan(
		&budget.ID,
		&budget.UserID,
		&budget.Name,
		&budget.Description,
		&budget.IsActive,
		&budget.CreatedAt,
		&budget.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("budget not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query budget: %w", err)
	}

	return &budget, nil
}

// CreateBudget creates a new budget for a user
func CreateBudget(ctx context.Context, userID uuid.UUID, req CreateBudgetRequest) (*Budget, error) {
	query := `
		INSERT INTO budget.budgets (user_id, name, description)
		VALUES ($1, $2, $3)
		RETURNING id, user_id, name, description, is_active, created_at, updated_at
	`

	var budget Budget
	err := database.DB.QueryRow(
		ctx,
		query,
		userID,
		req.Name,
		req.Description,
	).Scan(
		&budget.ID,
		&budget.UserID,
		&budget.Name,
		&budget.Description,
		&budget.IsActive,
		&budget.CreatedAt,
		&budget.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create budget: %w", err)
	}

	return &budget, nil
}

// UpdateBudget updates an existing budget
func UpdateBudget(ctx context.Context, budgetID uuid.UUID, userID uuid.UUID, req UpdateBudgetRequest) (*Budget, error) {
	query := `UPDATE budget.budgets SET updated_at = CURRENT_TIMESTAMP`
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		query += fmt.Sprintf(", name = $%d", argIndex)
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Description != nil {
		query += fmt.Sprintf(", description = $%d", argIndex)
		args = append(args, *req.Description)
		argIndex++
	}
	if req.IsActive != nil {
		query += fmt.Sprintf(", is_active = $%d", argIndex)
		args = append(args, *req.IsActive)
		argIndex++
	}

	query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d", argIndex, argIndex+1)
	args = append(args, budgetID, userID)
	query += ` RETURNING id, user_id, name, description, is_active, created_at, updated_at`

	var budget Budget
	err := database.DB.QueryRow(ctx, query, args...).Scan(
		&budget.ID,
		&budget.UserID,
		&budget.Name,
		&budget.Description,
		&budget.IsActive,
		&budget.CreatedAt,
		&budget.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("budget not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update budget: %w", err)
	}

	return &budget, nil
}

// DeleteBudget deletes a budget (soft delete by setting is_active = false)
func DeleteBudget(ctx context.Context, budgetID uuid.UUID, userID uuid.UUID) error {
	query := `
		UPDATE budget.budgets
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2
	`

	result, err := database.DB.Exec(ctx, query, budgetID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete budget: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("budget not found")
	}

	return nil
}

// GetBudgetWithEntries retrieves a budget with all its entries
func GetBudgetWithEntries(ctx context.Context, budgetID uuid.UUID, userID uuid.UUID) (*BudgetWithEntries, error) {
	// First get the budget
	budget, err := GetBudgetByID(ctx, budgetID, userID)
	if err != nil {
		return nil, err
	}

	// Then get all entries for this budget
	entries, err := GetBudgetEntriesByBudgetID(ctx, budgetID, userID)
	if err != nil {
		return nil, err
	}

	return &BudgetWithEntries{
		Budget:  *budget,
		Entries: entries,
	}, nil
}

// GetActiveBudget retrieves the user's active budget (if any)
func GetActiveBudget(ctx context.Context, userID uuid.UUID) (*Budget, error) {
	query := `
		SELECT id, user_id, name, description, is_active, created_at, updated_at
		FROM budget.budgets
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
		LIMIT 1
	`

	var budget Budget
	err := database.DB.QueryRow(ctx, query, userID).Scan(
		&budget.ID,
		&budget.UserID,
		&budget.Name,
		&budget.Description,
		&budget.IsActive,
		&budget.CreatedAt,
		&budget.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, nil // No active budget is not an error
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query active budget: %w", err)
	}

	return &budget, nil
}

// GetBudgetEntriesByBudgetID retrieves all budget entries for a specific budget
func GetBudgetEntriesByBudgetID(ctx context.Context, budgetID uuid.UUID, userID uuid.UUID) ([]BudgetEntry, error) {
	// First verify the budget belongs to the user
	_, err := GetBudgetByID(ctx, budgetID, userID)
	if err != nil {
		return nil, err
	}

	query := `
		SELECT id, budget_id, category_id, name, description, amount, entry_type,
		       frequency, day_of_month, day_of_week, start_date::text, end_date::text,
		       matching_rules, is_active, created_at, updated_at
		FROM budget.budget_entries
		WHERE budget_id = $1 AND is_active = true
		ORDER BY entry_type DESC, amount DESC
	`

	rows, err := database.DB.Query(ctx, query, budgetID)
	if err != nil {
		return nil, fmt.Errorf("failed to query budget entries: %w", err)
	}
	defer rows.Close()

	var entries []BudgetEntry
	for rows.Next() {
		var entry BudgetEntry
		var matchingRulesJSON []byte

		err := rows.Scan(
			&entry.ID,
			&entry.BudgetID,
			&entry.CategoryID,
			&entry.Name,
			&entry.Description,
			&entry.Amount,
			&entry.EntryType,
			&entry.Frequency,
			&entry.DayOfMonth,
			&entry.DayOfWeek,
			&entry.StartDate,
			&entry.EndDate,
			&matchingRulesJSON,
			&entry.IsActive,
			&entry.CreatedAt,
			&entry.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan budget entry: %w", err)
		}

		// Parse matching rules JSON
		if matchingRulesJSON != nil {
			err = json.Unmarshal(matchingRulesJSON, &entry.MatchingRules)
			if err != nil {
				return nil, fmt.Errorf("failed to unmarshal matching rules: %w", err)
			}
		}

		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating budget entries: %w", err)
	}

	return entries, nil
}

// CreateBudgetEntry creates a new budget entry
func CreateBudgetEntry(ctx context.Context, budgetID uuid.UUID, userID uuid.UUID, req CreateBudgetEntryRequest) (*BudgetEntry, error) {
	// Verify the budget belongs to the user
	_, err := GetBudgetByID(ctx, budgetID, userID)
	if err != nil {
		return nil, err
	}

	// Convert matching rules to JSON
	var matchingRulesJSON []byte
	if req.MatchingRules != nil {
		matchingRulesJSON, err = json.Marshal(req.MatchingRules)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal matching rules: %w", err)
		}
	}

	query := `
		INSERT INTO budget.budget_entries
		(budget_id, category_id, name, description, amount, entry_type, frequency,
		 day_of_month, day_of_week, start_date, end_date, matching_rules)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id, budget_id, category_id, name, description, amount, entry_type,
		          frequency, day_of_month, day_of_week, start_date::text, end_date::text,
		          matching_rules, is_active, created_at, updated_at
	`

	var entry BudgetEntry
	var matchingRulesResult []byte

	err = database.DB.QueryRow(
		ctx,
		query,
		budgetID,
		req.CategoryID,
		req.Name,
		req.Description,
		req.Amount,
		req.EntryType,
		req.Frequency,
		req.DayOfMonth,
		req.DayOfWeek,
		req.StartDate,
		req.EndDate,
		matchingRulesJSON,
	).Scan(
		&entry.ID,
		&entry.BudgetID,
		&entry.CategoryID,
		&entry.Name,
		&entry.Description,
		&entry.Amount,
		&entry.EntryType,
		&entry.Frequency,
		&entry.DayOfMonth,
		&entry.DayOfWeek,
		&entry.StartDate,
		&entry.EndDate,
		&matchingRulesResult,
		&entry.IsActive,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create budget entry: %w", err)
	}

	// Parse matching rules JSON
	if matchingRulesResult != nil {
		err = json.Unmarshal(matchingRulesResult, &entry.MatchingRules)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal matching rules: %w", err)
		}
	}

	return &entry, nil
}

// UpdateBudgetEntry updates an existing budget entry
func UpdateBudgetEntry(ctx context.Context, entryID uuid.UUID, budgetID uuid.UUID, userID uuid.UUID, req UpdateBudgetEntryRequest) (*BudgetEntry, error) {
	// Verify the budget belongs to the user
	_, err := GetBudgetByID(ctx, budgetID, userID)
	if err != nil {
		return nil, err
	}

	query := `UPDATE budget.budget_entries SET updated_at = CURRENT_TIMESTAMP`
	args := []interface{}{}
	argIndex := 1

	if req.CategoryID != nil {
		query += fmt.Sprintf(", category_id = $%d", argIndex)
		args = append(args, *req.CategoryID)
		argIndex++
	}
	if req.Name != nil {
		query += fmt.Sprintf(", name = $%d", argIndex)
		args = append(args, *req.Name)
		argIndex++
	}
	if req.Description != nil {
		query += fmt.Sprintf(", description = $%d", argIndex)
		args = append(args, *req.Description)
		argIndex++
	}
	if req.Amount != nil {
		query += fmt.Sprintf(", amount = $%d", argIndex)
		args = append(args, *req.Amount)
		argIndex++
	}
	if req.EntryType != nil {
		query += fmt.Sprintf(", entry_type = $%d", argIndex)
		args = append(args, *req.EntryType)
		argIndex++
	}
	if req.Frequency != nil {
		query += fmt.Sprintf(", frequency = $%d", argIndex)
		args = append(args, *req.Frequency)
		argIndex++
	}
	if req.DayOfMonth != nil {
		query += fmt.Sprintf(", day_of_month = $%d", argIndex)
		args = append(args, *req.DayOfMonth)
		argIndex++
	}
	if req.DayOfWeek != nil {
		query += fmt.Sprintf(", day_of_week = $%d", argIndex)
		args = append(args, *req.DayOfWeek)
		argIndex++
	}
	if req.StartDate != nil {
		query += fmt.Sprintf(", start_date = $%d", argIndex)
		args = append(args, *req.StartDate)
		argIndex++
	}
	if req.EndDate != nil {
		query += fmt.Sprintf(", end_date = $%d", argIndex)
		args = append(args, *req.EndDate)
		argIndex++
	}
	if req.MatchingRules != nil {
		matchingRulesJSON, err := json.Marshal(req.MatchingRules)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal matching rules: %w", err)
		}
		query += fmt.Sprintf(", matching_rules = $%d", argIndex)
		args = append(args, matchingRulesJSON)
		argIndex++
	}
	if req.IsActive != nil {
		query += fmt.Sprintf(", is_active = $%d", argIndex)
		args = append(args, *req.IsActive)
		argIndex++
	}

	query += fmt.Sprintf(" WHERE id = $%d AND budget_id = $%d", argIndex, argIndex+1)
	args = append(args, entryID, budgetID)
	query += ` RETURNING id, budget_id, category_id, name, description, amount, entry_type,
	           frequency, day_of_month, day_of_week, start_date::text, end_date::text,
	           matching_rules, is_active, created_at, updated_at`

	var entry BudgetEntry
	var matchingRulesResult []byte

	err = database.DB.QueryRow(ctx, query, args...).Scan(
		&entry.ID,
		&entry.BudgetID,
		&entry.CategoryID,
		&entry.Name,
		&entry.Description,
		&entry.Amount,
		&entry.EntryType,
		&entry.Frequency,
		&entry.DayOfMonth,
		&entry.DayOfWeek,
		&entry.StartDate,
		&entry.EndDate,
		&matchingRulesResult,
		&entry.IsActive,
		&entry.CreatedAt,
		&entry.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("budget entry not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update budget entry: %w", err)
	}

	// Parse matching rules JSON
	if matchingRulesResult != nil {
		err = json.Unmarshal(matchingRulesResult, &entry.MatchingRules)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal matching rules: %w", err)
		}
	}

	return &entry, nil
}

// DeleteBudgetEntry deletes a budget entry (soft delete)
func DeleteBudgetEntry(ctx context.Context, entryID uuid.UUID, budgetID uuid.UUID, userID uuid.UUID) error {
	// Verify the budget belongs to the user
	_, err := GetBudgetByID(ctx, budgetID, userID)
	if err != nil {
		return err
	}

	query := `
		UPDATE budget.budget_entries
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND budget_id = $2
	`

	result, err := database.DB.Exec(ctx, query, entryID, budgetID)
	if err != nil {
		return fmt.Errorf("failed to delete budget entry: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("budget entry not found")
	}

	return nil
}
