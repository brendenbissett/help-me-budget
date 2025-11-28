import type { LayoutServerLoad } from './$types';

export const load: LayoutServerLoad = async ({ locals: { supabase } }) => {
	// Use getUser() instead of getSession() for security
	// This authenticates the session by contacting Supabase Auth server
	const {
		data: { user }
	} = await supabase.auth.getUser();

	return {
		user
	};
};
