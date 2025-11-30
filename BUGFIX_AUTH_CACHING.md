# Bugfix: Auth/Session Caching Issue

## Problem
When logging out and logging in as a different user, using the browser's back button would show the previous user's data. This is a critical security and UX issue.

## Root Cause
1. **Browser Caching**: SvelteKit was allowing browser to cache user-specific pages
2. **Incomplete Logout**: signOut() was only clearing Supabase session but not invalidating cached page data
3. **Back Button**: Browser's back/forward cache (bfcache) was serving stale authenticated pages

## Solution Implemented

### 1. Force Full Page Reload on Logout
**File**: `src/lib/supabase.client.ts`

Updated `signOut()` function to:
- Sign out from Supabase
- Force a full page reload to `/` using `window.location.href`
- This clears all cached data and forces fresh server requests

```typescript
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
```

### 2. Add No-Cache Headers to Authenticated Pages
**Files Updated**:
- `src/routes/dashboard/+layout.server.ts`
- `src/routes/dashboard/accounts/+page.server.ts`
- `src/routes/dashboard/categories/+page.server.ts`

Added cache control headers to prevent browser caching:
```typescript
setHeaders({
	'cache-control': 'private, no-cache, no-store, must-revalidate',
	'pragma': 'no-cache',
	'expires': '0'
});
```

**What these headers do**:
- `private`: Only browser can cache, not CDN/proxies
- `no-cache`: Must revalidate with server before using cached copy
- `no-store`: Don't store this in cache at all
- `must-revalidate`: Can't serve stale content
- `pragma: no-cache`: HTTP/1.0 backward compatibility
- `expires: 0`: Immediately expired

### 3. Simplified Logout Handler
**File**: `src/routes/dashboard/+layout.svelte`

Removed redundant `goto('/')` since `signOut()` now handles the redirect:
```typescript
async function handleLogout() {
	try {
		// signOut() will handle redirect and cache invalidation
		await signOut();
	} catch (err) {
		console.error('Error logging out:', err);
	}
}
```

## Testing

To verify the fix:
1. Log in as User A
2. Navigate to Accounts or Categories
3. Log out
4. Log in as User B
5. Navigate to Accounts or Categories
6. Click browser back button
7. **Expected**: Should stay on User B's page or redirect to login
8. **Before Fix**: Would show User A's cached data

## Technical Details

### Why Full Page Reload?
- SvelteKit's client-side navigation can keep data in memory
- `invalidateAll()` only works for server load functions
- Client-side state (like Svelte stores) needs to be cleared
- Full reload ensures completely clean slate

### Why These Specific Headers?
- Multiple headers for maximum browser compatibility
- Prevents all forms of caching (browser, proxy, CDN)
- Works across HTTP/1.0, HTTP/1.1, and HTTP/2
- Critical for pages with user-specific data

### Alternative Approaches Considered
1. ❌ **Just `invalidateAll()`**: Doesn't clear client-side state
2. ❌ **Session-based cache keys**: Complex and error-prone
3. ✅ **Full reload + no-cache headers**: Simple and bulletproof

## Security Implications

This fix addresses:
- **Privacy**: Previous user's data not visible to new user
- **Security**: Prevents unauthorized data access via browser history
- **Compliance**: Required for apps handling sensitive financial data

## Files Changed
- `src/lib/supabase.client.ts` (signOut function)
- `src/routes/dashboard/+layout.svelte` (handleLogout)
- `src/routes/dashboard/+layout.server.ts` (cache headers)
- `src/routes/dashboard/accounts/+page.server.ts` (cache headers)
- `src/routes/dashboard/categories/+page.server.ts` (cache headers)

## Status
✅ **Fixed and Ready for Testing**

---

**Note**: All future authenticated pages should include the same no-cache headers in their `+page.server.ts` or `+layout.server.ts` files.
