import { authenticatedFetchWithUser } from '../api-client';

// ============================================================================
// Types
// ============================================================================

export interface SpendingTrend {
	month: string; // YYYY-MM format
	category_id: string;
	category: string;
	amount: number;
}

export interface BudgetVariance {
	entry_id: string;
	entry_name: string;
	category: string;
	budgeted: number;
	actual: number;
	variance: number; // positive = under budget, negative = over budget
	variance_pct: number;
}

export interface DailyCashFlowProjection {
	date: string; // YYYY-MM-DD
	projected_income: number;
	projected_expenses: number;
	projected_balance: number;
}

export interface TopExpense {
	category_id: string;
	category_name: string;
	total_amount: number;
	percentage: number; // Percentage of total expenses
	count: number; // Number of transactions
}

// ============================================================================
// API Functions
// ============================================================================

/**
 * Get spending trends by category over time
 * @param userId - User ID
 * @param startDate - Start date (YYYY-MM-DD), defaults to 6 months ago
 * @param endDate - End date (YYYY-MM-DD), defaults to today
 */
export async function getSpendingTrends(
	userId: string,
	startDate?: string,
	endDate?: string
): Promise<SpendingTrend[]> {
	const params = new URLSearchParams();
	if (startDate) params.append('start_date', startDate);
	if (endDate) params.append('end_date', endDate);

	const queryString = params.toString();
	const url = `/api/reports/spending-trends${queryString ? '?' + queryString : ''}`;

	const response = await authenticatedFetchWithUser(url, userId, {
		method: 'GET'
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to get spending trends');
	}

	return response.json();
}

/**
 * Get budget vs actual comparison for a month
 * @param userId - User ID
 * @param month - Month in YYYY-MM format, defaults to current month
 */
export async function getBudgetVariance(
	userId: string,
	month?: string
): Promise<BudgetVariance[]> {
	const params = new URLSearchParams();
	if (month) params.append('month', month);

	const queryString = params.toString();
	const url = `/api/reports/budget-variance${queryString ? '?' + queryString : ''}`;

	const response = await authenticatedFetchWithUser(url, userId, {
		method: 'GET'
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to get budget variance');
	}

	return response.json();
}

/**
 * Get projected cash flow based on budget entries
 * @param userId - User ID
 * @param days - Number of days to project (default 90)
 * @param startingBalance - Starting balance (default 0)
 */
export async function getCashFlowProjection(
	userId: string,
	days: number = 90,
	startingBalance: number = 0
): Promise<DailyCashFlowProjection[]> {
	const params = new URLSearchParams();
	params.append('days', days.toString());
	params.append('starting_balance', startingBalance.toString());

	const url = `/api/reports/cash-flow-projection?${params.toString()}`;

	const response = await authenticatedFetchWithUser(url, userId, {
		method: 'GET'
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to get cash flow projection');
	}

	return response.json();
}

/**
 * Get top spending categories
 * @param userId - User ID
 * @param startDate - Start date (YYYY-MM-DD), defaults to start of current month
 * @param endDate - End date (YYYY-MM-DD), defaults to today
 * @param limit - Number of top expenses to return (default 10, max 50)
 */
export async function getTopExpenses(
	userId: string,
	startDate?: string,
	endDate?: string,
	limit: number = 10
): Promise<TopExpense[]> {
	const params = new URLSearchParams();
	if (startDate) params.append('start_date', startDate);
	if (endDate) params.append('end_date', endDate);
	params.append('limit', limit.toString());

	const url = `/api/reports/top-expenses?${params.toString()}`;

	const response = await authenticatedFetchWithUser(url, userId, {
		method: 'GET'
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to get top expenses');
	}

	return response.json();
}
