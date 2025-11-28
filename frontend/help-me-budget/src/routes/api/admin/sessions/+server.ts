import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getLocalUserId } from '$lib/server/auth-helpers';
import { createClient } from '@supabase/supabase-js';
import { SUPABASE_URL, SUPABASE_SERVICE_ROLE_KEY } from '$env/static/private';

export const GET: RequestHandler = async ({ locals: { supabase } }) => {
	try {
		// Verify user is admin
		await getLocalUserId(supabase);

		// Use Supabase Admin API to list sessions
		const supabaseAdmin = createClient(SUPABASE_URL, SUPABASE_SERVICE_ROLE_KEY, {
			auth: {
				autoRefreshToken: false,
				persistSession: false
			}
		});

		// Get all users with their last sign in time
		// Note: Supabase doesn't expose active sessions directly, but we can show users
		const { data: users, error: usersError } = await supabaseAdmin.auth.admin.listUsers();

		if (usersError) {
			console.error('Failed to list users from Supabase:', usersError);
			throw error(500, 'Failed to fetch sessions');
		}

		// Transform to match the expected session format
		const sessions = users.users
			.filter((user) => user.last_sign_in_at) // Only show users who have logged in
			.map((user) => ({
				key: `session:${user.id}`,
				user_id: user.id,
				email: user.email,
				name: user.user_metadata?.full_name || user.user_metadata?.name || user.email,
				provider: user.app_metadata?.provider || 'email',
				login_at: user.last_sign_in_at,
				// Note: We can't get exact session expiry from Supabase without service role
				// Sessions typically last 1 hour by default
			}))
			.sort((a, b) => new Date(b.login_at).getTime() - new Date(a.login_at).getTime());

		return json({
			sessions,
			total: sessions.length
		});
	} catch (err: any) {
		console.error('Admin sessions API error:', err);
		if (err.status) {
			throw err;
		}
		throw error(500, 'Internal server error');
	}
};
