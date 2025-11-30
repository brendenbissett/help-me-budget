import type { PageServerLoad } from './$types';
import { getLocalUserId } from '$lib/server/auth-helpers';
import { getDashboardSummary } from '$lib/server/budget/dashboard';

export const load: PageServerLoad = async ({ locals }) => {
	const session = await locals.safeGetSession();

	if (!session?.user) {
		return {
			summary: null
		};
	}

	try {
		const userId = await getLocalUserId(locals.supabase);
		const summary = await getDashboardSummary(userId);

		return {
			summary
		};
	} catch (error) {
		console.error('Error loading dashboard:', error);
		return {
			summary: null,
			error: 'Failed to load dashboard data'
		};
	}
};
