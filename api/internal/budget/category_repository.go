package budget

import (
	"context"
	"fmt"

	"github.com/brendenbissett/help-me-budget/api/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// GetCategoriesByUserID retrieves all categories for a specific user
func GetCategoriesByUserID(ctx context.Context, userID uuid.UUID) ([]Category, error) {
	query := `
		SELECT id, user_id, name, category_type, color, icon, parent_category_id, is_active, created_at, updated_at
		FROM budget.categories
		WHERE user_id = $1
		ORDER BY category_type, name ASC
	`

	rows, err := database.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.CategoryType,
			&category.Color,
			&category.Icon,
			&category.ParentCategoryID,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}

// GetCategoryByID retrieves a specific category by ID
func GetCategoryByID(ctx context.Context, categoryID uuid.UUID, userID uuid.UUID) (*Category, error) {
	query := `
		SELECT id, user_id, name, category_type, color, icon, parent_category_id, is_active, created_at, updated_at
		FROM budget.categories
		WHERE id = $1 AND user_id = $2
	`

	var category Category
	err := database.DB.QueryRow(ctx, query, categoryID, userID).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.CategoryType,
		&category.Color,
		&category.Icon,
		&category.ParentCategoryID,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("category not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to query category: %w", err)
	}

	return &category, nil
}

// CreateCategory creates a new category for a user
func CreateCategory(ctx context.Context, userID uuid.UUID, req CreateCategoryRequest) (*Category, error) {
	query := `
		INSERT INTO budget.categories (user_id, name, category_type, color, icon, parent_category_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, name, category_type, color, icon, parent_category_id, is_active, created_at, updated_at
	`

	var category Category
	err := database.DB.QueryRow(
		ctx,
		query,
		userID,
		req.Name,
		req.CategoryType,
		req.Color,
		req.Icon,
		req.ParentCategoryID,
	).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.CategoryType,
		&category.Color,
		&category.Icon,
		&category.ParentCategoryID,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create category: %w", err)
	}

	return &category, nil
}

// UpdateCategory updates an existing category
func UpdateCategory(ctx context.Context, categoryID uuid.UUID, userID uuid.UUID, req UpdateCategoryRequest) (*Category, error) {
	// Build dynamic update query based on provided fields
	query := `UPDATE budget.categories SET updated_at = CURRENT_TIMESTAMP`
	args := []interface{}{}
	argIndex := 1

	if req.Name != nil {
		query += fmt.Sprintf(", name = $%d", argIndex)
		args = append(args, *req.Name)
		argIndex++
	}
	if req.CategoryType != nil {
		query += fmt.Sprintf(", category_type = $%d", argIndex)
		args = append(args, *req.CategoryType)
		argIndex++
	}
	if req.Color != nil {
		query += fmt.Sprintf(", color = $%d", argIndex)
		args = append(args, *req.Color)
		argIndex++
	}
	if req.Icon != nil {
		query += fmt.Sprintf(", icon = $%d", argIndex)
		args = append(args, *req.Icon)
		argIndex++
	}
	if req.ParentCategoryID != nil {
		query += fmt.Sprintf(", parent_category_id = $%d", argIndex)
		args = append(args, *req.ParentCategoryID)
		argIndex++
	}
	if req.IsActive != nil {
		query += fmt.Sprintf(", is_active = $%d", argIndex)
		args = append(args, *req.IsActive)
		argIndex++
	}

	query += fmt.Sprintf(" WHERE id = $%d AND user_id = $%d", argIndex, argIndex+1)
	args = append(args, categoryID, userID)
	query += ` RETURNING id, user_id, name, category_type, color, icon, parent_category_id, is_active, created_at, updated_at`

	var category Category
	err := database.DB.QueryRow(ctx, query, args...).Scan(
		&category.ID,
		&category.UserID,
		&category.Name,
		&category.CategoryType,
		&category.Color,
		&category.Icon,
		&category.ParentCategoryID,
		&category.IsActive,
		&category.CreatedAt,
		&category.UpdatedAt,
	)

	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("category not found")
	}
	if err != nil {
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return &category, nil
}

