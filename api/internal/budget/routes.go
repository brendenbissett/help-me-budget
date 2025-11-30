package budget

import (
	"github.com/gofiber/fiber/v2"
)

// SetupBudgetRoutes configures all budget-related API endpoints
func SetupBudgetRoutes(app *fiber.App) {
	// Account management routes
	accounts := app.Group("/api/accounts")
	accounts.Get("/", GetAccountsHandler)                     // List all accounts
	accounts.Get("/:id", GetAccountHandler)                   // Get specific account
	accounts.Post("/", CreateAccountHandler)                  // Create new account
	accounts.Put("/:id", UpdateAccountHandler)                // Update account
	accounts.Delete("/:id", DeleteAccountHandler)             // Delete account (soft delete)
	accounts.Get("/balance/total", GetTotalBalanceHandler)    // Get total balance across all accounts

	// Category management routes
	categories := app.Group("/api/categories")
	categories.Get("/", GetCategoriesHandler)                 // List all categories (supports ?type=income|expense filter)
	categories.Get("/:id", GetCategoryHandler)                // Get specific category
	categories.Post("/", CreateCategoryHandler)               // Create new category
	categories.Put("/:id", UpdateCategoryHandler)             // Update category
	categories.Delete("/:id", DeleteCategoryHandler)          // Delete category (soft delete)
	categories.Post("/seed", SeedDefaultCategoriesHandler)    // Seed default categories for new user

	// Budget management routes
	budgets := app.Group("/api/budgets")
	budgets.Get("/", GetBudgetsHandler)                       // List all budgets
	budgets.Get("/:id", GetBudgetHandler)                     // Get specific budget
	budgets.Get("/:id/full", GetBudgetWithEntriesHandler)     // Get budget with all entries
	budgets.Post("/", CreateBudgetHandler)                    // Create new budget
	budgets.Put("/:id", UpdateBudgetHandler)                  // Update budget
	budgets.Delete("/:id", DeleteBudgetHandler)               // Delete budget (soft delete)
	budgets.Get("/:id/summary", GetBudgetSummaryHandler)      // Get budget summary (income/expense totals, health)
	budgets.Get("/:id/projection", ProjectCashFlowHandler)    // Project cash flow (supports ?days=90&starting_balance=1000)

	// Budget entry routes (nested under budgets)
	budgets.Get("/:id/entries", GetBudgetEntriesHandler)      // List all entries for a budget
	budgets.Post("/:id/entries", CreateBudgetEntryHandler)    // Create new budget entry
	budgets.Put("/:id/entries/:entryId", UpdateBudgetEntryHandler)   // Update budget entry
	budgets.Delete("/:id/entries/:entryId", DeleteBudgetEntryHandler) // Delete budget entry

	// Transaction management routes
	transactions := app.Group("/api/transactions")
	transactions.Get("/", GetTransactionsHandler)                      // List all transactions (supports filters: ?account_id=&category_id=&start_date=&end_date=)
	transactions.Get("/unmatched", GetUnmatchedTransactionsHandler)    // Get unmatched transactions
	transactions.Get("/:id", GetTransactionHandler)                    // Get specific transaction
	transactions.Post("/", CreateTransactionHandler)                   // Create new transaction
	transactions.Put("/:id", UpdateTransactionHandler)                 // Update transaction
	transactions.Delete("/:id", DeleteTransactionHandler)              // Delete transaction
	transactions.Post("/:id/categorize", CategorizeTransactionHandler) // Assign category to transaction
	transactions.Post("/:id/link", LinkTransactionHandler)             // Link transaction to budget entry

	// Dashboard routes
	dashboard := app.Group("/api/dashboard")
	dashboard.Get("/summary", GetDashboardSummaryHandler)                // Get comprehensive dashboard overview
	dashboard.Get("/recent-activity", GetRecentActivityHandler)          // Get recent transactions (supports ?limit=20)
	dashboard.Get("/spending-by-category", GetSpendingByCategoryHandler) // Get spending breakdown (supports ?start_date=&end_date=)

	// Matching/automation routes
	matching := app.Group("/api/matching")
	matching.Get("/suggestions/:id", GetSuggestedMatchesHandler)           // Get match suggestions for a transaction
	matching.Post("/auto-match/:id", AutoMatchTransactionHandler)          // Auto-match a single transaction
	matching.Post("/bulk-auto-match", BulkAutoMatchHandler)                // Auto-match all unmatched transactions
	matching.Post("/teach/:id", TeachMatchHandler)                         // Link transaction + create matching rules

	// Budget entry matching rules (nested under budgets)
	budgets.Post("/:id/entries/:entryId/matching-rules", UpdateBudgetEntryMatchingRulesHandler) // Update matching rules for budget entry

	// Reports and analytics routes
	reports := app.Group("/api/reports")
	reports.Get("/spending-trends", GetSpendingTrendsHandler)           // Get spending trends by category over time (supports ?start_date=&end_date=)
	reports.Get("/budget-variance", GetBudgetVarianceHandler)           // Get budget vs actual comparison (supports ?month=YYYY-MM)
	reports.Get("/cash-flow-projection", GetCashFlowProjectionHandler) // Get projected cash flow (supports ?days=90&starting_balance=1000)
	reports.Get("/top-expenses", GetTopExpensesHandler)                 // Get top spending categories (supports ?start_date=&end_date=&limit=10)
}
