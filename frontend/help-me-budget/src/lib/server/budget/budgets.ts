import { authenticatedFetchWithUser } from '../api-client';

export interface Budget {
	id: string;
	user_id: string;
	name: string;
	description?: string | null;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface BudgetEntry {
	id: string;
	budget_id: string;
	category_id?: string | null;
	name: string;
	description?: string | null;
	amount: number;
	entry_type: 'income' | 'expense';
	frequency: 'once_off' | 'daily' | 'weekly' | 'fortnightly' | 'monthly' | 'annually';
	day_of_month?: number | null;
	day_of_week?: number | null;
	start_date: string;
	end_date?: string | null;
	matching_rules?: Record<string, any> | null;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface BudgetWithEntries extends Budget {
	entries: BudgetEntry[];
}

export interface BudgetSummary {
	budget_id: string;
	total_monthly_income: number;
	total_monthly_expenses: number;
	monthly_surplus_deficit: number;
	total_annual_income: number;
	total_annual_expenses: number;
	annual_surplus_deficit: number;
	income_entries_count: number;
	expense_entries_count: number;
}

export interface BudgetHealthStatus {
	score: number; // 0-100
	status: 'excellent' | 'good' | 'fair' | 'poor' | 'critical';
	message: string;
	color: string;
}

export interface DailyProjection {
	date: string;
	balance: number;
	daily_income: number;
	daily_expenses: number;
	daily_net: number;
}

export interface MonthlyBreakdown {
	month: string;
	income: number;
	expenses: number;
	net: number;
	ending_balance: number;
}

export interface CashFlowProjection {
	start_date: string;
	end_date: string;
	starting_balance: number;
	ending_balance: number;
	total_income: number;
	total_expenses: number;
	net_cash_flow: number;
	daily_projections: DailyProjection[];
	monthly_breakdown: MonthlyBreakdown[];
}

export interface CreateBudgetRequest {
	name: string;
	description?: string;
}

export interface UpdateBudgetRequest {
	name?: string;
	description?: string;
	is_active?: boolean;
}

export interface CreateBudgetEntryRequest {
	category_id?: string;
	name: string;
	description?: string;
	amount: number;
	entry_type: 'income' | 'expense';
	frequency: 'once_off' | 'daily' | 'weekly' | 'fortnightly' | 'monthly' | 'annually';
	day_of_month?: number;
	day_of_week?: number;
	start_date: string;
	end_date?: string;
	matching_rules?: Record<string, any>;
}

export interface UpdateBudgetEntryRequest {
	category_id?: string;
	name?: string;
	description?: string;
	amount?: number;
	entry_type?: 'income' | 'expense';
	frequency?: 'once_off' | 'daily' | 'weekly' | 'fortnightly' | 'monthly' | 'annually';
	day_of_month?: number;
	day_of_week?: number;
	start_date?: string;
	end_date?: string;
	matching_rules?: Record<string, any>;
	is_active?: boolean;
}

/**
 * Get all budgets for a user
 */
export async function getBudgets(userId: string): Promise<Budget[]> {
	const response = await authenticatedFetchWithUser('/api/budgets', userId);

	if (!response.ok) {
		throw new Error(`Failed to fetch budgets: ${response.statusText}`);
	}

	const data = await response.json();
	return data.budgets || [];
}

/**
 * Get a specific budget by ID
 */
export async function getBudget(userId: string, budgetId: string): Promise<Budget> {
	const response = await authenticatedFetchWithUser(`/api/budgets/${budgetId}`, userId);

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Budget not found');
		}
		throw new Error(`Failed to fetch budget: ${response.statusText}`);
	}

	return await response.json();
}

/**
 * Get a budget with all its entries
 */
export async function getBudgetWithEntries(
	userId: string,
	budgetId: string
): Promise<BudgetWithEntries> {
	const response = await authenticatedFetchWithUser(`/api/budgets/${budgetId}/full`, userId);

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Budget not found');
		}
		throw new Error(`Failed to fetch budget: ${response.statusText}`);
	}

	return await response.json();
}

/**
 * Create a new budget
 */
