package budget

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// DashboardSummary represents the main dashboard overview
type DashboardSummary struct {
	TotalBalance           float64                  `json:"total_balance"`
	AccountCount           int                      `json:"account_count"`
	MonthToDateIncome      float64                  `json:"month_to_date_income"`
	MonthToDateExpenses    float64                  `json:"month_to_date_expenses"`
	MonthToDateNet         float64                  `json:"month_to_date_net"`
	BudgetedMonthlyIncome  float64                  `json:"budgeted_monthly_income"`
	BudgetedMonthlyExpense float64                  `json:"budgeted_monthly_expense"`
	BudgetHealthScore      int                      `json:"budget_health_score"`
	BudgetHealthStatus     string                   `json:"budget_health_status"`
	BudgetHealthMessage    string                   `json:"budget_health_message"`
	BudgetHealthColor      string                   `json:"budget_health_color"`
	UpcomingBills          []UpcomingBill           `json:"upcoming_bills"`
	RecentTransactions     []Transaction            `json:"recent_transactions"`
	SpendingByCategory     []CategorySpending       `json:"spending_by_category"`
}

// UpcomingBill represents a budget entry that's due soon
type UpcomingBill struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Amount      float64   `json:"amount"`
	DueDate     string    `json:"due_date"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	IsOverdue   bool      `json:"is_overdue"`
}

// CategorySpending represents spending grouped by category
type CategorySpending struct {
	CategoryID   *uuid.UUID `json:"category_id,omitempty"`
	CategoryName string     `json:"category_name"`
	TotalAmount  float64    `json:"total_amount"`
	Percentage   float64    `json:"percentage"`
	Color        *string    `json:"color,omitempty"`
}

// GetDashboardSummaryHandler returns comprehensive dashboard data
func GetDashboardSummaryHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get current month start/end dates
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Second)

	// Get all accounts for total balance
	accounts, err := GetAccountsByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve accounts",
		})
	}

	totalBalance := 0.0
	for _, account := range accounts {
		if account.IsActive {
			totalBalance += account.Balance
		}
	}

	// Get month-to-date transaction stats
	monthStartStr := monthStart.Format("2006-01-02")
	monthEndStr := monthEnd.Format("2006-01-02")

	monthIncome := 0.0
	monthExpenses := 0.0

	transactions, err := GetTransactionsByUserID(c.Context(), userID, nil, nil, &monthStartStr, &monthEndStr)
	if err == nil {
		for _, t := range transactions {
			if t.TransactionType == "income" {
				monthIncome += t.Amount
			} else {
				monthExpenses += t.Amount
			}
		}
	}

	// Get active budget summary
	var budgetedIncome float64
	var budgetedExpenses float64
	var budgetHealth int
	var healthStatus, healthMessage, healthColor string

	activeBudget, err := GetActiveBudget(c.Context(), userID)
	if err == nil && activeBudget != nil {
		summary, err := CalculateBudgetSummary(c.Context(), activeBudget.ID, userID)
		if err == nil {
			budgetedIncome = summary.TotalMonthlyIncome
			budgetedExpenses = summary.TotalMonthlyExpenses

			// Calculate health
			health := GetBudgetHealthStatus(summary)
			budgetHealth = health.Score
			healthStatus = health.Status
			healthMessage = health.Message
			healthColor = health.Color
		}
	}

	// Get upcoming bills (next 30 days from budget entries)
	upcomingBills := []UpcomingBill{}
	if activeBudget != nil {
		entries, err := GetBudgetEntriesByBudgetID(c.Context(), activeBudget.ID, userID)
		if err == nil {
			// Get upcoming expenses only
			for _, entry := range entries {
				if entry.EntryType == "expense" && entry.IsActive {
					// Simple logic: include all recurring entries
					// In a real implementation, calculate actual due dates based on frequency
					upcomingBills = append(upcomingBills, UpcomingBill{
						ID:         entry.ID,
						Name:       entry.Name,
						Amount:     entry.Amount,
						DueDate:    entry.StartDate,
						CategoryID: entry.CategoryID,
						IsOverdue:  false,
					})
				}
			}
		}
	}

	// Get recent transactions (last 10)
	recentTransactions, err := GetTransactionsByUserID(c.Context(), userID, nil, nil, nil, nil)
	if err != nil {
		recentTransactions = []Transaction{}
	}
	// Limit to last 10
	if len(recentTransactions) > 10 {
		recentTransactions = recentTransactions[:10]
	}

	// Get spending by category for this month
	spendingByCategory := []CategorySpending{}
	categories, err := GetCategoriesByUserID(c.Context(), userID)
	if err == nil && len(transactions) > 0 {
		categoryMap := make(map[string]*CategorySpending)

		for _, t := range transactions {
			if t.TransactionType == "expense" {
				var key string
				var categoryName string
				var categoryID *uuid.UUID
				var color *string

				if t.CategoryID != nil {
					key = t.CategoryID.String()
					categoryID = t.CategoryID

					// Find category name and color
					for _, cat := range categories {
						if cat.ID == *t.CategoryID {
							categoryName = cat.Name
							color = cat.Color
							break
						}
					}
				} else {
					key = "uncategorized"
					categoryName = "Uncategorized"
				}

				if _, exists := categoryMap[key]; !exists {
					categoryMap[key] = &CategorySpending{
						CategoryID:   categoryID,
						CategoryName: categoryName,
						TotalAmount:  0,
						Color:        color,
					}
				}
				categoryMap[key].TotalAmount += t.Amount
			}
		}

		// Calculate percentages
		for _, spending := range categoryMap {
			if monthExpenses > 0 {
				spending.Percentage = (spending.TotalAmount / monthExpenses) * 100
			}
			spendingByCategory = append(spendingByCategory, *spending)
		}
	}

	summary := DashboardSummary{
		TotalBalance:           totalBalance,
		AccountCount:           len(accounts),
		MonthToDateIncome:      monthIncome,
		MonthToDateExpenses:    monthExpenses,
		MonthToDateNet:         monthIncome - monthExpenses,
		BudgetedMonthlyIncome:  budgetedIncome,
		BudgetedMonthlyExpense: budgetedExpenses,
		BudgetHealthScore:      budgetHealth,
		BudgetHealthStatus:     healthStatus,
		BudgetHealthMessage:    healthMessage,
		BudgetHealthColor:      healthColor,
		UpcomingBills:          upcomingBills,
		RecentTransactions:     recentTransactions,
		SpendingByCategory:     spendingByCategory,
	}

	return c.JSON(summary)
}

// GetRecentActivityHandler returns recent transactions
func GetRecentActivityHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	limit := 20 // Default to last 20 transactions
	if limitParam := c.QueryInt("limit", 20); limitParam > 0 && limitParam <= 100 {
		limit = limitParam
	}

	transactions, err := GetTransactionsByUserID(c.Context(), userID, nil, nil, nil, nil)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve recent transactions",
		})
	}

	// Limit results
	if len(transactions) > limit {
		transactions = transactions[:limit]
	}

	return c.JSON(fiber.Map{
		"transactions": transactions,
	})
}

// GetSpendingByCategoryHandler returns spending breakdown by category
func GetSpendingByCategoryHandler(c *fiber.Ctx) error {
	userID := getUserIDFromContext(c)
	if userID == uuid.Nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	// Get date range from query params (default to current month)
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Second)

	startDate := c.Query("start_date", monthStart.Format("2006-01-02"))
	endDate := c.Query("end_date", monthEnd.Format("2006-01-02"))

	// Get transactions for the period
	transactions, err := GetTransactionsByUserID(c.Context(), userID, nil, nil, &startDate, &endDate)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve transactions",
		})
	}

	// Get categories for name lookup
	categories, err := GetCategoriesByUserID(c.Context(), userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve categories",
		})
	}

	// Group by category
	categoryMap := make(map[string]*CategorySpending)
	totalExpenses := 0.0

	for _, t := range transactions {
		if t.TransactionType == "expense" {
			totalExpenses += t.Amount

			var key string
			var categoryName string
			var categoryID *uuid.UUID
			var color *string

			if t.CategoryID != nil {
				key = t.CategoryID.String()
				categoryID = t.CategoryID

				for _, cat := range categories {
					if cat.ID == *t.CategoryID {
						categoryName = cat.Name
						color = cat.Color
						break
					}
				}
			} else {
				key = "uncategorized"
				categoryName = "Uncategorized"
			}

			if _, exists := categoryMap[key]; !exists {
				categoryMap[key] = &CategorySpending{
					CategoryID:   categoryID,
					CategoryName: categoryName,
					TotalAmount:  0,
					Color:        color,
				}
			}
			categoryMap[key].TotalAmount += t.Amount
		}
	}

	// Calculate percentages and convert to slice
	spendingByCategory := []CategorySpending{}
	for _, spending := range categoryMap {
		if totalExpenses > 0 {
			spending.Percentage = (spending.TotalAmount / totalExpenses) * 100
		}
		spendingByCategory = append(spendingByCategory, *spending)
	}

	return c.JSON(fiber.Map{
		"spending_by_category": spendingByCategory,
		"total_expenses":       totalExpenses,
		"start_date":           startDate,
		"end_date":             endDate,
	})
}
