import { json } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';

// This endpoint logs out the user by clearing the user_data cookie
export const POST: RequestHandler = async ({ cookies }) => {
	try {
		cookies.delete('user_data', { path: '/' });
		return json({ success: true }, { status: 200 });
	} catch (error) {
		console.error('Error logging out:', error);
		return json({ error: 'Failed to logout' }, { status: 500 });
	}
};
