-- Create budgets table (users can have multiple budgets for different scenarios)
CREATE TABLE budget.budgets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create categories table (for organizing budget entries and transactions)
CREATE TABLE budget.categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    category_type VARCHAR(10) NOT NULL, -- 'income' or 'expense'
    color VARCHAR(7), -- Hex color for UI (e.g., '#FF5733')
    icon VARCHAR(50), -- Icon identifier for UI
    parent_category_id UUID REFERENCES budget.categories(id) ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_category_type CHECK (category_type IN ('income', 'expense'))
);

-- Create budget_entries table (planned recurring income/expenses)
CREATE TABLE budget.budget_entries (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    budget_id UUID NOT NULL REFERENCES budget.budgets(id) ON DELETE CASCADE,
    category_id UUID REFERENCES budget.categories(id) ON DELETE SET NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    amount NUMERIC(20, 2) NOT NULL,
    entry_type VARCHAR(10) NOT NULL, -- 'income' or 'expense'
    frequency VARCHAR(20) NOT NULL, -- 'once_off', 'daily', 'weekly', 'fortnightly', 'monthly', 'annually'
    day_of_month INT, -- 1-31 for monthly/annually, null for other frequencies
    day_of_week INT, -- 0-6 (Sunday-Saturday) for weekly/fortnightly, null otherwise
    start_date DATE NOT NULL,
    end_date DATE, -- Optional: when this recurring entry stops
    matching_rules JSONB, -- Flexible matching criteria: {"description_contains": ["NETFLIX"], "amount_tolerance": 1.00, "merchant_name": "Netflix Inc"}
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_entry_type CHECK (entry_type IN ('income', 'expense')),
    CONSTRAINT valid_frequency CHECK (frequency IN ('once_off', 'daily', 'weekly', 'fortnightly', 'monthly', 'annually')),
    CONSTRAINT valid_day_of_month CHECK (day_of_month IS NULL OR (day_of_month >= 1 AND day_of_month <= 31)),
    CONSTRAINT valid_day_of_week CHECK (day_of_week IS NULL OR (day_of_week >= 0 AND day_of_week <= 6))
);

-- Create accounts table (bank accounts, credit cards, cash, etc.)
CREATE TABLE budget.accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    account_type VARCHAR(50) NOT NULL, -- 'checking', 'savings', 'credit_card', 'cash', 'investment'
    balance NUMERIC(20, 2) DEFAULT 0.00,
    currency VARCHAR(3) DEFAULT 'USD',
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create transactions table (actual transactions that occurred)
CREATE TABLE budget.transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    account_id UUID NOT NULL REFERENCES budget.accounts(id) ON DELETE CASCADE,
    category_id UUID REFERENCES budget.categories(id) ON DELETE SET NULL,
    budget_entry_id UUID REFERENCES budget.budget_entries(id) ON DELETE SET NULL, -- Link to planned entry (only when match is confirmed)
    amount NUMERIC(20, 2) NOT NULL,
    transaction_type VARCHAR(10) NOT NULL, -- 'income' or 'expense'
    description TEXT,
    transaction_date DATE NOT NULL,
    notes TEXT,
    match_confidence VARCHAR(20) DEFAULT 'unmatched', -- 'manual', 'auto_high', 'auto_low', 'unmatched'
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT valid_transaction_type CHECK (transaction_type IN ('income', 'expense')),
    CONSTRAINT valid_match_confidence CHECK (match_confidence IN ('manual', 'auto_high', 'auto_low', 'unmatched'))
);

-- Indexes for performance
CREATE INDEX idx_budgets_user_id ON budget.budgets(user_id);
CREATE INDEX idx_categories_user_id ON budget.categories(user_id);
CREATE INDEX idx_categories_parent ON budget.categories(parent_category_id);
CREATE INDEX idx_budget_entries_budget_id ON budget.budget_entries(budget_id);
CREATE INDEX idx_budget_entries_category_id ON budget.budget_entries(category_id);
CREATE INDEX idx_budget_entries_dates ON budget.budget_entries(start_date, end_date);
CREATE INDEX idx_budget_entries_matching_rules ON budget.budget_entries USING gin(matching_rules); -- GIN index for JSONB queries
CREATE INDEX idx_accounts_user_id ON budget.accounts(user_id);
CREATE INDEX idx_transactions_user_id ON budget.transactions(user_id);
CREATE INDEX idx_transactions_account_id ON budget.transactions(account_id);
CREATE INDEX idx_transactions_category_id ON budget.transactions(category_id);
CREATE INDEX idx_transactions_budget_entry_id ON budget.transactions(budget_entry_id);
CREATE INDEX idx_transactions_date ON budget.transactions(transaction_date);
CREATE INDEX idx_transactions_match_confidence ON budget.transactions(match_confidence);

-- Create updated_at triggers
CREATE TRIGGER update_budgets_updated_at
    BEFORE UPDATE ON budget.budgets
    FOR EACH ROW
    EXECUTE FUNCTION auth.update_updated_at_column();

CREATE TRIGGER update_categories_updated_at
    BEFORE UPDATE ON budget.categories
    FOR EACH ROW
    EXECUTE FUNCTION auth.update_updated_at_column();

CREATE TRIGGER update_budget_entries_updated_at
    BEFORE UPDATE ON budget.budget_entries
    FOR EACH ROW
    EXECUTE FUNCTION auth.update_updated_at_column();

CREATE TRIGGER update_accounts_updated_at
    BEFORE UPDATE ON budget.accounts
    FOR EACH ROW
    EXECUTE FUNCTION auth.update_updated_at_column();

CREATE TRIGGER update_transactions_updated_at
    BEFORE UPDATE ON budget.transactions
    FOR EACH ROW
    EXECUTE FUNCTION auth.update_updated_at_column();
