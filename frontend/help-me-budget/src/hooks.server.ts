import type { Handle } from '@sveltejs/kit';
import { paraglideMiddleware } from '$lib/paraglide/server';
import { createSupabaseServerClient } from '$lib/supabase.server';
import { checkApiHealth } from '$lib/server/api-client';
import { sequence } from '@sveltejs/kit/hooks';

// Check API health on server startup
let apiHealthChecked = false;
let isApiAvailable = false;

async function checkApiOnStartup() {
	if (!apiHealthChecked) {
		apiHealthChecked = true;
		isApiAvailable = await checkApiHealth();

		if (isApiAvailable) {
			console.log('✅ API server is available at http://localhost:3000');
		} else {
			console.error('❌ API server is NOT available at http://localhost:3000');
			console.error('Please ensure the Go API server is running: cd api && go run ./cmd/server');
		}
	}
}

const handleParaglide: Handle = ({ event, resolve }) => paraglideMiddleware(event.request, ({ request, locale }) => {
	event.request = request;

	return resolve(event, {
		transformPageChunk: ({ html }) => html.replace('%paraglide.lang%', locale)
	});
});

const handleSupabase: Handle = async ({ event, resolve }) => {
	// Check API health on first request
	await checkApiOnStartup();

	// Create a Supabase client for this request
	event.locals.supabase = createSupabaseServerClient(event);

	// Add safeGetSession helper
	event.locals.safeGetSession = async () => {
		const {
			data: { session },
		} = await event.locals.supabase.auth.getSession();
		if (!session) {
			return { session: null, user: null };
		}

		const {
			data: { user },
			error,
		} = await event.locals.supabase.auth.getUser();
		if (error) {
			return { session: null, user: null };
		}

		return { session, user };
	};

	return resolve(event, {
		filterSerializedResponseHeaders(name) {
			// Allow Supabase auth cookies to be sent to the client
			return name === 'content-range' || name === 'x-supabase-api-version';
		}
	});
};

// Combine multiple handlers using sequence
export const handle: Handle = sequence(handleSupabase, handleParaglide);
