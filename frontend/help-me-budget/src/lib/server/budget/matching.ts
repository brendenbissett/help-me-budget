import { authenticatedFetchWithUser } from '../api-client';

// ============================================================================
// Types
// ============================================================================

export interface MatchSuggestion {
	budget_entry: {
		id: string;
		budget_id: string;
		name: string;
		amount: number;
		entry_type: 'income' | 'expense';
		frequency: string;
		category_id: string | null;
		is_active: boolean;
	};
	confidence_score: number;
	confidence_level: 'auto_high' | 'auto_low';
	match_reasons: string[];
}

export interface MatchSuggestionsResponse {
	transaction: {
		id: string;
		user_id: string;
		account_id: string;
		transaction_type: 'income' | 'expense';
		amount: number;
		transaction_date: string;
		description: string | null;
		category_id: string | null;
		budget_entry_id: string | null;
		match_confidence: string | null;
	};
	suggestions: MatchSuggestion[];
}

export interface AutoMatchResponse {
	transaction: {
		id: string;
		user_id: string;
		account_id: string;
		transaction_type: 'income' | 'expense';
		amount: number;
		transaction_date: string;
		description: string | null;
		category_id: string | null;
		budget_entry_id: string | null;
		match_confidence: string | null;
	};
	matched: boolean;
}

export interface BulkAutoMatchResponse {
	matched_count: number;
	message: string;
}

export interface TeachMatchRequest {
	budget_entry_id: string;
	create_rules: boolean;
	amount_tolerance?: number;
}

export interface TeachMatchResponse {
	transaction: {
		id: string;
		user_id: string;
		account_id: string;
		transaction_type: 'income' | 'expense';
		amount: number;
		transaction_date: string;
		description: string | null;
		category_id: string | null;
		budget_entry_id: string | null;
		match_confidence: string | null;
	};
	rules_created: boolean;
}

export interface MatchingRules {
	description_contains?: string[];
	merchant_name?: string;
	amount_tolerance?: number;
}

// ============================================================================
// API Functions
// ============================================================================

/**
 * Get match suggestions for a transaction
 */
export async function getMatchSuggestions(
	userId: string,
	transactionId: string
): Promise<MatchSuggestionsResponse> {
	const response = await authenticatedFetchWithUser(
		`/api/matching/suggestions/${transactionId}`,
		userId,
		{
			method: 'GET'
		}
	);

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to get match suggestions');
	}

	return response.json();
}

/**
 * Auto-match a single transaction (only if confidence >= 70%)
 */
export async function autoMatchTransaction(
	userId: string,
	transactionId: string
): Promise<AutoMatchResponse> {
	const response = await authenticatedFetchWithUser(
		`/api/matching/auto-match/${transactionId}`,
		userId,
		{
			method: 'POST'
		}
	);

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to auto-match transaction');
	}

	return response.json();
}

/**
 * Auto-match all unmatched transactions
 */
export async function bulkAutoMatch(userId: string): Promise<BulkAutoMatchResponse> {
	const response = await authenticatedFetchWithUser(`/api/matching/bulk-auto-match`, userId, {
		method: 'POST'
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to bulk auto-match transactions');
	}

	return response.json();
}

/**
 * Teach mode: Link transaction to budget entry and optionally create matching rules
 */
export async function teachMatch(
	userId: string,
	transactionId: string,
	data: TeachMatchRequest
): Promise<TeachMatchResponse> {
	const response = await authenticatedFetchWithUser(
		`/api/matching/teach/${transactionId}`,
		userId,
		{
			method: 'POST',
			body: JSON.stringify(data)
		}
	);

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to teach match');
	}

	return response.json();
}

/**
 * Update matching rules for a budget entry
 */
export async function updateBudgetEntryMatchingRules(
	userId: string,
	budgetId: string,
	entryId: string,
	rules: MatchingRules
): Promise<void> {
	const response = await authenticatedFetchWithUser(
		`/api/budgets/${budgetId}/entries/${entryId}/matching-rules`,
		userId,
		{
			method: 'POST',
			body: JSON.stringify({ matching_rules: rules })
		}
	);

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to update matching rules');
	}
}
