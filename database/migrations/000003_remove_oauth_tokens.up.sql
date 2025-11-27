-- Remove OAuth token storage (not needed for auth-only flows)
ALTER TABLE auth.user_oauth_providers
    DROP COLUMN IF EXISTS access_token,
    DROP COLUMN IF EXISTS refresh_token,
    DROP COLUMN IF EXISTS token_expiry;
