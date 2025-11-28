import { json, error } from '@sveltejs/kit';
import type { RequestHandler } from './$types';
import { getLocalUserId } from '$lib/server/auth-helpers';
import { createClient } from '@supabase/supabase-js';
import { SUPABASE_URL, SUPABASE_SERVICE_ROLE_KEY } from '$env/static/private';

export const DELETE: RequestHandler = async ({ locals: { supabase }, params }) => {
	try {
		// Verify user is admin
		await getLocalUserId(supabase);

		// Extract user ID from session key (format: "session:user-id")
		const userId = params.key.replace('session:', '');

		// Use Supabase Admin API to delete the user's sessions
		const supabaseAdmin = createClient(SUPABASE_URL, SUPABASE_SERVICE_ROLE_KEY, {
			auth: {
				autoRefreshToken: false,
				persistSession: false
			}
		});

		// Delete all sessions for this user by updating their password
		// This invalidates all existing sessions
		// Note: We generate a random password since we're using OAuth (user won't use it)
		const randomPassword = crypto.randomUUID();
		const { error: updateError } = await supabaseAdmin.auth.admin.updateUserById(userId, {
			password: randomPassword
		});

		if (updateError) {
			console.error('Failed to invalidate sessions for user:', updateError);
			throw error(500, 'Failed to terminate session');
		}

		return json({
			success: true,
			message: 'Session terminated successfully'
		});
	} catch (err: any) {
		console.error('Admin kill session API error:', err);
		if (err.status) {
			throw err;
		}
		throw error(500, 'Internal server error');
	}
};
