import { json } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';

const API_URL = 'http://localhost:3000';

// This endpoint logs out the user by clearing the user_data cookie and Redis session
export const POST: RequestHandler = async ({ cookies }) => {
	try {
		// Get user data before deleting cookie
		const userCookie = cookies.get('user_data');
		if (userCookie) {
			try {
				const userData = JSON.parse(userCookie);
				const userId = userData.user_id;

				// Delete Redis session via Go API
				await fetch(`${API_URL}/auth/logout/${userId}`, {
					method: 'DELETE'
				});
			} catch (error) {
				console.error('Error deleting Redis session:', error);
				// Continue with logout even if Redis deletion fails
			}
		}

		// Delete cookie
		cookies.delete('user_data', { path: '/' });
		return json({ success: true }, { status: 200 });
	} catch (error) {
		console.error('Error logging out:', error);
		return json({ error: 'Failed to logout' }, { status: 500 });
	}
};
