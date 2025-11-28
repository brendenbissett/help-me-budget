-- ============================================================================
-- Supabase Auth → Custom Schema Sync Trigger
-- ============================================================================
--
-- PURPOSE: This SQL script syncs Supabase's auth.users to your custom auth.users table
--
-- WHEN TO RUN:
--   - Run this ONCE in your Supabase project's SQL Editor
--   - Located at: https://app.supabase.com/project/_/sql
--
-- WHAT IT DOES:
--   1. Creates your custom auth.users table (matching your local schema)
--   2. Creates auth.user_oauth_providers table
--   3. Sets up a trigger to automatically sync Supabase auth → your tables
--
-- ============================================================================

-- First, ensure we're working in the correct database
-- You should see this in the Supabase SQL Editor

-- ============================================================================
-- Step 1: Create custom auth.users table (if not exists)
-- ============================================================================

-- Note: Supabase has its own auth.users table managed by them
-- We're creating a SEPARATE table to store application-specific user data

CREATE TABLE IF NOT EXISTS public.app_users (
    id UUID PRIMARY KEY REFERENCES auth.users(id) ON DELETE CASCADE,
    email VARCHAR(255) NOT NULL UNIQUE,
    email_verified BOOLEAN DEFAULT FALSE,
    name VARCHAR(255),
    avatar_url TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_login_at TIMESTAMP WITH TIME ZONE
);

-- Create index for faster lookups
CREATE INDEX IF NOT EXISTS idx_app_users_email ON public.app_users(email);

-- ============================================================================
-- Step 2: Create OAuth providers table
-- ============================================================================

CREATE TABLE IF NOT EXISTS public.user_oauth_providers (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES public.app_users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL,
    provider_user_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_used_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(provider, provider_user_id)
);

CREATE INDEX IF NOT EXISTS idx_oauth_providers_user_id ON public.user_oauth_providers(user_id);
CREATE INDEX IF NOT EXISTS idx_oauth_providers_lookup ON public.user_oauth_providers(provider, provider_user_id);

-- ============================================================================
-- Step 3: Create updated_at trigger function
-- ============================================================================

CREATE OR REPLACE FUNCTION public.update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Create triggers to automatically update updated_at
DROP TRIGGER IF EXISTS update_app_users_updated_at ON public.app_users;
CREATE TRIGGER update_app_users_updated_at
    BEFORE UPDATE ON public.app_users
    FOR EACH ROW
    EXECUTE FUNCTION public.update_updated_at_column();

DROP TRIGGER IF EXISTS update_oauth_providers_updated_at ON public.user_oauth_providers;
CREATE TRIGGER update_oauth_providers_updated_at
    BEFORE UPDATE ON public.user_oauth_providers
    FOR EACH ROW
    EXECUTE FUNCTION public.update_updated_at_column();

-- ============================================================================
-- Step 4: Create sync function (Supabase auth.users → public.app_users)
-- ============================================================================

CREATE OR REPLACE FUNCTION public.handle_new_user()
RETURNS TRIGGER AS $$
DECLARE
    v_provider TEXT;
    v_provider_id TEXT;
