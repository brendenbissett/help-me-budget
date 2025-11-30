import type { PageServerLoad, Actions } from './$types';
import { fail } from '@sveltejs/kit';
import { getLocalUserId } from '$lib/server/auth-helpers';
import { getUnmatchedTransactions } from '$lib/server/budget/transactions';
import { getBudgets, getBudgetWithEntries } from '$lib/server/budget/budgets';
import { getAccounts } from '$lib/server/budget/accounts';
import { getCategories } from '$lib/server/budget/categories';
import {
	getMatchSuggestions,
	teachMatch,
	bulkAutoMatch,
	type TeachMatchRequest
} from '$lib/server/budget/matching';

export const load: PageServerLoad = async ({ locals }) => {
	try {
		const userId = await getLocalUserId(locals.supabase);

		// Fetch unmatched transactions
		const unmatchedTransactions = await getUnmatchedTransactions(userId);

		// Fetch active budget with entries
		let budget = null;
		try {
			const budgets = await getBudgets(userId);
			const activeBudget = budgets.find((b) => b.is_active);
			if (activeBudget) {
				budget = await getBudgetWithEntries(userId, activeBudget.id);
			}
		} catch (err) {
			console.error('Error fetching budget:', err);
		}

		// Fetch accounts and categories for display
		const accounts = await getAccounts(userId);
		const categories = await getCategories(userId);

		return {
			unmatchedTransactions,
			budget,
			accounts,
			categories
		};
	} catch (error) {
		console.error('Error loading matching review page:', error);
		return {
			unmatchedTransactions: [],
			budget: null,
			accounts: [],
			categories: []
		};
	}
};

export const actions: Actions = {
	// Get suggestions for a specific transaction
	getSuggestions: async ({ locals, request }) => {
		try {
			const userId = await getLocalUserId(locals.supabase);
			const formData = await request.formData();
			const transactionId = formData.get('transaction_id') as string;

			if (!transactionId) {
				return fail(400, { error: 'Transaction ID is required' });
			}

			const suggestions = await getMatchSuggestions(userId, transactionId);

			return {
				success: true,
				suggestions
			};
		} catch (error) {
			console.error('Error getting match suggestions:', error);
			return fail(500, {
				error: error instanceof Error ? error.message : 'Failed to get suggestions'
			});
		}
	},

	// Teach mode: Link transaction and create rules
	teach: async ({ locals, request }) => {
		try {
			const userId = await getLocalUserId(locals.supabase);
			const formData = await request.formData();
			const transactionId = formData.get('transaction_id') as string;
			const budgetEntryId = formData.get('budget_entry_id') as string;
			const createRules = formData.get('create_rules') === 'true';
			const amountTolerance = formData.get('amount_tolerance');

			if (!transactionId || !budgetEntryId) {
				return fail(400, { error: 'Transaction ID and Budget Entry ID are required' });
			}

			const teachData: TeachMatchRequest = {
				budget_entry_id: budgetEntryId,
				create_rules: createRules
			};

			if (amountTolerance) {
				teachData.amount_tolerance = parseFloat(amountTolerance as string);
			}

			await teachMatch(userId, transactionId, teachData);

			return {
				success: true,
				message: createRules
					? 'Transaction linked and matching rules created'
					: 'Transaction linked successfully'
			};
		} catch (error) {
			console.error('Error teaching match:', error);
			return fail(500, {
				error: error instanceof Error ? error.message : 'Failed to teach match'
			});
		}
	},

	// Bulk auto-match all unmatched transactions
	bulkAutoMatch: async ({ locals }) => {
		try {
			const userId = await getLocalUserId(locals.supabase);

			const result = await bulkAutoMatch(userId);

			return {
				success: true,
				matchedCount: result.matched_count,
				message: result.message
			};
		} catch (error) {
			console.error('Error bulk auto-matching:', error);
			return fail(500, {
				error: error instanceof Error ? error.message : 'Failed to auto-match transactions'
			});
		}
	}
};
