import { json } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';

// This endpoint returns the current authenticated user from the stored cookie
export const GET: RequestHandler = async ({ cookies }) => {
	try {
		const userDataStr = cookies.get('user_data');
		if (!userDataStr) {
			return json({ user: null }, { status: 200 });
		}

		const userData = JSON.parse(userDataStr);
		return json({ user: userData }, { status: 200 });
	} catch (error) {
		console.error('Error retrieving user data:', error);
		return json({ user: null }, { status: 200 });
	}
};
