package budget

import (
	"time"

	"github.com/google/uuid"
)

// Account represents a financial account (bank, credit card, cash, etc.)
type Account struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	AccountType string    `json:"account_type"` // 'checking', 'savings', 'credit_card', 'cash', 'investment'
	Balance     float64   `json:"balance"`
	Currency    string    `json:"currency"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Category represents an income or expense category
type Category struct {
	ID               uuid.UUID  `json:"id"`
	UserID           uuid.UUID  `json:"user_id"`
	Name             string     `json:"name"`
	CategoryType     string     `json:"category_type"` // 'income' or 'expense'
	Color            *string    `json:"color,omitempty"`
	Icon             *string    `json:"icon,omitempty"`
	ParentCategoryID *uuid.UUID `json:"parent_category_id,omitempty"`
	IsActive         bool       `json:"is_active"`
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// Budget represents a budget plan
type Budget struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description *string   `json:"description,omitempty"`
	IsActive    bool      `json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// BudgetEntry represents a planned recurring income or expense
type BudgetEntry struct {
	ID            uuid.UUID              `json:"id"`
	BudgetID      uuid.UUID              `json:"budget_id"`
	CategoryID    *uuid.UUID             `json:"category_id,omitempty"`
	Name          string                 `json:"name"`
	Description   *string                `json:"description,omitempty"`
	Amount        float64                `json:"amount"`
	EntryType     string                 `json:"entry_type"` // 'income' or 'expense'
	Frequency     string                 `json:"frequency"`  // 'once_off', 'daily', 'weekly', 'fortnightly', 'monthly', 'annually'
	DayOfMonth    *int                   `json:"day_of_month,omitempty"`
	DayOfWeek     *int                   `json:"day_of_week,omitempty"`
	StartDate     string                 `json:"start_date"` // DATE format
	EndDate       *string                `json:"end_date,omitempty"`
	MatchingRules map[string]interface{} `json:"matching_rules,omitempty"`
	IsActive      bool                   `json:"is_active"`
	CreatedAt     time.Time              `json:"created_at"`
	UpdatedAt     time.Time              `json:"updated_at"`
}

// CreateAccountRequest represents the request body for creating an account
type CreateAccountRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	AccountType string  `json:"account_type" validate:"required,oneof=checking savings credit_card cash investment"`
	Balance     float64 `json:"balance"`
	Currency    string  `json:"currency" validate:"required,len=3"`
}

// UpdateAccountRequest represents the request body for updating an account
type UpdateAccountRequest struct {
	Name        *string  `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	AccountType *string  `json:"account_type,omitempty" validate:"omitempty,oneof=checking savings credit_card cash investment"`
	Balance     *float64 `json:"balance,omitempty"`
	Currency    *string  `json:"currency,omitempty" validate:"omitempty,len=3"`
	IsActive    *bool    `json:"is_active,omitempty"`
}

// CreateCategoryRequest represents the request body for creating a category
type CreateCategoryRequest struct {
	Name             string     `json:"name" validate:"required,min=1,max=255"`
	CategoryType     string     `json:"category_type" validate:"required,oneof=income expense"`
	Color            *string    `json:"color,omitempty" validate:"omitempty,len=7"`
	Icon             *string    `json:"icon,omitempty" validate:"omitempty,max=50"`
	ParentCategoryID *uuid.UUID `json:"parent_category_id,omitempty"`
}

// UpdateCategoryRequest represents the request body for updating a category
type UpdateCategoryRequest struct {
	Name             *string    `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	CategoryType     *string    `json:"category_type,omitempty" validate:"omitempty,oneof=income expense"`
	Color            *string    `json:"color,omitempty" validate:"omitempty,len=7"`
	Icon             *string    `json:"icon,omitempty" validate:"omitempty,max=50"`
	ParentCategoryID *uuid.UUID `json:"parent_category_id,omitempty"`
	IsActive         *bool      `json:"is_active,omitempty"`
}

// CreateBudgetRequest represents the request body for creating a budget
type CreateBudgetRequest struct {
	Name        string  `json:"name" validate:"required,min=1,max=255"`
	Description *string `json:"description,omitempty"`
}

// UpdateBudgetRequest represents the request body for updating a budget
type UpdateBudgetRequest struct {
	Name        *string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description,omitempty"`
	IsActive    *bool   `json:"is_active,omitempty"`
}

