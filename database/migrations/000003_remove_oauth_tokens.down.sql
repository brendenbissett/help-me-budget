-- Restore OAuth token columns (for rollback)
ALTER TABLE auth.user_oauth_providers
    ADD COLUMN IF NOT EXISTS access_token TEXT,
    ADD COLUMN IF NOT EXISTS refresh_token TEXT,
    ADD COLUMN IF NOT EXISTS token_expiry TIMESTAMP WITH TIME ZONE;
