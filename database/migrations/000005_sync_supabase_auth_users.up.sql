-- Migration for Supabase Auth integration
--
-- NOTE: This migration is a placeholder for local development.
-- The actual Supabase sync trigger must be created in your Supabase project's SQL Editor.
--
-- For local development, your existing auth.users table will continue to work as-is.
-- When you deploy to production with Supabase, run the SQL in:
-- database/supabase/sync_trigger.sql in your Supabase dashboard.

-- This migration creates a note in the database about the Supabase integration
DO $$
BEGIN
    RAISE NOTICE 'Supabase Auth integration prepared.';
    RAISE NOTICE 'To complete setup, run database/supabase/sync_trigger.sql in Supabase SQL Editor.';
END $$;

-- No actual schema changes needed for local development
-- Your existing auth.users and auth.user_oauth_providers tables are already compatible
