import { json } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';

// This endpoint logs out the user by clearing Supabase session
// Note: Session management is now fully handled by Supabase
export const POST: RequestHandler = async ({ locals: { supabase } }) => {
	try {
		// Sign out from Supabase - this is the only source of truth for sessions
		const { error: signOutError } = await supabase.auth.signOut();
		if (signOutError) {
			console.error('Error signing out from Supabase:', signOutError);
			return json({ error: 'Failed to logout' }, { status: 500 });
		}

		return json({ success: true }, { status: 200 });
	} catch (error) {
		console.error('Error logging out:', error);
		return json({ error: 'Failed to logout' }, { status: 500 });
	}
};
