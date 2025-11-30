import { redirect } from '@sveltejs/kit';
import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals: { supabase }, setHeaders }) => {
	// Prevent caching of user-specific data
	setHeaders({
		'cache-control': 'private, no-cache, no-store, must-revalidate',
		'pragma': 'no-cache',
		'expires': '0'
	});

	// Use getUser() to verify authentication
	const {
		data: { user }
	} = await supabase.auth.getUser();

	// If not authenticated, redirect to login
	if (!user) {
		throw redirect(303, '/auth');
	}

	return {
		user
	};
};
