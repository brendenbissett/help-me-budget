import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getLocalUserId } from '$lib/server/auth-helpers';
import { authenticatedFetchWithUser } from '$lib/server/api-client';

export const POST: RequestHandler = async ({ locals: { supabase }, params, request }) => {
	try {
		const localUserId = await getLocalUserId(supabase);
		const body = await request.json();

		const response = await authenticatedFetchWithUser(`/admin/users/${params.id}/deactivate`, localUserId, {
			method: 'POST',
			body: JSON.stringify(body)
		});

		const data = await response.json();

		if (!response.ok) {
			return json(data, { status: response.status });
		}

		return json(data);
	} catch (err: any) {
		console.error('Admin deactivate user API error:', err);
		if (err.status) {
			throw err;
		}
		throw error(500, 'Internal server error');
	}
};
