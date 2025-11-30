package budget

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// CalculateBudgetSummary calculates summary statistics for a budget
func CalculateBudgetSummary(ctx context.Context, budgetID uuid.UUID, userID uuid.UUID) (*BudgetSummary, error) {
	// Get all budget entries
	entries, err := GetBudgetEntriesByBudgetID(ctx, budgetID, userID)
	if err != nil {
		return nil, err
	}

	summary := &BudgetSummary{
		BudgetID: budgetID,
	}

	for _, entry := range entries {
		monthlyAmount := calculateMonthlyAmount(entry.Amount, entry.Frequency)
		annualAmount := monthlyAmount * 12

		if entry.EntryType == "income" {
			summary.TotalMonthlyIncome += monthlyAmount
			summary.TotalAnnualIncome += annualAmount
			summary.IncomeEntriesCount++
		} else {
			summary.TotalMonthlyExpenses += monthlyAmount
			summary.TotalAnnualExpenses += annualAmount
			summary.ExpenseEntriesCount++
		}
	}

	summary.MonthlySurplusDeficit = summary.TotalMonthlyIncome - summary.TotalMonthlyExpenses
	summary.AnnualSurplusDeficit = summary.TotalAnnualIncome - summary.TotalAnnualExpenses

	return summary, nil
}

// calculateMonthlyAmount converts any frequency to a monthly equivalent
func calculateMonthlyAmount(amount float64, frequency string) float64 {
	switch frequency {
	case "once_off":
		return 0 // One-time expenses don't count toward recurring monthly budget
	case "daily":
		return amount * 30.44 // Average days per month
	case "weekly":
		return amount * 4.33 // Average weeks per month
	case "fortnightly":
		return amount * 2.17 // Average fortnights per month
	case "monthly":
		return amount
	case "annually":
		return amount / 12
	default:
		return 0
	}
}

// CashFlowProjection represents projected cash flow over time
type CashFlowProjection struct {
	StartDate         string                  `json:"start_date"`
	EndDate           string                  `json:"end_date"`
	StartingBalance   float64                 `json:"starting_balance"`
	EndingBalance     float64                 `json:"ending_balance"`
	TotalIncome       float64                 `json:"total_income"`
	TotalExpenses     float64                 `json:"total_expenses"`
	NetCashFlow       float64                 `json:"net_cash_flow"`
	DailyProjections  []DailyProjection       `json:"daily_projections"`
	MonthlyBreakdown  []MonthlyBreakdown      `json:"monthly_breakdown"`
}

// DailyProjection represents projected balance for a specific day
type DailyProjection struct {
	Date            string  `json:"date"`
	Balance         float64 `json:"balance"`
	DailyIncome     float64 `json:"daily_income"`
	DailyExpenses   float64 `json:"daily_expenses"`
	DailyNet        float64 `json:"daily_net"`
}

// MonthlyBreakdown represents monthly summary
type MonthlyBreakdown struct {
	Month         string  `json:"month"` // "2025-01"
	Income        float64 `json:"income"`
	Expenses      float64 `json:"expenses"`
	Net           float64 `json:"net"`
	EndingBalance float64 `json:"ending_balance"`
}

