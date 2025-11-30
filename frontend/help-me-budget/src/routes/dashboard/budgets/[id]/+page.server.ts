import type { PageServerLoad, Actions } from './$types';
import { error, fail } from '@sveltejs/kit';
import { getLocalUserId } from '$lib/server/auth-helpers';
import {
	getBudgetWithEntries,
	getBudgetSummary,
	createBudgetEntry,
	updateBudgetEntry,
	deleteBudgetEntry,
	type CreateBudgetEntryRequest,
	type UpdateBudgetEntryRequest
} from '$lib/server/budget/budgets';

export const load: PageServerLoad = async ({ params, locals }) => {
	const session = await locals.safeGetSession();

	if (!session?.user) {
		throw error(401, 'Unauthorized');
	}

	try {
		const userId = await getLocalUserId(locals.supabase);
		const budgetId = params.id;

		// Fetch budget with entries and summary in parallel
		const [budgetWithEntries, summaryData] = await Promise.all([
			getBudgetWithEntries(userId, budgetId),
			getBudgetSummary(userId, budgetId)
		]);

		return {
			budget: {
				...budgetWithEntries,
				entries: budgetWithEntries.entries || []
			},
			summary: summaryData.summary,
			health: summaryData.health
		};
	} catch (err) {
		console.error('Error loading budget:', err);
		throw error(404, 'Budget not found');
	}
};

export const actions: Actions = {
	createEntry: async ({ request, params, locals }) => {
		const session = await locals.safeGetSession();

		if (!session?.user) {
			return fail(401, { error: 'Unauthorized' });
		}

		try {
			const formData = await request.formData();
			const userId = await getLocalUserId(locals.supabase);
			const budgetId = params.id;

			const name = formData.get('name') as string;
			const amount = parseFloat(formData.get('amount') as string);
			const entryType = formData.get('entry_type') as 'income' | 'expense';
			const frequency = formData.get('frequency') as string;
			const startDate = formData.get('start_date') as string;
			const description = formData.get('description') as string | null;
			const categoryId = formData.get('category_id') as string | null;
			const endDate = formData.get('end_date') as string | null;
			const dayOfMonth = formData.get('day_of_month') as string | null;
			const dayOfWeek = formData.get('day_of_week') as string | null;

			if (!name || isNaN(amount) || !entryType || !frequency || !startDate) {
				return fail(400, { error: 'Missing required fields' });
			}

			const entryData: CreateBudgetEntryRequest = {
				name,
				amount,
				entry_type: entryType,
				frequency: frequency as any,
				start_date: startDate
			};

			if (description) entryData.description = description;
			if (categoryId) entryData.category_id = categoryId;
			if (endDate) entryData.end_date = endDate;
			if (dayOfMonth) entryData.day_of_month = parseInt(dayOfMonth);
			if (dayOfWeek) entryData.day_of_week = parseInt(dayOfWeek);

			const entry = await createBudgetEntry(userId, budgetId, entryData);

			return { success: true, entry };
		} catch (err) {
			console.error('Error creating budget entry:', err);
			return fail(500, { error: 'Failed to create budget entry' });
		}
	},

	updateEntry: async ({ request, params, locals }) => {
		const session = await locals.safeGetSession();

		if (!session?.user) {
			return fail(401, { error: 'Unauthorized' });
		}

		try {
			const formData = await request.formData();
			const userId = await getLocalUserId(locals.supabase);
			const budgetId = params.id;
			const entryId = formData.get('entry_id') as string;

			if (!entryId) {
				return fail(400, { error: 'Entry ID is required' });
			}

			const updates: UpdateBudgetEntryRequest = {};

			const name = formData.get('name') as string | null;
			const amount = formData.get('amount') as string | null;
			const entryType = formData.get('entry_type') as string | null;
			const frequency = formData.get('frequency') as string | null;
			const startDate = formData.get('start_date') as string | null;
			const description = formData.get('description') as string | null;
			const categoryId = formData.get('category_id') as string | null;
			const endDate = formData.get('end_date') as string | null;
			const dayOfMonth = formData.get('day_of_month') as string | null;
			const dayOfWeek = formData.get('day_of_week') as string | null;
			const isActive = formData.get('is_active') as string | null;

			if (name) updates.name = name;
			if (amount) updates.amount = parseFloat(amount);
			if (entryType) updates.entry_type = entryType as 'income' | 'expense';
			if (frequency) updates.frequency = frequency as any;
			if (startDate) updates.start_date = startDate;
			if (description !== null) updates.description = description;
			if (categoryId) updates.category_id = categoryId;
			if (endDate !== null) updates.end_date = endDate;
			if (dayOfMonth) updates.day_of_month = parseInt(dayOfMonth);
			if (dayOfWeek) updates.day_of_week = parseInt(dayOfWeek);
			if (isActive !== null) updates.is_active = isActive === 'true';

			const entry = await updateBudgetEntry(userId, budgetId, entryId, updates);

			return { success: true, entry };
		} catch (err) {
			console.error('Error updating budget entry:', err);
			return fail(500, { error: 'Failed to update budget entry' });
		}
	},

	deleteEntry: async ({ request, params, locals }) => {
		const session = await locals.safeGetSession();

		if (!session?.user) {
			return fail(401, { error: 'Unauthorized' });
		}

		try {
			const formData = await request.formData();
			const userId = await getLocalUserId(locals.supabase);
			const budgetId = params.id;
			const entryId = formData.get('entry_id') as string;

			if (!entryId) {
				return fail(400, { error: 'Entry ID is required' });
			}

			await deleteBudgetEntry(userId, budgetId, entryId);

			return { success: true };
		} catch (err) {
			console.error('Error deleting budget entry:', err);
			return fail(500, { error: 'Failed to delete budget entry' });
		}
	}
};
