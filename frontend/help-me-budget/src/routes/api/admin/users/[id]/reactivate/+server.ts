import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getLocalUserId } from '$lib/server/auth-helpers';

const API_URL = 'http://localhost:3000';

export const POST: RequestHandler = async ({ locals: { supabase }, params }) => {
	try {
		const localUserId = await getLocalUserId(supabase);

		const response = await globalThis.fetch(`${API_URL}/admin/users/${params.id}/reactivate`, {
			method: 'POST',
			headers: {
				'X-User-ID': localUserId,
				'Content-Type': 'application/json'
			}
		});

		const data = await response.json();

		if (!response.ok) {
			return json(data, { status: response.status });
		}

		return json(data);
	} catch (err: any) {
		console.error('Admin reactivate user API error:', err);
		if (err.status) {
			throw err;
		}
		throw error(500, 'Internal server error');
	}
};
