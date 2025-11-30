package budget

import (
	"context"
	"fmt"
	"time"

	"github.com/brendenbissett/help-me-budget/api/internal/database"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ============================================================================
// Report Types
// ============================================================================

// SpendingTrend represents spending for a category over time
type SpendingTrend struct {
	Month      string  `json:"month"`       // YYYY-MM format
	CategoryID string  `json:"category_id"`
	Category   string  `json:"category"`
	Amount     float64 `json:"amount"`
}

// BudgetVariance shows budget vs actual for each entry
type BudgetVariance struct {
	EntryID     uuid.UUID `json:"entry_id"`
	EntryName   string    `json:"entry_name"`
	Category    string    `json:"category"`
	Budgeted    float64   `json:"budgeted"`
	Actual      float64   `json:"actual"`
	Variance    float64   `json:"variance"` // positive = under budget, negative = over budget
	VariancePct float64   `json:"variance_pct"`
}

// DailyCashFlowProjection represents projected balance for a single day
type DailyCashFlowProjection struct {
	Date              string  `json:"date"`
	ProjectedIncome   float64 `json:"projected_income"`
	ProjectedExpenses float64 `json:"projected_expenses"`
	ProjectedBalance  float64 `json:"projected_balance"`
}

// TopExpense represents a high-spending category
type TopExpense struct {
	CategoryID   string  `json:"category_id"`
	CategoryName string  `json:"category_name"`
	TotalAmount  float64 `json:"total_amount"`
	Percentage   float64 `json:"percentage"` // Percentage of total expenses
	Count        int     `json:"count"`      // Number of transactions
}

// ============================================================================
// Handler Functions
// ============================================================================

// GetSpendingTrendsHandler returns spending trends by category over time
func GetSpendingTrendsHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get optional date range from query params
	startDate := c.Query("start_date") // YYYY-MM-DD
	endDate := c.Query("end_date")     // YYYY-MM-DD

	// Default to last 6 months if not provided
	if startDate == "" {
		startDate = time.Now().AddDate(0, -6, 0).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	trends, err := GetSpendingTrends(c.Context(), userID, startDate, endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get spending trends",
		})
	}

	return c.JSON(trends)
}

// GetBudgetVarianceHandler returns budget vs actual comparison
func GetBudgetVarianceHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get optional month from query params (defaults to current month)
	month := c.Query("month") // YYYY-MM format
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	variance, err := GetBudgetVariance(c.Context(), userID, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get budget variance",
		})
	}

	return c.JSON(variance)
}

// GetCashFlowProjectionHandler returns projected cash flow
func GetCashFlowProjectionHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get optional parameters
	daysStr := c.Query("days", "90")               // Default 90 days
	startingBalanceStr := c.Query("starting_balance", "0")

	var days int
	var startingBalance float64
	fmt.Sscanf(daysStr, "%d", &days)
	fmt.Sscanf(startingBalanceStr, "%f", &startingBalance)

	if days <= 0 || days > 365 {
		days = 90
	}

	projection, err := GetCashFlowProjection(c.Context(), userID, days, startingBalance)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get cash flow projection",
		})
	}

	return c.JSON(projection)
}

// GetTopExpensesHandler returns the biggest spending categories
func GetTopExpensesHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get optional date range and limit
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")
	limitStr := c.Query("limit", "10")

	// Default to current month if not provided
	if startDate == "" {
		now := time.Now()
		startDate = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Format("2006-01-02")
	}
	if endDate == "" {
		endDate = time.Now().Format("2006-01-02")
	}

	var limit int
	fmt.Sscanf(limitStr, "%d", &limit)
	if limit <= 0 || limit > 50 {
		limit = 10
	}

	topExpenses, err := GetTopExpenses(c.Context(), userID, startDate, endDate, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to get top expenses",
		})
	}

	return c.JSON(topExpenses)
}

// ============================================================================
// Repository Functions
// ============================================================================

