import { authenticatedFetchWithUser } from '../api-client';

export interface Category {
	id: string;
	user_id: string;
	name: string;
	category_type: 'income' | 'expense';
	color?: string | null;
	icon?: string | null;
	parent_category_id?: string | null;
	is_active: boolean;
	created_at: string;
	updated_at: string;
}

export interface CreateCategoryRequest {
	name: string;
	category_type: 'income' | 'expense';
	color?: string;
	icon?: string;
	parent_category_id?: string;
}

export interface UpdateCategoryRequest {
	name?: string;
	category_type?: 'income' | 'expense';
	color?: string;
	icon?: string;
	parent_category_id?: string;
	is_active?: boolean;
}

/**
 * Get all categories for a user
 * @param type - Optional filter by category type ('income' or 'expense')
 */
export async function getCategories(
	userId: string,
	type?: 'income' | 'expense'
): Promise<Category[]> {
	const endpoint = type ? `/api/categories?type=${type}` : '/api/categories';
	const response = await authenticatedFetchWithUser(endpoint, userId);

	if (!response.ok) {
		throw new Error(`Failed to fetch categories: ${response.statusText}`);
	}

	const data = await response.json();
	return data.categories || [];
}

/**
 * Get a specific category by ID
 */
export async function getCategory(userId: string, categoryId: string): Promise<Category> {
	const response = await authenticatedFetchWithUser(`/api/categories/${categoryId}`, userId);

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Category not found');
		}
		throw new Error(`Failed to fetch category: ${response.statusText}`);
	}

	return await response.json();
}

/**
 * Create a new category
 */
export async function createCategory(
	userId: string,
	category: CreateCategoryRequest
): Promise<Category> {
	const response = await authenticatedFetchWithUser('/api/categories', userId, {
		method: 'POST',
		body: JSON.stringify(category)
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to create category');
	}

	return await response.json();
}

/**
 * Update an existing category
 */
export async function updateCategory(
	userId: string,
	categoryId: string,
	updates: UpdateCategoryRequest
): Promise<Category> {
	const response = await authenticatedFetchWithUser(`/api/categories/${categoryId}`, userId, {
		method: 'PUT',
		body: JSON.stringify(updates)
	});

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Category not found');
		}
		const error = await response.json();
		throw new Error(error.error || 'Failed to update category');
	}

	return await response.json();
}

/**
 * Delete a category (soft delete)
 */
export async function deleteCategory(userId: string, categoryId: string): Promise<void> {
	const response = await authenticatedFetchWithUser(`/api/categories/${categoryId}`, userId, {
		method: 'DELETE'
	});

	if (!response.ok) {
		if (response.status === 404) {
			throw new Error('Category not found');
		}
		const error = await response.json();
		throw new Error(error.error || 'Failed to delete category');
	}
}

/**
 * Seed default categories for a new user
 */
export async function seedDefaultCategories(userId: string): Promise<Category[]> {
	const response = await authenticatedFetchWithUser('/api/categories/seed', userId, {
		method: 'POST'
	});

	if (!response.ok) {
		const error = await response.json();
		throw new Error(error.error || 'Failed to seed default categories');
	}

	const data = await response.json();
	return data.categories || [];
}
