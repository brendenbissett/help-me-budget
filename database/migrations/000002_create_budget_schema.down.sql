-- Drop triggers
DROP TRIGGER IF EXISTS update_transactions_updated_at ON budget.transactions;
DROP TRIGGER IF EXISTS update_accounts_updated_at ON budget.accounts;
DROP TRIGGER IF EXISTS update_budget_entries_updated_at ON budget.budget_entries;
DROP TRIGGER IF EXISTS update_categories_updated_at ON budget.categories;
DROP TRIGGER IF EXISTS update_budgets_updated_at ON budget.budgets;

-- Drop indexes
DROP INDEX IF EXISTS budget.idx_transactions_match_confidence;
DROP INDEX IF EXISTS budget.idx_transactions_date;
DROP INDEX IF EXISTS budget.idx_transactions_budget_entry_id;
DROP INDEX IF EXISTS budget.idx_transactions_category_id;
DROP INDEX IF EXISTS budget.idx_transactions_account_id;
DROP INDEX IF EXISTS budget.idx_transactions_user_id;
DROP INDEX IF EXISTS budget.idx_accounts_user_id;
DROP INDEX IF EXISTS budget.idx_budget_entries_matching_rules;
DROP INDEX IF EXISTS budget.idx_budget_entries_dates;
DROP INDEX IF EXISTS budget.idx_budget_entries_category_id;
DROP INDEX IF EXISTS budget.idx_budget_entries_budget_id;
DROP INDEX IF EXISTS budget.idx_categories_parent;
DROP INDEX IF EXISTS budget.idx_categories_user_id;
DROP INDEX IF EXISTS budget.idx_budgets_user_id;

-- Drop tables (in reverse order of dependencies)
DROP TABLE IF EXISTS budget.transactions;
DROP TABLE IF EXISTS budget.accounts;
DROP TABLE IF EXISTS budget.budget_entries;
DROP TABLE IF EXISTS budget.categories;
DROP TABLE IF EXISTS budget.budgets;
