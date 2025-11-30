import type { PageServerLoad, Actions } from './$types';
import { fail } from '@sveltejs/kit';
import { getLocalUserId } from '$lib/server/auth-helpers';
import {
	getTransactions,
	createTransaction,
	updateTransaction,
	deleteTransaction,
	categorizeTransaction,
	type CreateTransactionRequest,
	type UpdateTransactionRequest,
	type TransactionFilters
} from '$lib/server/budget/transactions';
import { getAccounts } from '$lib/server/budget/accounts';
import { getCategories } from '$lib/server/budget/categories';

export const load: PageServerLoad = async ({ locals, url }) => {
	const session = await locals.safeGetSession();

	if (!session?.user) {
		return {
			transactions: [],
			accounts: [],
			categories: []
		};
	}

	try {
		const userId = await getLocalUserId(locals.supabase);

		// Parse filters from URL query params
		const filters: TransactionFilters = {};
		const accountId = url.searchParams.get('account_id');
		const categoryId = url.searchParams.get('category_id');
		const startDate = url.searchParams.get('start_date');
		const endDate = url.searchParams.get('end_date');

		if (accountId) filters.account_id = accountId;
		if (categoryId) filters.category_id = categoryId;
		if (startDate) filters.start_date = startDate;
		if (endDate) filters.end_date = endDate;

		// Fetch transactions, accounts, and categories in parallel
		const [transactions, accounts, categories] = await Promise.all([
			getTransactions(userId, filters),
			getAccounts(userId),
			getCategories(userId)
		]);

		return {
			transactions: transactions || [],
			accounts: accounts || [],
			categories: categories || []
		};
	} catch (error) {
		console.error('Error loading transactions:', error);
		return {
			transactions: [],
			accounts: [],
			categories: [],
			loadError: 'Failed to load transactions. Please try refreshing the page.'
		};
	}
};

export const actions: Actions = {
	create: async ({ request, locals }) => {
		const session = await locals.safeGetSession();

		if (!session?.user) {
			return fail(401, { error: 'Unauthorized' });
		}

		try {
			const formData = await request.formData();
			const userId = await getLocalUserId(locals.supabase);

			const accountId = formData.get('account_id') as string;
			const amount = parseFloat(formData.get('amount') as string);
			const transactionType = formData.get('transaction_type') as 'income' | 'expense';
			const transactionDate = formData.get('transaction_date') as string;
			const categoryId = formData.get('category_id') as string | null;
			const description = formData.get('description') as string | null;
			const notes = formData.get('notes') as string | null;

			if (!accountId || isNaN(amount) || !transactionType || !transactionDate) {
				return fail(400, { error: 'Missing required fields' });
			}

			const transactionData: CreateTransactionRequest = {
				account_id: accountId,
				amount,
				transaction_type: transactionType,
				transaction_date: transactionDate
			};

			if (categoryId) transactionData.category_id = categoryId;
			if (description) transactionData.description = description;
			if (notes) transactionData.notes = notes;

			const transaction = await createTransaction(userId, transactionData);

			return { success: true, transaction };
		} catch (err) {
			console.error('Error creating transaction:', err);
			return fail(500, { error: 'Failed to create transaction' });
		}
	},

	update: async ({ request, locals }) => {
		const session = await locals.safeGetSession();

		if (!session?.user) {
			return fail(401, { error: 'Unauthorized' });
		}

		try {
			const formData = await request.formData();
			const userId = await getLocalUserId(locals.supabase);
			const transactionId = formData.get('transaction_id') as string;

			if (!transactionId) {
				return fail(400, { error: 'Transaction ID is required' });
			}

			const updates: UpdateTransactionRequest = {};

			const accountId = formData.get('account_id') as string | null;
			const categoryId = formData.get('category_id') as string | null;
			const amount = formData.get('amount') as string | null;
			const transactionType = formData.get('transaction_type') as string | null;
			const transactionDate = formData.get('transaction_date') as string | null;
			const description = formData.get('description') as string | null;
			const notes = formData.get('notes') as string | null;

			if (accountId) updates.account_id = accountId;
			if (categoryId) updates.category_id = categoryId;
			if (amount) updates.amount = parseFloat(amount);
			if (transactionType) updates.transaction_type = transactionType as 'income' | 'expense';
			if (transactionDate) updates.transaction_date = transactionDate;
			if (description !== null) updates.description = description;
			if (notes !== null) updates.notes = notes;

			const transaction = await updateTransaction(userId, transactionId, updates);

			return { success: true, transaction };
		} catch (err) {
			console.error('Error updating transaction:', err);
			return fail(500, { error: 'Failed to update transaction' });
		}
	},

	delete: async ({ request, locals }) => {
		const session = await locals.safeGetSession();

		if (!session?.user) {
			return fail(401, { error: 'Unauthorized' });
		}

		try {
			const formData = await request.formData();
			const userId = await getLocalUserId(locals.supabase);
			const transactionId = formData.get('transaction_id') as string;

			if (!transactionId) {
				return fail(400, { error: 'Transaction ID is required' });
			}

			await deleteTransaction(userId, transactionId);

			return { success: true };
		} catch (err) {
			console.error('Error deleting transaction:', err);
			return fail(500, { error: 'Failed to delete transaction' });
		}
	},

	categorize: async ({ request, locals }) => {
		const session = await locals.safeGetSession();

		if (!session?.user) {
			return fail(401, { error: 'Unauthorized' });
		}

		try {
			const formData = await request.formData();
			const userId = await getLocalUserId(locals.supabase);
			const transactionId = formData.get('transaction_id') as string;
			const categoryId = formData.get('category_id') as string;

			if (!transactionId || !categoryId) {
				return fail(400, { error: 'Transaction ID and Category ID are required' });
			}

			const transaction = await categorizeTransaction(userId, transactionId, categoryId);

			return { success: true, transaction };
		} catch (err) {
			console.error('Error categorizing transaction:', err);
			return fail(500, { error: 'Failed to categorize transaction' });
		}
	}
};
