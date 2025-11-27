import { json } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';

const API_URL = 'http://localhost:3000';

// This endpoint returns the current authenticated user from the stored cookie
// and validates that the Redis session still exists
export const GET: RequestHandler = async ({ cookies }) => {
	try {
		const userDataStr = cookies.get('user_data');
		if (!userDataStr) {
			return json({ user: null }, { status: 200 });
		}

		const userData = JSON.parse(userDataStr);

		// Validate that the Redis session still exists
		const sessionCheck = await fetch(`${API_URL}/auth/session/${userData.user_id}`);
		if (!sessionCheck.ok) {
			// Session doesn't exist in Redis - clear cookie and return null
			cookies.delete('user_data', { path: '/' });
			return json({ user: null }, { status: 200 });
		}

		return json({ user: userData }, { status: 200 });
	} catch (error) {
		console.error('Error retrieving user data:', error);
		return json({ user: null }, { status: 200 });
	}
};
