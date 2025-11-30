import type { PageServerLoad, Actions } from './$types';
import { getLocalUserId } from '$lib/server/auth-helpers';
import {
	getBudgets,
	createBudget,
	updateBudget,
	deleteBudget,
	type CreateBudgetRequest,
	type UpdateBudgetRequest
} from '$lib/server/budget/budgets';
import { fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
	try {
		const session = await locals.safeGetSession();

		if (!session?.user) {
			return {
				budgets: []
			};
		}

		const userId = await getLocalUserId(locals.supabase);
		const budgets = await getBudgets(userId);

		return {
			budgets: budgets || []
		};
	} catch (error) {
		console.error('Error loading budgets:', error);
		console.error('Error details:', error instanceof Error ? error.message : String(error));
		return {
			budgets: [],
			loadError: 'Failed to load budgets. Please try refreshing the page.'
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
			const description = formData.get('description') as string | null;

			if (!name) {
				return fail(400, { error: 'Budget name is required' });
			}

			const userId = await getLocalUserId(locals.supabase);

			const budgetData: CreateBudgetRequest = {
				name
			};

			if (description) {
				budgetData.description = description;
			}

			const budget = await createBudget(userId, budgetData);

			return { success: true, budget };
		} catch (error) {
			console.error('Error creating budget:', error);
			return fail(500, { error: 'Failed to create budget' });
		}
	},

	update: async ({ request, locals }) => {
		const session = await locals.safeGetSession();

		if (!session?.user) {
			return fail(401, { error: 'Unauthorized' });
		}

		try {
			const formData = await request.formData();
			const budgetId = formData.get('id') as string;
			const name = formData.get('name') as string | null;
			const description = formData.get('description') as string | null;
			const isActive = formData.get('is_active') as string | null;

			if (!budgetId) {
				return fail(400, { error: 'Budget ID is required' });
			}

			const userId = await getLocalUserId(locals.supabase);

			const updates: UpdateBudgetRequest = {};
			if (name) updates.name = name;
			if (description !== null) updates.description = description;
			if (isActive !== null) updates.is_active = isActive === 'true';

			const budget = await updateBudget(userId, budgetId, updates);

			return { success: true, budget };
		} catch (error) {
			console.error('Error updating budget:', error);
			return fail(500, { error: 'Failed to update budget' });
		}
	},

	delete: async ({ request, locals }) => {
		const session = await locals.safeGetSession();

		if (!session?.user) {
			return fail(401, { error: 'Unauthorized' });
		}

		try {
			const formData = await request.formData();
			const budgetId = formData.get('id') as string;

			if (!budgetId) {
				return fail(400, { error: 'Budget ID is required' });
			}

			const userId = await getLocalUserId(locals.supabase);
			await deleteBudget(userId, budgetId);

			return { success: true };
		} catch (error) {
			console.error('Error deleting budget:', error);
			return fail(500, { error: 'Failed to delete budget' });
		}
	}
};
