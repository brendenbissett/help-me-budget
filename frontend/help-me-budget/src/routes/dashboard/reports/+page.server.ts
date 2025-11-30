import type { PageServerLoad } from './$types';
import { getLocalUserId } from '$lib/server/auth-helpers';
import {
	getSpendingTrends,
	getBudgetVariance,
	getCashFlowProjection,
	getTopExpenses
} from '$lib/server/budget/reports';
import { getAccounts } from '$lib/server/budget/accounts';

export const load: PageServerLoad = async ({ locals, url }) => {
	try {
		const userId = await getLocalUserId(locals.supabase);

		// Get query parameters
		const month = url.searchParams.get('month') || undefined; // YYYY-MM
		const days = parseInt(url.searchParams.get('days') || '90');

		// Get current month for defaults
		const now = new Date();
		const currentMonth = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}`;
		const startOfMonth = `${now.getFullYear()}-${String(now.getMonth() + 1).padStart(2, '0')}-01`;
		const endOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0)
			.toISOString()
			.split('T')[0];

		// Get 6 months ago for trends
		const sixMonthsAgo = new Date();
		sixMonthsAgo.setMonth(sixMonthsAgo.getMonth() - 6);
		const trendsStartDate = sixMonthsAgo.toISOString().split('T')[0];

		// Fetch all reports data in parallel
		const [spendingTrends, budgetVariance, cashFlowProjection, topExpenses, accounts] =
			await Promise.all([
				getSpendingTrends(userId, trendsStartDate, endOfMonth),
				getBudgetVariance(userId, month || currentMonth),
				getCashFlowProjection(userId, days, 0), // Starting balance of 0
				getTopExpenses(userId, startOfMonth, endOfMonth, 10),
				getAccounts(userId)
			]);

		// Calculate total balance from accounts for cash flow projection starting point
		const totalBalance = accounts.reduce((sum, account) => sum + account.balance, 0);

		// Re-fetch cash flow with actual starting balance
		const cashFlowWithBalance = await getCashFlowProjection(userId, days, totalBalance);

		return {
			spendingTrends,
			budgetVariance,
			cashFlowProjection: cashFlowWithBalance,
			topExpenses,
			totalBalance,
			currentMonth
		};
	} catch (error) {
		console.error('Error loading reports:', error);
		return {
			spendingTrends: [],
			budgetVariance: [],
			cashFlowProjection: [],
			topExpenses: [],
			totalBalance: 0,
			currentMonth: new Date().toISOString().slice(0, 7),
			loadError: error instanceof Error ? error.message : 'Failed to load reports'
		};
	}
};
