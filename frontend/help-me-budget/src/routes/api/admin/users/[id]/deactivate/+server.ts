import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getLocalUserId } from '$lib/server/auth-helpers';

const API_URL = 'http://localhost:3000';

export const POST: RequestHandler = async ({ locals: { supabase }, params, request }) => {
	try {
		const localUserId = await getLocalUserId(supabase);
		const body = await request.json();

		const response = await globalThis.fetch(`${API_URL}/admin/users/${params.id}/deactivate`, {
			method: 'POST',
			headers: {
				'X-User-ID': localUserId,
				'Content-Type': 'application/json'
			},
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
