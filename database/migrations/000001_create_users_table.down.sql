-- Drop triggers
DROP TRIGGER IF EXISTS update_oauth_providers_updated_at ON auth.user_oauth_providers;
DROP TRIGGER IF EXISTS update_users_updated_at ON auth.users;

-- Drop trigger function
DROP FUNCTION IF EXISTS auth.update_updated_at_column();

-- Drop indexes
DROP INDEX IF EXISTS auth.idx_oauth_providers_lookup;
DROP INDEX IF EXISTS auth.idx_oauth_providers_user_id;
DROP INDEX IF EXISTS auth.idx_users_email;

-- Drop tables
DROP TABLE IF EXISTS auth.user_oauth_providers;
DROP TABLE IF EXISTS auth.users;