export async function createBudget(
	userId: string,
	budget: CreateBudgetRequest
): Promise<Budget> {
	const response = await authenticatedFetchWithUser('/api/budgets', userId, {
		method: 'POST',
		body: JSON.stringify(budget)
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to create budget');
	}

	return await response.json();
}

/**
 * Update an existing budget
 */
export async function updateBudget(
	userId: string,
	budgetId: string,
	updates: UpdateBudgetRequest
): Promise<Budget> {
	const response = await authenticatedFetchWithUser(`/api/budgets/${budgetId}`, userId, {
		method: 'PUT',
		body: JSON.stringify(updates)
	});

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Budget not found');
		}
		const error = await response.json();
		throw new Error(error.error || 'Failed to update budget');
	}

	return await response.json();
}

/**
 * Delete a budget (soft delete)
 */
export async function deleteBudget(userId: string, budgetId: string): Promise<void> {
	const response = await authenticatedFetchWithUser(`/api/budgets/${budgetId}`, userId, {
		method: 'DELETE'
	});

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Budget not found');
		}
		const error = await response.json();
		throw new Error(error.error || 'Failed to delete budget');
	}
}

/**
 * Get budget summary (income/expense totals, health status)
 */
export async function getBudgetSummary(
	userId: string,
	budgetId: string
): Promise<{ summary: BudgetSummary; health: BudgetHealthStatus }> {
	const response = await authenticatedFetchWithUser(`/api/budgets/${budgetId}/summary`, userId);

	if (!response.ok) {
		throw new Error(`Failed to fetch budget summary: ${response.statusText}`);
	}

	return await response.json();
}

/**
 * Project cash flow for a budget
 */
export async function projectCashFlow(
	userId: string,
	budgetId: string,
	startingBalance: number = 0,
	days: number = 90
): Promise<CashFlowProjection> {
	const response = await authenticatedFetchWithUser(
		`/api/budgets/${budgetId}/projection?starting_balance=${startingBalance}&days=${days}`,
		userId
	);

	if (!response.ok) {
		throw new Error(`Failed to project cash flow: ${response.statusText}`);
	}

	return await response.json();
}

/**
 * Get all entries for a budget
 */
export async function getBudgetEntries(userId: string, budgetId: string): Promise<BudgetEntry[]> {
	const response = await authenticatedFetchWithUser(`/api/budgets/${budgetId}/entries`, userId);

	if (!response.ok) {
		throw new Error(`Failed to fetch budget entries: ${response.statusText}`);
	}

	const data = await response.json();
	return data.entries || [];
}

/**
 * Create a new budget entry
 */
export async function createBudgetEntry(
	userId: string,
	budgetId: string,
	entry: CreateBudgetEntryRequest
): Promise<BudgetEntry> {
	const response = await authenticatedFetchWithUser(`/api/budgets/${budgetId}/entries`, userId, {
		method: 'POST',
		body: JSON.stringify(entry)
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to create budget entry');
	}

	return await response.json();
}

/**
 * Update an existing budget entry
 */
export async function updateBudgetEntry(
	userId: string,
	budgetId: string,
	entryId: string,
	updates: UpdateBudgetEntryRequest
): Promise<BudgetEntry> {
	const response = await authenticatedFetchWithUser(
		`/api/budgets/${budgetId}/entries/${entryId}`,
		userId,
		{
			method: 'PUT',
			body: JSON.stringify(updates)
		}
	);

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Budget entry not found');
		}
		const error = await response.json();
		throw new Error(error.error || 'Failed to update budget entry');
	}

	return await response.json();
}

/**
 * Delete a budget entry (soft delete)
 */
export async function deleteBudgetEntry(
	userId: string,
	budgetId: string,
	entryId: string
): Promise<void> {
	const response = await authenticatedFetchWithUser(
		`/api/budgets/${budgetId}/entries/${entryId}`,
		userId,
		{
			method: 'DELETE'
		}
	);

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Budget entry not found');
		}
		const error = await response.json();
		throw new Error(error.error || 'Failed to delete budget entry');
	}
}
