-- Rollback migration: Remove Supabase auth sync trigger and function

-- Drop the trigger
DROP TRIGGER IF EXISTS on_auth_user_created ON auth.users;

-- Drop the function
DROP FUNCTION IF EXISTS public.handle_supabase_auth_user();

-- Revoke permissions (if needed)
-- Note: Be careful with this in production as it might affect other functionality
-- REVOKE SELECT, INSERT, UPDATE ON auth.users FROM service_role;
-- REVOKE SELECT, INSERT, UPDATE ON auth.user_oauth_providers FROM service_role;
