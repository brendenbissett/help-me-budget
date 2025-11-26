import type { RequestHandler } from '@sveltejs/kit';

const GO_API_URL = 'http://localhost:3000';

export const GET: RequestHandler = async ({ params }) => {
	const { provider } = params;

	// Call the Go API to start OAuth flow
	const response = await fetch(`${GO_API_URL}/auth/${provider}`, {
		redirect: 'manual',
	});

	// Get the auth URL from the redirect response
	if (response.status === 307 || response.status === 302) {
		const authUrl = response.headers.get('location');
		if (authUrl) {
			// Extract Set-Cookie headers from Go API response
			const setCookieHeader = response.headers.get('set-cookie');
			const headers: Record<string, string> = {
				Location: authUrl,
			};

			// Forward the session cookie from Go API to the browser
			if (setCookieHeader) {
				headers['Set-Cookie'] = setCookieHeader;
			}

			// Return a redirect response with the auth URL and cookies
			return new Response(null, {
				status: 302,
				headers,
			});
		}
	}

	throw new Error('Failed to get auth URL');
};
