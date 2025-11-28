import { API_SECRET_KEY } from '$env/static/private';

export const API_URL = 'http://localhost:3000';

/**
 * Authenticated fetch wrapper that automatically includes API key
 *
 * @param endpoint - API endpoint path (e.g., '/auth/sync')
 * @param options - Standard fetch options
 * @returns Promise<Response>
 */
export async function authenticatedFetch(
	endpoint: string,
	options: RequestInit = {}
): Promise<Response> {
	const headers = new Headers(options.headers);
	headers.set('X-API-Key', API_SECRET_KEY);
	headers.set('Content-Type', 'application/json');

	return globalThis.fetch(`${API_URL}${endpoint}`, {
		...options,
		headers
	});
}

/**
 * Authenticated fetch with user context
 * Includes both API key and user ID in request headers
 *
 * @param endpoint - API endpoint path (e.g., '/admin/users')
 * @param userId - Local PostgreSQL user ID
 * @param options - Standard fetch options
 * @returns Promise<Response>
 */
export async function authenticatedFetchWithUser(
	endpoint: string,
	userId: string,
	options: RequestInit = {}
): Promise<Response> {
	const headers = new Headers(options.headers);
	headers.set('X-API-Key', API_SECRET_KEY);
	headers.set('X-User-ID', userId);
	headers.set('Content-Type', 'application/json');

	return globalThis.fetch(`${API_URL}${endpoint}`, {
		...options,
		headers
	});
}
