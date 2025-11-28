import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals: { supabase } }) => {
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