// GetSpendingTrends calculates spending by category over time
func GetSpendingTrends(ctx context.Context, userID uuid.UUID, startDate, endDate string) ([]SpendingTrend, error) {
	query := `
		SELECT
			TO_CHAR(t.transaction_date::date, 'YYYY-MM') as month,
			COALESCE(t.category_id::text, 'uncategorized') as category_id,
			COALESCE(c.name, 'Uncategorized') as category,
			SUM(t.amount) as amount
		FROM budget.transactions t
		LEFT JOIN budget.categories c ON t.category_id = c.id
		WHERE t.user_id = $1
			AND t.transaction_type = 'expense'
			AND t.transaction_date >= $2
			AND t.transaction_date <= $3
			AND t.deleted_at IS NULL
		GROUP BY month, category_id, category
		ORDER BY month DESC, amount DESC
	`

	rows, err := database.DB.Query(ctx, query, userID, startDate, endDate)
	if err != nil {
		return nil, fmt.Errorf("failed to query spending trends: %w", err)
	}
	defer rows.Close()

	var trends []SpendingTrend
	for rows.Next() {
		var trend SpendingTrend
		if err := rows.Scan(&trend.Month, &trend.CategoryID, &trend.Category, &trend.Amount); err != nil {
			return nil, fmt.Errorf("failed to scan spending trend: %w", err)
		}
		trends = append(trends, trend)
	}

	return trends, nil
}

// GetBudgetVariance compares budget vs actual for a given month
func GetBudgetVariance(ctx context.Context, userID uuid.UUID, month string) ([]BudgetVariance, error) {
	// Get active budget
	activeBudget, err := GetActiveBudget(ctx, userID)
	if err != nil || activeBudget == nil {
		return []BudgetVariance{}, nil // No active budget
	}

	// Get budget entries
	entries, err := GetBudgetEntriesByBudgetID(ctx, activeBudget.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget entries: %w", err)
	}

	// Parse month (YYYY-MM)
	startOfMonth, err := time.Parse("2006-01", month)
	if err != nil {
		return nil, fmt.Errorf("invalid month format: %w", err)
	}
	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Second)

	// Calculate variance for each entry
	var variances []BudgetVariance
	for _, entry := range entries {
		// Get actual spending for this entry
		actual, err := getActualForEntry(ctx, userID, entry.ID, startOfMonth, endOfMonth)
		if err != nil {
			continue // Skip on error
		}

		categoryName := "Uncategorized"
		if entry.CategoryID != nil {
			category, _ := GetCategoryByID(ctx, *entry.CategoryID, userID)
			if category != nil {
				categoryName = category.Name
			}
		}

		variance := entry.Amount - actual
		variancePct := 0.0
		if entry.Amount > 0 {
			variancePct = (variance / entry.Amount) * 100
		}

		variances = append(variances, BudgetVariance{
			EntryID:     entry.ID,
			EntryName:   entry.Name,
			Category:    categoryName,
			Budgeted:    entry.Amount,
			Actual:      actual,
			Variance:    variance,
			VariancePct: variancePct,
		})
	}

	return variances, nil
}

