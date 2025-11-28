import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

const GO_API_URL = 'http://localhost:3000';

export const GET: RequestHandler = async ({ locals: { supabase } }) => {
	// Get authenticated user
	const {
		data: { user }
	} = await supabase.auth.getUser();

	if (!user) {
		throw error(401, 'Unauthorized');
	}

	try {
		// Call the Go API to get user roles by email
		// We use email because Supabase user ID != local PostgreSQL user ID
		const response = await globalThis.fetch(`${GO_API_URL}/auth/roles/by-email?email=${encodeURIComponent(user.email || '')}`, {
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		});

		if (!response.ok) {
			const errorText = await response.text();
			console.error('Failed to fetch user roles:', response.status, errorText);
			throw error(response.status, 'Failed to fetch user roles from API');
		}

		const data = await response.json();
		return json(data);
	} catch (err: any) {
		console.error('Error fetching user roles:', err);
		throw error(500, 'Failed to fetch user roles');
	}
};
