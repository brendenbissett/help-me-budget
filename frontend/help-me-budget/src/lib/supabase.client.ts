import { createBrowserClient } from '@supabase/ssr';
import { PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY } from '$env/static/public';

/**
 * Creates a Supabase client for browser-side use
 * This is safe to call multiple times - it will return the same client instance
 */
export const createSupabaseBrowserClient = () => {
	return createBrowserClient(PUBLIC_SUPABASE_URL, PUBLIC_SUPABASE_ANON_KEY);
};

/**
 * Helper function to sign in with OAuth provider (Google, Facebook, etc.)
 */
export const signInWithOAuth = async (provider: 'google' | 'facebook') => {
	const supabase = createSupabaseBrowserClient();
	const { data, error } = await supabase.auth.signInWithOAuth({
		provider,
		options: {
			redirectTo: `${window.location.origin}/auth/callback`
		}
	});

	if (error) {
		throw error;
	}

	return data;
};

/**
 * Helper function to sign up with email and password
 */
export const signUpWithEmail = async (email: string, password: string, metadata?: any) => {
	const supabase = createSupabaseBrowserClient();
	const { data, error } = await supabase.auth.signUp({
		email,
		password,
		options: {
			emailRedirectTo: `${window.location.origin}/auth/callback`,
			data: metadata
		}
	});

	if (error) {
		throw error;
	}

	return data;
};

/**
 * Helper function to sign in with email and password
 */
export const signInWithEmail = async (email: string, password: string) => {
	const supabase = createSupabaseBrowserClient();
	const { data, error } = await supabase.auth.signInWithPassword({
		email,
		password
	});

	if (error) {
		throw error;
	}

	return data;
};

/**
 * Helper function to sign in with magic link (passwordless)
 */
export const signInWithMagicLink = async (email: string) => {
	const supabase = createSupabaseBrowserClient();
	const { data, error } = await supabase.auth.signInWithOtp({
		email,
		options: {
			emailRedirectTo: `${window.location.origin}/auth/callback`
		}
	});

	if (error) {
		throw error;
	}

	return data;
};

/**
 * Helper function to sign out
 * Clears all authentication and invalidates cached data
 */
export const signOut = async () => {
	const supabase = createSupabaseBrowserClient();
	const { error } = await supabase.auth.signOut();

	if (error) {
		throw error;
	}

	// Clear any cached data by invalidating all routes
	if (typeof window !== 'undefined') {
		// Force a full page reload to clear all cached data
		window.location.href = '/';
	}
};

/**
 * Helper function to get the current user
 * Uses getUser() instead of getSession() for security
 * This authenticates the session by contacting Supabase Auth server
 */
export const getUser = async () => {
	const supabase = createSupabaseBrowserClient();
	const {
		data: { user }
	} = await supabase.auth.getUser();
	return user;
};
