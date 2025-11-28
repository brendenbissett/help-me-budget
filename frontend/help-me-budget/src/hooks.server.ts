import type { Handle } from '@sveltejs/kit';
import { paraglideMiddleware } from '$lib/paraglide/server';
import { createSupabaseServerClient } from '$lib/supabase.server';
import { sequence } from '@sveltejs/kit/hooks';

const handleParaglide: Handle = ({ event, resolve }) => paraglideMiddleware(event.request, ({ request, locale }) => {
	event.request = request;

	return resolve(event, {
		transformPageChunk: ({ html }) => html.replace('%paraglide.lang%', locale)
	});
});

const handleSupabase: Handle = async ({ event, resolve }) => {
	// Create a Supabase client for this request
	event.locals.supabase = createSupabaseServerClient(event);

	return resolve(event, {
		filterSerializedResponseHeaders(name) {
			// Allow Supabase auth cookies to be sent to the client
			return name === 'content-range' || name === 'x-supabase-api-version';
		}
	});
};

// Combine multiple handlers using sequence
export const handle: Handle = sequence(handleSupabase, handleParaglide);
