import type { PageServerLoad, Actions } from './$types';
import { getLocalUserId } from '$lib/server/auth-helpers';
import {
	getCategories,
	createCategory,
	updateCategory,
	deleteCategory,
	seedDefaultCategories,
	type Category,
	type CreateCategoryRequest,
	type UpdateCategoryRequest
} from '$lib/server/budget/categories';
import { fail } from '@sveltejs/kit';

export const load: PageServerLoad = async ({ locals }) => {
	const session = await locals.safeGetSession();

	if (!session?.user) {
		return {
			categories: []
		};
	}

	try {
		const userId = await getLocalUserId(locals.supabase);
		const categories = await getCategories(userId);

		return {
			categories
		};
	} catch (error) {
		console.error('Error loading categories:', error);
		return {
			categories: [],
			error: 'Failed to load categories'
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
			const categoryType = formData.get('category_type') as string;
			const color = formData.get('color') as string | null;
			const icon = formData.get('icon') as string | null;
			const parentCategoryId = formData.get('parent_category_id') as string | null;

			if (!name || !categoryType) {
				return fail(400, { error: 'Name and category type are required' });
			}

			const userId = await getLocalUserId(locals.supabase);

			const categoryData: CreateCategoryRequest = {
				name,
				category_type: categoryType as any
			};

			if (color) categoryData.color = color;
			if (icon) categoryData.icon = icon;
			if (parentCategoryId) categoryData.parent_category_id = parentCategoryId;

			const category = await createCategory(userId, categoryData);

			return { success: true, category };
		} catch (error) {
			console.error('Error creating category:', error);
			return fail(500, { error: 'Failed to create category' });
		}
	},

	update: async ({ request, locals }) => {
		const session = await locals.safeGetSession();

		if (!session?.user) {
			return fail(401, { error: 'Unauthorized' });
		}

		try {
			const formData = await request.formData();
			const categoryId = formData.get('id') as string;
			const name = formData.get('name') as string | null;
			const categoryType = formData.get('category_type') as string | null;
			const color = formData.get('color') as string | null;
			const icon = formData.get('icon') as string | null;

			if (!categoryId) {
				return fail(400, { error: 'Category ID is required' });
			}

			const userId = await getLocalUserId(locals.supabase);

			const updates: UpdateCategoryRequest = {};
			if (name) updates.name = name;
			if (categoryType) updates.category_type = categoryType as any;
			if (color) updates.color = color;
			if (icon) updates.icon = icon;

			const category = await updateCategory(userId, categoryId, updates);

			return { success: true, category };
		} catch (error) {
			console.error('Error updating category:', error);
			return fail(500, { error: 'Failed to update category' });
		}
	},

	delete: async ({ request, locals }) => {
		const session = await locals.safeGetSession();

		if (!session?.user) {
			return fail(401, { error: 'Unauthorized' });
		}

		try {
			const formData = await request.formData();
			const categoryId = formData.get('id') as string;

			if (!categoryId) {
				return fail(400, { error: 'Category ID is required' });
			}

			const userId = await getLocalUserId(locals.supabase);
			await deleteCategory(userId, categoryId);

			return { success: true };
		} catch (error) {
			console.error('Error deleting category:', error);
			return fail(500, { error: 'Failed to delete category' });
		}
	},

	seedDefaults: async ({ locals }) => {
		const session = await locals.safeGetSession();

		if (!session?.user) {
			return fail(401, { error: 'Unauthorized' });
		}

		try {
			const userId = await getLocalUserId(locals.supabase);
			const categories = await seedDefaultCategories(userId);

			return { success: true, seeded: true, categories };
		} catch (error) {
			console.error('Error seeding categories:', error);
			return fail(500, { error: 'Failed to seed default categories' });
		}
	}
};
