import { error } from '@sveltejs/kit';
import type { SupabaseClient } from '@supabase/supabase-js';

const API_URL = 'http://localhost:3000';

/**
 * Gets the local PostgreSQL user ID for a Supabase-authenticated user
 * This bridges the gap between Supabase auth (external) and local database (internal)
 */
export async function getLocalUserId(supabase: SupabaseClient): Promise<string> {
	// Get authenticated user from Supabase
	const {
		data: { user }
	} = await supabase.auth.getUser();

	if (!user || !user.email) {
		throw error(401, 'Unauthorized');
	}

	// Lookup user in local PostgreSQL by email
	const rolesResponse = await globalThis.fetch(
		`${API_URL}/auth/roles/by-email?email=${encodeURIComponent(user.email)}`,
		{
			method: 'GET',
			headers: {
				'Content-Type': 'application/json'
			}
		}
	);

	if (!rolesResponse.ok) {
		throw error(403, 'User not found in local database');
	}

	const rolesData = await rolesResponse.json();
	return rolesData.user_id;
}
