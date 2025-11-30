import type { PageServerLoad } from './$types';
import { redirect } from '@sveltejs/kit';
import { getLocalUserId } from '$lib/server/auth-helpers';
import { getAccounts } from '$lib/server/budget/accounts';
import { getCategories } from '$lib/server/budget/categories';

export const load: PageServerLoad = async ({ locals }) => {
	const session = await locals.safeGetSession();

	if (!session?.user) {
		throw redirect(303, '/auth');
	}

	try {
		const userId = await getLocalUserId(locals.supabase);

		// Check if user has already completed onboarding
		const accounts = await getAccounts(userId);
		const categories = await getCategories(userId);

		// If user already has accounts or categories, redirect to dashboard
		if (accounts.length > 0 || categories.length > 0) {
			throw redirect(303, '/dashboard');
		}

		return {
			user: session.user
		};
	} catch (error) {
		console.error('Error loading onboarding:', error);
		// If there's an error but not a redirect, continue to onboarding
		if (error instanceof Response) {
			throw error;
		}
		return {
			user: session.user
		};
	}
};
