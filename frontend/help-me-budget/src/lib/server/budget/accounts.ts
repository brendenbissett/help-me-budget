import { authenticatedFetchWithUser } from '../api-client';

export interface Account {
	id: string;
	user_id: string;
	name: string;
	account_type: 'checking' | 'savings' | 'credit_card' | 'cash' | 'investment';
	balance: number;
	currency: string;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface CreateAccountRequest {
	name: string;
	account_type: 'checking' | 'savings' | 'credit_card' | 'cash' | 'investment';
	balance: number;
	currency: string;
}

export interface UpdateAccountRequest {
	name?: string;
	account_type?: 'checking' | 'savings' | 'credit_card' | 'cash' | 'investment';
	balance?: number;
	currency?: string;
	is_active?: boolean;
}

/**
 * Get all accounts for a user
 */
export async function getAccounts(userId: string): Promise<Account[]> {
	const response = await authenticatedFetchWithUser('/api/accounts', userId);

	if (!response.ok) {
		throw new Error(`Failed to fetch accounts: ${response.statusText}`);
	}

	const data = await response.json();
	return data.accounts || [];
}

/**
 * Get a specific account by ID
 */
export async function getAccount(userId: string, accountId: string): Promise<Account> {
	const response = await authenticatedFetchWithUser(`/api/accounts/${accountId}`, userId);

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Account not found');
		}
		throw new Error(`Failed to fetch account: ${response.statusText}`);
	}

	return await response.json();
}

/**
 * Create a new account
 */
export async function createAccount(
	userId: string,
	account: CreateAccountRequest
): Promise<Account> {
	const response = await authenticatedFetchWithUser('/api/accounts', userId, {
		method: 'POST',
		body: JSON.stringify(account)
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to create account');
	}

	return await response.json();
}

/**
 * Update an existing account
 */
export async function updateAccount(
	userId: string,
	accountId: string,
	updates: UpdateAccountRequest
): Promise<Account> {
	const response = await authenticatedFetchWithUser(`/api/accounts/${accountId}`, userId, {
		method: 'PUT',
		body: JSON.stringify(updates)
	});

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Account not found');
		}
		const error = await response.json();
		throw new Error(error.error || 'Failed to update account');
	}

	return await response.json();
}

/**
 * Delete an account (soft delete)
 */
export async function deleteAccount(userId: string, accountId: string): Promise<void> {
	const response = await authenticatedFetchWithUser(`/api/accounts/${accountId}`, userId, {
		method: 'DELETE'
	});

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Account not found');
		}
		const error = await response.json();
		throw new Error(error.error || 'Failed to delete account');
	}
}

/**
 * Get total balance across all active accounts
 */
export async function getTotalBalance(
	userId: string
): Promise<{ total_balance: number; currency: string }> {
	const response = await authenticatedFetchWithUser('/api/accounts/balance/total', userId);

	if (!response.ok) {
		throw new Error(`Failed to fetch total balance: ${response.statusText}`);
	}

	return await response.json();
}
