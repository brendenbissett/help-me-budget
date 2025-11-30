import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getLocalUserId } from '$lib/server/auth-helpers';
import { createAccount, type CreateAccountRequest } from '$lib/server/budget/accounts';

export const POST: RequestHandler = async ({ request, locals }) => {
	const session = await locals.safeGetSession();

	if (!session?.user) {
		throw error(401, 'Unauthorized');
	}

	try {
		const body = await request.json();
		const { name, account_type, balance, currency } = body;

		if (!name || !account_type) {
			throw error(400, 'Name and account type are required');
		}

		const userId = await getLocalUserId(locals.supabase);

		const accountData: CreateAccountRequest = {
			name,
			account_type,
			balance: balance || 0,
			currency: currency || 'USD'
		};

		const account = await createAccount(userId, accountData);

		return json({ success: true, account });
	} catch (err: any) {
		console.error('Error creating account during onboarding:', err);
		throw error(500, err.message || 'Failed to create account');
	}
};
