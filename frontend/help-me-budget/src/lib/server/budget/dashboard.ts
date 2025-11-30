import { authenticatedFetchWithUser } from '../api-client';

// Dashboard types matching backend models
export interface DashboardSummary {
	total_balance: number;
	account_count: number;
	month_to_date_income: number;
	month_to_date_expenses: number;
	month_to_date_net: number;
	budgeted_monthly_income: number;
	budgeted_monthly_expense: number;
	budget_health_score: number;
	budget_health_status: string;
	budget_health_message: string;
	budget_health_color: string;
	upcoming_bills: UpcomingBill[];
	recent_transactions: RecentTransaction[];
	spending_by_category: CategorySpending[];
}

export interface UpcomingBill {
	id: string;
	name: string;
	amount: number;
	due_date: string;
	category_id?: string | null;
	is_overdue: boolean;
}

export interface RecentTransaction {
	id: string;
	user_id: string;
	account_id: string;
	category_id?: string | null;
	budget_entry_id?: string | null;
	amount: number;
	transaction_type: 'income' | 'expense';
	description?: string | null;
	transaction_date: string;
	notes?: string | null;
	match_confidence: 'manual' | 'auto_high' | 'auto_low' | 'unmatched';
	created_at: string;
	updated_at: string;
}

export interface CategorySpending {
	category_id?: string | null;
	category_name: string;
	total_amount: number;
	percentage: number;
	color?: string | null;
}

export interface SpendingByCategoryResponse {
	spending_by_category: CategorySpending[];
	total_expenses: number;
	start_date: string;
	end_date: string;
}

/**
 * Get comprehensive dashboard summary with all metrics
 */
export async function getDashboardSummary(userId: string): Promise<DashboardSummary> {
	const response = await authenticatedFetchWithUser('/api/dashboard/summary', userId, {
		method: 'GET'
	});

	if (!response.ok) {
		throw new Error(`Failed to fetch dashboard summary: ${response.statusText}`);
	}

	return await response.json();
}

/**
 * Get recent transactions with optional limit
 */
export async function getRecentActivity(
	userId: string,
	limit: number = 20
): Promise<RecentTransaction[]> {
	const url = `/api/dashboard/recent-activity?limit=${limit}`;

	const response = await authenticatedFetchWithUser(url, userId, {
		method: 'GET'
	});

	if (!response.ok) {
		throw new Error(`Failed to fetch recent activity: ${response.statusText}`);
	}

	const data = await response.json();
	return data.transactions || [];
}

/**
 * Get spending breakdown by category for a date range
 */
export async function getSpendingByCategory(
	userId: string,
	startDate?: string,
	endDate?: string
): Promise<SpendingByCategoryResponse> {
	const params = new URLSearchParams();
	if (startDate) params.append('start_date', startDate);
	if (endDate) params.append('end_date', endDate);

	const url = `/api/dashboard/spending-by-category${params.toString() ? `?${params.toString()}` : ''}`;

	const response = await authenticatedFetchWithUser(url, userId, {
		method: 'GET'
	});

	if (!response.ok) {
		throw new Error(`Failed to fetch spending by category: ${response.statusText}`);
	}

	return await response.json();
}