// DeleteCategory deletes a category (soft delete by setting is_active = false)
func DeleteCategory(ctx context.Context, categoryID uuid.UUID, userID uuid.UUID) error {
	query := `
		UPDATE budget.categories
		SET is_active = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2
	`

	result, err := database.DB.Exec(ctx, query, categoryID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	if result.RowsAffected() == 0 {
		return fmt.Errorf("category not found")
	}

	return nil
}

// SeedDefaultCategories creates a set of default categories for a new user
func SeedDefaultCategories(ctx context.Context, userID uuid.UUID) error {
	// Default expense categories
	expenseCategories := []struct {
		name  string
		color string
		icon  string
	}{
		{"Housing", "#FF6B6B", "home"},
		{"Transportation", "#4ECDC4", "car"},
		{"Food & Groceries", "#95E1D3", "shopping-cart"},
		{"Utilities", "#F38181", "zap"},
		{"Healthcare", "#AA96DA", "heart"},
		{"Entertainment", "#FCBAD3", "film"},
		{"Shopping", "#A8D8EA", "shopping-bag"},
		{"Personal Care", "#FFCFDF", "user"},
		{"Education", "#C7CEEA", "book"},
		{"Insurance", "#B5EAD7", "shield"},
		{"Subscriptions", "#FFD3B6", "repeat"},
		{"Dining Out", "#FFA8B6", "coffee"},
		{"Other Expenses", "#D4A5A5", "more-horizontal"},
	}

	// Default income categories
	incomeCategories := []struct {
		name  string
		color string
		icon  string
	}{
		{"Salary", "#52B788", "dollar-sign"},
		{"Freelance", "#74C69D", "briefcase"},
		{"Investments", "#95D5B2", "trending-up"},
		{"Gifts", "#B7E4C7", "gift"},
		{"Other Income", "#D8F3DC", "plus-circle"},
	}

	// Insert expense categories
	for _, cat := range expenseCategories {
		color := cat.color
		icon := cat.icon
		_, err := CreateCategory(ctx, userID, CreateCategoryRequest{
			Name:         cat.name,
			CategoryType: "expense",
			Color:        &color,
			Icon:         &icon,
		})
		if err != nil {
			return fmt.Errorf("failed to create expense category %s: %w", cat.name, err)
		}
	}

	// Insert income categories
	for _, cat := range incomeCategories {
		color := cat.color
		icon := cat.icon
		_, err := CreateCategory(ctx, userID, CreateCategoryRequest{
			Name:         cat.name,
			CategoryType: "income",
			Color:        &color,
			Icon:         &icon,
		})
		if err != nil {
			return fmt.Errorf("failed to create income category %s: %w", cat.name, err)
		}
	}

	return nil
}

// GetCategoriesByType retrieves categories filtered by type (income or expense)
func GetCategoriesByType(ctx context.Context, userID uuid.UUID, categoryType string) ([]Category, error) {
	query := `
		SELECT id, user_id, name, category_type, color, icon, parent_category_id, is_active, created_at, updated_at
		FROM budget.categories
		WHERE user_id = $1 AND category_type = $2 AND is_active = true
		ORDER BY name ASC
	`

	rows, err := database.DB.Query(ctx, query, userID, categoryType)
	if err != nil {
		return nil, fmt.Errorf("failed to query categories by type: %w", err)
	}
	defer rows.Close()

	var categories []Category
	for rows.Next() {
		var category Category
		err := rows.Scan(
			&category.ID,
			&category.UserID,
			&category.Name,
			&category.CategoryType,
			&category.Color,
			&category.Icon,
			&category.ParentCategoryID,
			&category.IsActive,
			&category.CreatedAt,
			&category.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan category: %w", err)
		}
		categories = append(categories, category)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating categories: %w", err)
	}

	return categories, nil
}
