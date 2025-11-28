# Supabase Authentication Setup Guide

This guide will help you complete the Supabase authentication migration for Help-Me-Budget.

## Prerequisites

- Supabase project created at [https://app.supabase.com](https://app.supabase.com)
- Access to your Google and Facebook OAuth credentials

## Step 1: Get Supabase Credentials

1. Go to your Supabase project dashboard
2. Navigate to **Project Settings** → **API**
3. Copy the following values:
   - **Project URL** (e.g., `https://xxxxx.supabase.co`)
   - **Anon/Public Key** (starts with `eyJ...`)
   - **Service Role Key** (starts with `eyJ...`, keep this secret!)

## Step 2: Configure Frontend Environment Variables

1. Navigate to `frontend/help-me-budget/`
2. Create a `.env` file (or `.env.local`) based on `.env.example`:

```bash
# Public variables (safe to expose to the browser)
PUBLIC_SUPABASE_URL=https://your-project-ref.supabase.co
PUBLIC_SUPABASE_ANON_KEY=your-anon-key-here

# Private variables (server-side only)
SUPABASE_URL=https://your-project-ref.supabase.co
SUPABASE_ANON_KEY=your-anon-key-here
SUPABASE_SERVICE_ROLE_KEY=your-service-role-key-here
```

3. Replace the placeholder values with your actual Supabase credentials

## Step 3: Configure OAuth Providers in Supabase

### Google OAuth Setup

1. In Supabase Dashboard, go to **Authentication** → **Providers**
2. Find **Google** and click to enable
3. Enter your OAuth credentials:
   - **Client ID**: Your `GOOGLE_KEY` from `api/.env`
   - **Client Secret**: Your `GOOGLE_SECRET` from `api/.env`
4. Add authorized redirect URI in Google Cloud Console:
   - `https://YOUR-PROJECT-REF.supabase.co/auth/v1/callback`
5. Save changes

### Facebook OAuth Setup

1. In Supabase Dashboard, still in **Authentication** → **Providers**
2. Find **Facebook** and click to enable
3. Enter your OAuth credentials:
   - **Client ID**: Your `FACEBOOK_KEY` from `api/.env`
   - **Client Secret**: Your `FACEBOOK_SECRET` from `api/.env`
4. Add authorized redirect URI in Facebook Developers:
   - `https://YOUR-PROJECT-REF.supabase.co/auth/v1/callback`
5. Save changes

### Email Provider Setup

1. In Supabase Dashboard, still in **Authentication** → **Providers**
2. **Email** should be enabled by default
3. Optionally configure:
   - **Enable email confirmations** (recommended for production)
   - **Enable secure email change**
   - **Customize email templates** under **Authentication** → **Email Templates**

## Step 4: Run Database Migration

The sync trigger needs to be created in your Supabase database:

1. Start your local PostgreSQL (if testing locally):
   ```bash
   cd database && make up
   ```

2. Run the migration:
   ```bash
   cd database && make migrate-up
   ```

**OR** if using Supabase hosted database:

1. Go to **SQL Editor** in Supabase Dashboard
2. Copy the contents of `database/migrations/000005_sync_supabase_auth_users.up.sql`
3. Paste and execute the SQL

This creates a trigger that automatically syncs Supabase's `auth.users` to your custom `auth.users` table.

## Step 5: Configure Site URL and Redirect URLs

1. In Supabase Dashboard, go to **Authentication** → **URL Configuration**
2. Set **Site URL** to: `http://localhost:5173` (development) or your production URL
3. Add **Redirect URLs**:
   - `http://localhost:5173/auth/callback`
   - `http://localhost:5173/` (optional)
   - Add your production URLs when deploying

## Step 6: Test Authentication

1. Start the frontend development server:
   ```bash
   cd frontend/help-me-budget
   npm run dev
   ```

2. Navigate to `http://localhost:5173`

3. Test the following flows:
   - ✅ Sign up with email/password
   - ✅ Check email for confirmation link
   - ✅ Sign in with email/password
   - ✅ Sign in with Google OAuth
   - ✅ Sign in with Facebook OAuth
   - ✅ Sign in with magic link (passwordless)
   - ✅ Log out

## Step 7: Verify Database Sync

After a successful login, verify the sync trigger is working:

1. Connect to your database
2. Check that users appear in both tables:
   ```sql
   -- Supabase's auth.users (managed by Supabase)
   SELECT id, email, created_at FROM auth.users;

   -- Your custom auth.users (synced via trigger)
   SELECT id, email, name, created_at FROM auth.users;
   ```

3. If using OAuth, check `auth.user_oauth_providers`:
   ```sql
   SELECT user_id, provider, created_at FROM auth.user_oauth_providers;
   ```

## Step 8: (Optional) Remove Old Go API Auth Code

Once Supabase auth is working, you can:

1. Remove old auth routes from the Go API
2. Keep the Go API running for other endpoints (if needed)
3. Update `CLAUDE.md` and `README.md` documentation

## Troubleshooting

### "Invalid API key" error
- Check that your `SUPABASE_ANON_KEY` and `PUBLIC_SUPABASE_ANON_KEY` match
- Verify the key is copied correctly (no extra spaces)

### OAuth redirect not working
- Verify redirect URLs in Supabase Dashboard match OAuth provider settings
- Check Site URL is set correctly
- Ensure OAuth provider credentials are correct

### Email confirmation not received
- Check Supabase email rate limits (development projects have limits)
- For production, configure a custom SMTP server in Supabase

### Database sync not working
- Verify migration ran successfully
- Check Supabase logs for trigger errors
- Ensure database permissions are correct

## Security Notes

- **NEVER** commit `.env` files to version control
- **NEVER** expose `SUPABASE_SERVICE_ROLE_KEY` to the browser
- Use environment-specific URLs (localhost for dev, production domain for prod)
- Enable email confirmation in production
- Consider enabling MFA for admin accounts
- Review Supabase Row Level Security (RLS) policies

## Production Deployment

Before deploying to production:

1. Update Site URL to your production domain
2. Add production redirect URLs
3. Update OAuth provider redirect URIs
4. Enable email confirmation
5. Configure custom SMTP (optional but recommended)
6. Set up Supabase Edge Functions for advanced auth flows (if needed)
7. Review and enable RLS policies on custom tables
8. Set up monitoring and alerts

## Additional Resources

- [Supabase Auth Documentation](https://supabase.com/docs/guides/auth)
- [SvelteKit SSR Auth Guide](https://supabase.com/docs/guides/auth/server-side/sveltekit)
- [OAuth Provider Setup](https://supabase.com/docs/guides/auth/social-login)
