import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

const API_URL = 'http://localhost:3000';

export const GET: RequestHandler = async ({ cookies }) => {
	const userCookie = cookies.get('user_data');
	if (!userCookie) {
		return json({ error: 'Unauthorized' }, { status: 401 });
	}

	try {
		const userData = JSON.parse(userCookie);
		const userId = userData.user_id;

		const response = await fetch(`${API_URL}/admin/sessions`, {
			method: 'GET',
			headers: {
				'X-User-ID': userId,
				'Content-Type': 'application/json'
			}
		});

		const data = await response.json();

		if (!response.ok) {
			return json(data, { status: response.status });
		}

		return json(data);
	} catch (error) {
		console.error('Admin sessions API error:', error);
		return json({ error: 'Internal server error' }, { status: 500 });
	}
};
