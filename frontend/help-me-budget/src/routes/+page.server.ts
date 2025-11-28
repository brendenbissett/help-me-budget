import { redirect } from '@sveltejs/kit';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ locals: { supabase } }) => {
	// Use getUser() instead of getSession() for security
	const {
		data: { user }
	} = await supabase.auth.getUser();

	// If user is already logged in, redirect to dashboard
	if (user) {
		throw redirect(303, '/dashboard');
	}

	return {};
};
