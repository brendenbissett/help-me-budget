import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getLocalUserId } from '$lib/server/auth-helpers';
import { seedDefaultCategories } from '$lib/server/budget/categories';

export const POST: RequestHandler = async ({ locals }) => {
	const session = await locals.safeGetSession();

	if (!session?.user) {
		throw error(401, 'Unauthorized');
	}

	try {
		const userId = await getLocalUserId(locals.supabase);
		const categories = await seedDefaultCategories(userId);

		return json({ success: true, categories });
	} catch (err: any) {
		console.error('Error seeding categories during onboarding:', err);
		throw error(500, err.message || 'Failed to seed categories');
	}
};