// CreateBudgetEntryRequest represents the request body for creating a budget entry
type CreateBudgetEntryRequest struct {
	CategoryID    *uuid.UUID             `json:"category_id,omitempty"`
	Name          string                 `json:"name" validate:"required,min=1,max=255"`
	Description   *string                `json:"description,omitempty"`
	Amount        float64                `json:"amount" validate:"required,gt=0"`
	EntryType     string                 `json:"entry_type" validate:"required,oneof=income expense"`
	Frequency     string                 `json:"frequency" validate:"required,oneof=once_off daily weekly fortnightly monthly annually"`
	DayOfMonth    *int                   `json:"day_of_month,omitempty" validate:"omitempty,gte=1,lte=31"`
	DayOfWeek     *int                   `json:"day_of_week,omitempty" validate:"omitempty,gte=0,lte=6"`
	StartDate     string                 `json:"start_date" validate:"required"`
	EndDate       *string                `json:"end_date,omitempty"`
	MatchingRules map[string]interface{} `json:"matching_rules,omitempty"`
}

// UpdateBudgetEntryRequest represents the request body for updating a budget entry
type UpdateBudgetEntryRequest struct {
	CategoryID    *uuid.UUID             `json:"category_id,omitempty"`
	Name          *string                `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description   *string                `json:"description,omitempty"`
	Amount        *float64               `json:"amount,omitempty" validate:"omitempty,gt=0"`
	EntryType     *string                `json:"entry_type,omitempty" validate:"omitempty,oneof=income expense"`
	Frequency     *string                `json:"frequency,omitempty" validate:"omitempty,oneof=once_off daily weekly fortnightly monthly annually"`
	DayOfMonth    *int                   `json:"day_of_month,omitempty" validate:"omitempty,gte=1,lte=31"`
	DayOfWeek     *int                   `json:"day_of_week,omitempty" validate:"omitempty,gte=0,lte=6"`
	StartDate     *string                `json:"start_date,omitempty"`
	EndDate       *string                `json:"end_date,omitempty"`
	MatchingRules map[string]interface{} `json:"matching_rules,omitempty"`
	IsActive      *bool                  `json:"is_active,omitempty"`
}

// BudgetWithEntries represents a budget with all its entries
type BudgetWithEntries struct {
	Budget
	Entries []BudgetEntry `json:"entries"`
}

// BudgetSummary represents summary statistics for a budget
type BudgetSummary struct {
	BudgetID                uuid.UUID `json:"budget_id"`
	TotalMonthlyIncome      float64   `json:"total_monthly_income"`
	TotalMonthlyExpenses    float64   `json:"total_monthly_expenses"`
	MonthlySurplusDeficit   float64   `json:"monthly_surplus_deficit"`
	TotalAnnualIncome       float64   `json:"total_annual_income"`
	TotalAnnualExpenses     float64   `json:"total_annual_expenses"`
	AnnualSurplusDeficit    float64   `json:"annual_surplus_deficit"`
	IncomeEntriesCount      int       `json:"income_entries_count"`
	ExpenseEntriesCount     int       `json:"expense_entries_count"`
}

// Transaction represents a financial transaction
type Transaction struct {
	ID              uuid.UUID  `json:"id"`
	UserID          uuid.UUID  `json:"user_id"`
	AccountID       uuid.UUID  `json:"account_id"`
	CategoryID      *uuid.UUID `json:"category_id,omitempty"`
	BudgetEntryID   *uuid.UUID `json:"budget_entry_id,omitempty"`
	Amount          float64    `json:"amount"`
	TransactionType string     `json:"transaction_type"` // 'income' or 'expense'
	Description     *string    `json:"description,omitempty"`
	TransactionDate string     `json:"transaction_date"` // DATE format
	Notes           *string    `json:"notes,omitempty"`
	MatchConfidence string     `json:"match_confidence"` // 'manual', 'auto_high', 'auto_low', 'unmatched'
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// CreateTransactionRequest represents the request body for creating a transaction
type CreateTransactionRequest struct {
	AccountID       uuid.UUID  `json:"account_id" validate:"required"`
	CategoryID      *uuid.UUID `json:"category_id,omitempty"`
	Amount          float64    `json:"amount" validate:"required,gt=0"`
	TransactionType string     `json:"transaction_type" validate:"required,oneof=income expense"`
	Description     *string    `json:"description,omitempty"`
	TransactionDate string     `json:"transaction_date" validate:"required"`
	Notes           *string    `json:"notes,omitempty"`
}

// UpdateTransactionRequest represents the request body for updating a transaction
type UpdateTransactionRequest struct {
	AccountID       *uuid.UUID `json:"account_id,omitempty"`
	CategoryID      *uuid.UUID `json:"category_id,omitempty"`
	BudgetEntryID   *uuid.UUID `json:"budget_entry_id,omitempty"`
	Amount          *float64   `json:"amount,omitempty" validate:"omitempty,gt=0"`
	TransactionType *string    `json:"transaction_type,omitempty" validate:"omitempty,oneof=income expense"`
	Description     *string    `json:"description,omitempty"`
	TransactionDate *string    `json:"transaction_date,omitempty"`
	Notes           *string    `json:"notes,omitempty"`
	MatchConfidence *string    `json:"match_confidence,omitempty" validate:"omitempty,oneof=manual auto_high auto_low unmatched"`
}
