import type { PageServerLoad, Actions } from './$types';
import { getLocalUserId } from '$lib/server/auth-helpers';
import {
	getAccounts,
	createAccount,
	updateAccount,
	deleteAccount,
	type Account,
	type CreateAccountRequest,
	type UpdateAccountRequest
} from '$lib/server/budget/accounts';
import { fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
	const session = await locals.safeGetSession();

	if (!session?.user) {
		return {
			accounts: []
		};
	}

	try {
		const userId = await getLocalUserId(locals.supabase);
		const accounts = await getAccounts(userId);

		return {
			accounts
		};
	} catch (error) {
		console.error('Error loading accounts:', error);
		return {
			accounts: [],
			error: 'Failed to load accounts'
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
			const name = formData.get('name') as string;
			const accountType = formData.get('account_type') as string;
			const balance = parseFloat(formData.get('balance') as string) || 0;
			const currency = (formData.get('currency') as string) || 'USD';

			if (!name || !accountType) {
				return fail(400, { error: 'Name and account type are required' });
			}

			const userId = await getLocalUserId(locals.supabase);

			const accountData: CreateAccountRequest = {
				name,
				account_type: accountType as any,
				balance,
				currency
			};

			const account = await createAccount(userId, accountData);

			return { success: true, account };
		} catch (error) {
			console.error('Error creating account:', error);
			return fail(500, { error: 'Failed to create account' });
		}
	},

	update: async ({ request, locals }) => {
		const session = await locals.safeGetSession();

		if (!session?.user) {
			return fail(401, { error: 'Unauthorized' });
		}

		try {
			const formData = await request.formData();
			const accountId = formData.get('id') as string;
			const name = formData.get('name') as string | null;
			const accountType = formData.get('account_type') as string | null;
			const balanceStr = formData.get('balance') as string | null;
			const currency = formData.get('currency') as string | null;

			if (!accountId) {
				return fail(400, { error: 'Account ID is required' });
			}

			const userId = await getLocalUserId(locals.supabase);

			const updates: UpdateAccountRequest = {};
			if (name) updates.name = name;
			if (accountType) updates.account_type = accountType as any;
			if (balanceStr !== null) updates.balance = parseFloat(balanceStr);
			if (currency) updates.currency = currency;

			const account = await updateAccount(userId, accountId, updates);

			return { success: true, account };
		} catch (error) {
			console.error('Error updating account:', error);
			return fail(500, { error: 'Failed to update account' });
		}
	},

	delete: async ({ request, locals }) => {
		const session = await locals.safeGetSession();

		if (!session?.user) {
			return fail(401, { error: 'Unauthorized' });
		}

		try {
			const formData = await request.formData();
			const accountId = formData.get('id') as string;

			if (!accountId) {
				return fail(400, { error: 'Account ID is required' });
			}

			const userId = await getLocalUserId(locals.supabase);
			await deleteAccount(userId, accountId);

			return { success: true };
		} catch (error) {
			console.error('Error deleting account:', error);
			return fail(500, { error: 'Failed to delete account' });
		}
	}
};