// ProjectCashFlow projects future cash flow based on budget entries
func ProjectCashFlow(ctx context.Context, budgetID uuid.UUID, userID uuid.UUID, startingBalance float64, days int) (*CashFlowProjection, error) {
	// Get all budget entries
	entries, err := GetBudgetEntriesByBudgetID(ctx, budgetID, userID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	startDate := now.Format("2006-01-02")
	endDate := now.AddDate(0, 0, days).Format("2006-01-02")

	projection := &CashFlowProjection{
		StartDate:       startDate,
		EndDate:         endDate,
		StartingBalance: startingBalance,
		DailyProjections: make([]DailyProjection, 0),
		MonthlyBreakdown: make([]MonthlyBreakdown, 0),
	}

	currentBalance := startingBalance
	monthlyTotals := make(map[string]*MonthlyBreakdown)

	// Iterate through each day in the projection period
	for i := 0; i < days; i++ {
		currentDate := now.AddDate(0, 0, i)
		dateStr := currentDate.Format("2006-01-02")
		monthStr := currentDate.Format("2006-01")

		dailyIncome := 0.0
		dailyExpenses := 0.0

		// Check each budget entry to see if it occurs on this date
		for _, entry := range entries {
			if shouldEntryOccurOnDate(entry, currentDate) {
				if entry.EntryType == "income" {
					dailyIncome += entry.Amount
				} else {
					dailyExpenses += entry.Amount
				}
			}
		}

		dailyNet := dailyIncome - dailyExpenses
		currentBalance += dailyNet

		// Add to daily projections
		projection.DailyProjections = append(projection.DailyProjections, DailyProjection{
			Date:          dateStr,
			Balance:       currentBalance,
			DailyIncome:   dailyIncome,
			DailyExpenses: dailyExpenses,
			DailyNet:      dailyNet,
		})

		// Update monthly totals
		if monthlyTotals[monthStr] == nil {
			monthlyTotals[monthStr] = &MonthlyBreakdown{
				Month: monthStr,
			}
		}
		monthlyTotals[monthStr].Income += dailyIncome
		monthlyTotals[monthStr].Expenses += dailyExpenses
		monthlyTotals[monthStr].Net += dailyNet
		monthlyTotals[monthStr].EndingBalance = currentBalance

		// Update projection totals
		projection.TotalIncome += dailyIncome
		projection.TotalExpenses += dailyExpenses
	}

	projection.EndingBalance = currentBalance
	projection.NetCashFlow = projection.TotalIncome - projection.TotalExpenses

	// Convert monthly totals to sorted array
	for _, breakdown := range monthlyTotals {
		projection.MonthlyBreakdown = append(projection.MonthlyBreakdown, *breakdown)
	}

	return projection, nil
}

// shouldEntryOccurOnDate determines if a budget entry should occur on a specific date
func shouldEntryOccurOnDate(entry BudgetEntry, date time.Time) bool {
	// Parse start date
	startDate, err := time.Parse("2006-01-02", entry.StartDate)
	if err != nil {
		return false
	}

	// Check if date is before start date
	if date.Before(startDate) {
		return false
	}

	// Check if date is after end date (if set)
	if entry.EndDate != nil {
		endDate, err := time.Parse("2006-01-02", *entry.EndDate)
		if err == nil && date.After(endDate) {
			return false
		}
	}

	// Check frequency
	switch entry.Frequency {
	case "once_off":
		// Only occurs on start date
		return date.Format("2006-01-02") == entry.StartDate

	case "daily":
		// Occurs every day
		return true

	case "weekly":
		// Occurs on the same day of week as start date
		if entry.DayOfWeek != nil {
			return int(date.Weekday()) == *entry.DayOfWeek
		}
		return int(date.Weekday()) == int(startDate.Weekday())

	case "fortnightly":
		// Occurs every 14 days from start date
		daysDiff := int(date.Sub(startDate).Hours() / 24)
		return daysDiff >= 0 && daysDiff%14 == 0

	case "monthly":
		// Occurs on the same day of month
		if entry.DayOfMonth != nil {
			// Use specified day of month
			return date.Day() == *entry.DayOfMonth
		}
		// Use start date's day of month
		return date.Day() == startDate.Day()

	case "annually":
		// Occurs on the same month and day each year
		return date.Month() == startDate.Month() && date.Day() == startDate.Day()

	default:
		return false
	}
}

// GetBudgetHealth calculates a health score (0-100) for a budget
func GetBudgetHealth(summary *BudgetSummary) int {
	// No entries = neutral score
	if summary.IncomeEntriesCount == 0 && summary.ExpenseEntriesCount == 0 {
		return 50
	}

	// No income = very unhealthy
	if summary.TotalMonthlyIncome == 0 {
		return 0
	}

	// Calculate surplus/deficit ratio
	ratio := summary.MonthlySurplusDeficit / summary.TotalMonthlyIncome

	// Convert to 0-100 scale
	// -100% (spending double income) = 0
	// 0% (break even) = 50
	// +50% (saving half of income) = 100
	score := 50 + (ratio * 100)

	// Cap at 0-100
	if score < 0 {
		return 0
	}
	if score > 100 {
		return 100
	}

	return int(score)
}

// BudgetHealthStatus represents the overall health of a budget
type BudgetHealthStatus struct {
	Score       int    `json:"score"`        // 0-100
	Status      string `json:"status"`       // "excellent", "good", "fair", "poor", "critical"
	Message     string `json:"message"`
	Color       string `json:"color"`        // For UI
}

// GetBudgetHealthStatus returns detailed health status
func GetBudgetHealthStatus(summary *BudgetSummary) *BudgetHealthStatus {
	score := GetBudgetHealth(summary)

	status := &BudgetHealthStatus{
		Score: score,
	}

	switch {
	case score >= 80:
		status.Status = "excellent"
		status.Message = "Your budget is in excellent shape! You're saving well."
		status.Color = "#10B981" // green

	case score >= 60:
		status.Status = "good"
		status.Message = "Your budget looks good. You have a healthy surplus."
		status.Color = "#3B82F6" // blue

	case score >= 40:
		status.Status = "fair"
		status.Message = "Your budget is balanced, but there's room for improvement."
		status.Color = "#F59E0B" // yellow

	case score >= 20:
		status.Status = "poor"
		status.Message = "Your expenses are close to or exceeding your income. Consider adjustments."
		status.Color = "#EF4444" // red

	default:
		status.Status = "critical"
		status.Message = "Your expenses significantly exceed your income. Immediate action needed."
		status.Color = "#DC2626" // dark red
	}

	return status
}
