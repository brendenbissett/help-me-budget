import { json } from '@sveltejs/kit';
import type { RequestHandler } from './$types';

const API_URL = 'http://localhost:3000';

export const GET: RequestHandler = async ({ cookies, url }) => {
	const userCookie = cookies.get('user_data');
	if (!userCookie) {
		return json({ error: 'Unauthorized' }, { status: 401 });
	}

	try {
		const userData = JSON.parse(userCookie);
		const userId = userData.user_id;

		const limit = url.searchParams.get('limit') || '50';
		const offset = url.searchParams.get('offset') || '0';

		const response = await fetch(`${API_URL}/admin/audit-logs?limit=${limit}&offset=${offset}`, {
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
		console.error('Admin audit logs API error:', error);
		return json({ error: 'Internal server error' }, { status: 500 });
	}
};
