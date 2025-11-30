import { authenticatedFetchWithUser } from '../api-client';

export interface Transaction {
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

export interface CreateTransactionRequest {
	account_id: string;
	category_id?: string;
	amount: number;
	transaction_type: 'income' | 'expense';
	description?: string;
	transaction_date: string;
	notes?: string;
}

export interface UpdateTransactionRequest {
	account_id?: string;
	category_id?: string;
	budget_entry_id?: string;
	amount?: number;
	transaction_type?: 'income' | 'expense';
	description?: string;
	transaction_date?: string;
	notes?: string;
	match_confidence?: 'manual' | 'auto_high' | 'auto_low' | 'unmatched';
}

export interface TransactionFilters {
	account_id?: string;
	category_id?: string;
	start_date?: string;
	end_date?: string;
}

/**
 * Get all transactions for a user with optional filters
 */
export async function getTransactions(
	userId: string,
	filters?: TransactionFilters
): Promise<Transaction[]> {
	let url = '/api/transactions';
	const params = new URLSearchParams();

	if (filters?.account_id) params.append('account_id', filters.account_id);
	if (filters?.category_id) params.append('category_id', filters.category_id);
	if (filters?.start_date) params.append('start_date', filters.start_date);
	if (filters?.end_date) params.append('end_date', filters.end_date);

	if (params.toString()) {
		url += `?${params.toString()}`;
	}

	const response = await authenticatedFetchWithUser(url, userId);

	if (!response.ok) {
		throw new Error(`Failed to fetch transactions: ${response.statusText}`);
	}

	const data = await response.json();
	return data.transactions || [];
}

/**
 * Get a specific transaction by ID
 */
export async function getTransaction(userId: string, transactionId: string): Promise<Transaction> {
	const response = await authenticatedFetchWithUser(`/api/transactions/${transactionId}`, userId);

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Transaction not found');
		}
		throw new Error(`Failed to fetch transaction: ${response.statusText}`);
	}

	return await response.json();
}

/**
 * Create a new transaction
 */
export async function createTransaction(
	userId: string,
	transaction: CreateTransactionRequest
): Promise<Transaction> {
	const response = await authenticatedFetchWithUser('/api/transactions', userId, {
		method: 'POST',
		body: JSON.stringify(transaction)
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to create transaction');
	}

	return await response.json();
}

/**
 * Update an existing transaction
 */
export async function updateTransaction(
	userId: string,
	transactionId: string,
	updates: UpdateTransactionRequest
): Promise<Transaction> {
	const response = await authenticatedFetchWithUser(`/api/transactions/${transactionId}`, userId, {
		method: 'PUT',
		body: JSON.stringify(updates)
	});

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Transaction not found');
		}
		const error = await response.json();
		throw new Error(error.error || 'Failed to update transaction');
	}

	return await response.json();
}

/**
 * Delete a transaction
 */
export async function deleteTransaction(userId: string, transactionId: string): Promise<void> {
	const response = await authenticatedFetchWithUser(`/api/transactions/${transactionId}`, userId, {
		method: 'DELETE'
	});

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Transaction not found');
		}
		const error = await response.json();
		throw new Error(error.error || 'Failed to delete transaction');
	}
}

/**
 * Get unmatched transactions (not linked to budget entries)
 */
export async function getUnmatchedTransactions(userId: string): Promise<Transaction[]> {
	const response = await authenticatedFetchWithUser('/api/transactions/unmatched', userId);

	if (!response.ok) {
		throw new Error(`Failed to fetch unmatched transactions: ${response.statusText}`);
	}

	const data = await response.json();
	return data.transactions || [];
}

/**
 * Categorize a transaction
 */
export async function categorizeTransaction(
	userId: string,
	transactionId: string,
	categoryId: string
): Promise<Transaction> {
	const response = await authenticatedFetchWithUser(
		`/api/transactions/${transactionId}/categorize`,
		userId,
		{
			method: 'POST',
			body: JSON.stringify({ category_id: categoryId })
		}
	);

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Transaction not found');
		}
		const error = await response.json();
		throw new Error(error.error || 'Failed to categorize transaction');
	}

	return await response.json();
}

/**
 * Link a transaction to a budget entry
 */
export async function linkTransactionToBudgetEntry(
	userId: string,
	transactionId: string,
	budgetEntryId: string,
	matchConfidence: 'manual' | 'auto_high' | 'auto_low' = 'manual'
): Promise<Transaction> {
	const response = await authenticatedFetchWithUser(
		`/api/transactions/${transactionId}/link`,
		userId,
		{
			method: 'POST',
			body: JSON.stringify({
				budget_entry_id: budgetEntryId,
				match_confidence: matchConfidence
			})
		}
	);

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Transaction not found');
		}
		const error = await response.json();
		throw new Error(error.error || 'Failed to link transaction');
	}

	return await response.json();
}
