import { createServerClient } from '@supabase/ssr';
import { SUPABASE_URL, SUPABASE_ANON_KEY } from '$env/static/private';
import type { Handle } from '@sveltejs/kit';

/**
 * Creates a Supabase client for server-side use with SvelteKit
 * This function is designed to be called from hooks.server.ts or +page.server.ts
 */
export const createSupabaseServerClient = (event: Parameters<Handle>[0]['event']) => {
	return createServerClient(SUPABASE_URL, SUPABASE_ANON_KEY, {
		cookies: {
			getAll: () => {
				return event.cookies.getAll();
			},
			setAll: (cookiesToSet) => {
				cookiesToSet.forEach(({ name, value, options }) => {
					event.cookies.set(name, value, { ...options, path: '/' });
				});
			}
		}
	});
};

/**
 * Helper function to get the current authenticated user
 * Uses getUser() instead of getSession() for security
 * This authenticates the session by contacting Supabase Auth server
 */
export const getServerUser = async (event: Parameters<Handle>[0]['event']) => {
	const supabase = createSupabaseServerClient(event);
	const {
		data: { user }
	} = await supabase.auth.getUser();
	return user;
};