// getActualForEntry gets total actual spending/income for a budget entry in a date range
func getActualForEntry(ctx context.Context, userID, entryID uuid.UUID, startDate, endDate time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(amount), 0)
		FROM budget.transactions
		WHERE user_id = $1
			AND budget_entry_id = $2
			AND transaction_date >= $3
			AND transaction_date <= $4
			AND deleted_at IS NULL
	`

	var total float64
	err := database.DB.QueryRow(ctx, query, userID, entryID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02")).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

// GetCashFlowProjection projects future balance based on budget entries
func GetCashFlowProjection(ctx context.Context, userID uuid.UUID, days int, startingBalance float64) ([]DailyCashFlowProjection, error) {
	// Get active budget
	activeBudget, err := GetActiveBudget(ctx, userID)
	if err != nil || activeBudget == nil {
		return []DailyCashFlowProjection{}, nil
	}

	// Get budget entries
	entries, err := GetBudgetEntriesByBudgetID(ctx, activeBudget.ID, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get budget entries: %w", err)
	}

	// Project daily balance
	projections := []DailyCashFlowProjection{}
	currentBalance := startingBalance
	today := time.Now()

	for i := 0; i < days; i++ {
		date := today.AddDate(0, 0, i)
		dailyIncome := 0.0
		dailyExpenses := 0.0

		// Calculate expected income/expenses for this day based on budget entries
		for _, entry := range entries {
			if shouldOccurOnDate(entry, date) {
				if entry.EntryType == "income" {
					dailyIncome += entry.Amount
				} else {
					dailyExpenses += entry.Amount
				}
			}
		}

		currentBalance += dailyIncome - dailyExpenses

		projections = append(projections, DailyCashFlowProjection{
			Date:              date.Format("2006-01-02"),
			ProjectedIncome:   dailyIncome,
			ProjectedExpenses: dailyExpenses,
			ProjectedBalance:  currentBalance,
		})
	}

	return projections, nil
}

// shouldOccurOnDate checks if a budget entry should occur on a given date
func shouldOccurOnDate(entry BudgetEntry, date time.Time) bool {
	startDate, err := time.Parse("2006-01-02", entry.StartDate)
	if err != nil || date.Before(startDate) {
		return false
	}

	// Check end date if set
	if entry.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *entry.EndDate)
		if err == nil && date.After(endDate) {
			return false
		}
	}

	switch entry.Frequency {
	case "once_off":
		return date.Format("2006-01-02") == entry.StartDate

	case "daily":
		return true

	case "weekly":
		if entry.DayOfWeek != nil {
			return int(date.Weekday()) == *entry.DayOfWeek
		}
		return false

	case "fortnightly":
		daysSince := date.Sub(startDate).Hours() / 24
		weeksSince := int(daysSince / 7)
		if weeksSince%2 == 0 && entry.DayOfWeek != nil {
			return int(date.Weekday()) == *entry.DayOfWeek
		}
		return false

	case "monthly":
		if entry.DayOfMonth != nil {
			return date.Day() == *entry.DayOfMonth
		}
		return false

	case "annually":
		return date.Month() == startDate.Month() && date.Day() == startDate.Day()

	default:
		return false
	}
}

// GetTopExpenses returns the highest spending categories
func GetTopExpenses(ctx context.Context, userID uuid.UUID, startDate, endDate string, limit int) ([]TopExpense, error) {
	// First get total expenses for percentage calculation
	var totalExpenses float64
	totalQuery := `
		SELECT COALESCE(SUM(amount), 0)
		FROM budget.transactions
		WHERE user_id = $1
			AND transaction_type = 'expense'
			AND transaction_date >= $2
			AND transaction_date <= $3
			AND deleted_at IS NULL
	`
	err := database.DB.QueryRow(ctx, totalQuery, userID, startDate, endDate).Scan(&totalExpenses)
	if err != nil {
		return nil, fmt.Errorf("failed to get total expenses: %w", err)
	}

	// Get top expenses by category
	query := `
		SELECT
			COALESCE(t.category_id::text, 'uncategorized') as category_id,
			COALESCE(c.name, 'Uncategorized') as category_name,
			SUM(t.amount) as total_amount,
			COUNT(*) as count
		FROM budget.transactions t
		LEFT JOIN budget.categories c ON t.category_id = c.id
		WHERE t.user_id = $1
			AND t.transaction_type = 'expense'
			AND t.transaction_date >= $2
			AND t.transaction_date <= $3
			AND t.deleted_at IS NULL
		GROUP BY category_id, category_name
		ORDER BY total_amount DESC
		LIMIT $4
	`

	rows, err := database.DB.Query(ctx, query, userID, startDate, endDate, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to query top expenses: %w", err)
	}
	defer rows.Close()

	var topExpenses []TopExpense
	for rows.Next() {
		var expense TopExpense
		if err := rows.Scan(&expense.CategoryID, &expense.CategoryName, &expense.TotalAmount, &expense.Count); err != nil {
			return nil, fmt.Errorf("failed to scan top expense: %w", err)
		}

		// Calculate percentage
		if totalExpenses > 0 {
			expense.Percentage = (expense.TotalAmount / totalExpenses) * 100
		}

		topExpenses = append(topExpenses, expense)
	}

	return topExpenses, nil
}
