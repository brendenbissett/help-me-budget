import { redirect } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { authenticatedFetch } from '$lib/server/api-client';

export const GET: RequestHandler = async ({ url, locals: { supabase } }) => {
	const code = url.searchParams.get('code');
	const next = url.searchParams.get('next') ?? '/dashboard';
	const error_description = url.searchParams.get('error_description');
	const error_code = url.searchParams.get('error');

	// Handle errors from Supabase
	if (error_code) {
		console.error('Auth callback error:', error_code, error_description);
		throw redirect(303, `/auth/auth-code-error?error=${encodeURIComponent(error_description || error_code)}`);
	}

	if (code) {
		const { data, error } = await supabase.auth.exchangeCodeForSession(code);

		if (error) {
			console.error('Failed to exchange code for session:', error);

			// Handle PKCE error specifically - this is expected for email confirmations
			if (error.message?.includes('code verifier') || error.message?.includes('code_verifier')) {
				// Email was verified, but PKCE flow failed (expected behavior)
				// Redirect to login with a success message
				throw redirect(303, '/auth?verified=true');
			}

			throw redirect(303, `/auth/auth-code-error?error=${encodeURIComponent(error.message)}`);
		}

		if (data.user) {
			// Sync the Supabase user to local PostgreSQL database
			try {
				// Extract provider from app_metadata (Supabase stores this)
				// For email signups, provider will be 'email'
				const provider = data.user.app_metadata?.provider || 'email';
				const providerUserId = data.user.app_metadata?.provider_id || data.user.id;

				// Call Go API to sync user to local PostgreSQL
				const syncResponse = await authenticatedFetch('/auth/sync', {
					method: 'POST',
					body: JSON.stringify({
						email: data.user.email,
						name: data.user.user_metadata?.full_name || data.user.user_metadata?.name || data.user.email?.split('@')[0] || 'User',
						avatar_url: data.user.user_metadata?.avatar_url || data.user.user_metadata?.picture || '',
						provider: provider,
						provider_user_id: providerUserId
					})
				});

				if (!syncResponse.ok) {
					console.error('Failed to sync user to local PostgreSQL:', await syncResponse.text());
					// Don't fail the auth flow if sync fails - user is authenticated in Supabase
				}
			} catch (syncError) {
				console.error('Error syncing user to local PostgreSQL:', syncError);
				// Don't fail the auth flow if sync fails
			}

			// Redirect to dashboard or requested page
			throw redirect(303, next);
		}
	}

	// Return the user to an error page with some instructions
	throw redirect(303, '/auth/auth-code-error?error=No authentication code provided');
};
