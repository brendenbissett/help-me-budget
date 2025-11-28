import { json } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';

// This endpoint returns the current authenticated user from Supabase
export const GET: RequestHandler = async ({ locals: { supabase } }) => {
	try {
		const {
			data: { user }
		} = await supabase.auth.getUser();

		if (!user) {
			return json({ user: null }, { status: 200 });
		}

		// Return user data in a format compatible with the old structure
		return json(
			{
				user: {
					email: user.email,
					name: user.user_metadata?.full_name || user.user_metadata?.name || user.email,
					avatar_url: user.user_metadata?.avatar_url || user.user_metadata?.picture || '',
					provider: user.app_metadata?.provider || 'email'
				}
			},
			{ status: 200 }
		);
	} catch (error) {
		console.error('Error retrieving user data:', error);
		return json({ user: null }, { status: 200 });
	}
};
