import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';

// This endpoint receives the user data from the Go API callback via query parameter
// The Go API redirects here with user data embedded in the URL after successful authentication
export const GET: RequestHandler = async ({ url, cookies }) => {
	const userParam = url.searchParams.get('user');
	if (!userParam) {
		return new Response('Missing user data', { status: 400 });
	}

	try {
		const userData = JSON.parse(userParam);

		// Store user data in a secure cookie
		cookies.set('user_data', JSON.stringify(userData), {
			path: '/',
			httpOnly: true,
			secure: false, // Set to true in production with HTTPS
			sameSite: 'lax',
			maxAge: 60 * 60 * 24, // 24 hours
		});
	} catch (error) {
		console.error('Error parsing or storing user data:', error);
		return new Response('Failed to process authentication', { status: 500 });
	}

	// Redirect to home page - moved outside try-catch so it's not caught
	redirect(302, '/');
};
