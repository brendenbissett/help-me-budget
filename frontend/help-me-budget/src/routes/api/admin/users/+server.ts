import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getLocalUserId } from '$lib/server/auth-helpers';
import { authenticatedFetchWithUser } from '$lib/server/api-client';

export const GET: RequestHandler = async ({ locals: { supabase }, url }) => {
	try {
		const localUserId = await getLocalUserId(supabase);

		// Forward request to Go API with local user ID header
		const limit = url.searchParams.get('limit') || '50';
		const offset = url.searchParams.get('offset') || '0';

		const response = await authenticatedFetchWithUser(`/admin/users?limit=${limit}&offset=${offset}`, localUserId, {
			method: 'GET'
		});

		const data = await response.json();

		if (!response.ok) {
			return json(data, { status: response.status });
		}

		return json(data);
	} catch (err: any) {
		console.error('Admin users API error:', err);
		if (err.status) {
			throw err;
		}
		throw error(500, 'Internal server error');
	}
};
