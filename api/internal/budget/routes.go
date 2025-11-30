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
}