BEGIN
    -- Extract provider information from Supabase auth metadata
    -- Supabase stores OAuth provider info in raw_app_meta_data.provider
    v_provider := COALESCE(NEW.raw_app_meta_data->>'provider', 'email');
    v_provider_id := NEW.id::TEXT;

    -- Insert or update the user in our custom app_users table
    INSERT INTO public.app_users (
        id,
        email,
        email_verified,
        name,
        avatar_url,
        created_at,
        updated_at,
        last_login_at
    ) VALUES (
        NEW.id,
        NEW.email,
        NEW.email_confirmed_at IS NOT NULL,
        COALESCE(
            NEW.raw_user_meta_data->>'full_name',
            NEW.raw_user_meta_data->>'name',
            NEW.raw_user_meta_data->>'user_name'
        ),
        COALESCE(
            NEW.raw_user_meta_data->>'avatar_url',
            NEW.raw_user_meta_data->>'picture'
        ),
        NEW.created_at,
        NEW.updated_at,
        NEW.last_sign_in_at
    )
    ON CONFLICT (id) DO UPDATE SET
        email = EXCLUDED.email,
        email_verified = EXCLUDED.email_verified,
        name = COALESCE(EXCLUDED.name, public.app_users.name),
        avatar_url = COALESCE(EXCLUDED.avatar_url, public.app_users.avatar_url),
        updated_at = EXCLUDED.updated_at,
        last_login_at = EXCLUDED.last_login_at;

    -- If this is an OAuth login (not email/password), sync to user_oauth_providers
    IF v_provider != 'email' THEN
        INSERT INTO public.user_oauth_providers (
            user_id,
            provider,
            provider_user_id,
            created_at,
            updated_at,
            last_used_at
        ) VALUES (
            NEW.id,
            v_provider,
            COALESCE(
                NEW.raw_user_meta_data->>'provider_id',
                NEW.raw_user_meta_data->>'sub',
                NEW.id::TEXT
            ),
            NEW.created_at,
            NEW.updated_at,
            NEW.last_sign_in_at
        )
        ON CONFLICT (provider, provider_user_id) DO UPDATE SET
            last_used_at = EXCLUDED.last_used_at,
            updated_at = EXCLUDED.updated_at;
    END IF;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql SECURITY DEFINER;

-- ============================================================================
-- Step 5: Create trigger on Supabase's auth.users
-- ============================================================================

-- Drop existing trigger if it exists
DROP TRIGGER IF EXISTS on_auth_user_created ON auth.users;

-- Create the trigger
CREATE TRIGGER on_auth_user_created
    AFTER INSERT OR UPDATE ON auth.users
    FOR EACH ROW
    EXECUTE FUNCTION public.handle_new_user();

-- ============================================================================
-- Step 6: Enable Row Level Security (RLS) - IMPORTANT for security!
-- ============================================================================

-- Enable RLS on custom tables
ALTER TABLE public.app_users ENABLE ROW LEVEL SECURITY;
ALTER TABLE public.user_oauth_providers ENABLE ROW LEVEL SECURITY;

-- Policy: Users can read their own data
CREATE POLICY "Users can view own profile" ON public.app_users
    FOR SELECT
    USING (auth.uid() = id);

-- Policy: Users can update their own data
CREATE POLICY "Users can update own profile" ON public.app_users
    FOR UPDATE
    USING (auth.uid() = id);

-- Policy: Users can view their own OAuth providers
CREATE POLICY "Users can view own oauth providers" ON public.user_oauth_providers
    FOR SELECT
    USING (auth.uid() = user_id);

-- Policy: Service role can do anything (for trigger)
CREATE POLICY "Service role has full access to app_users" ON public.app_users
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);

CREATE POLICY "Service role has full access to oauth_providers" ON public.user_oauth_providers
    FOR ALL
    TO service_role
    USING (true)
    WITH CHECK (true);

-- ============================================================================
-- Step 7: Grant necessary permissions
-- ============================================================================

-- Grant usage on schema
GRANT USAGE ON SCHEMA public TO authenticated, service_role;

-- Grant permissions on tables
GRANT SELECT, INSERT, UPDATE ON public.app_users TO authenticated, service_role;
GRANT SELECT ON public.user_oauth_providers TO authenticated, service_role;
GRANT INSERT, UPDATE ON public.user_oauth_providers TO service_role;

-- ============================================================================
-- Done!
-- ============================================================================

-- Add helpful comment
COMMENT ON FUNCTION public.handle_new_user() IS
'Automatically syncs Supabase auth.users to public.app_users table.
This trigger fires on INSERT or UPDATE to keep user data in sync.
OAuth provider information is extracted and stored in user_oauth_providers.';

-- Success message
DO $$
BEGIN
    RAISE NOTICE '✅ Supabase sync trigger installed successfully!';
    RAISE NOTICE 'Users will now be automatically synced to public.app_users';
END $$;
