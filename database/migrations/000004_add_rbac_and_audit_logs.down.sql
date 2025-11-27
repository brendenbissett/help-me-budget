-- Drop triggers
DROP TRIGGER IF EXISTS update_roles_updated_at ON auth.roles;

-- Drop indexes
DROP INDEX IF EXISTS auth.idx_audit_logs_created_at;
DROP INDEX IF EXISTS auth.idx_audit_logs_resource;
DROP INDEX IF EXISTS auth.idx_audit_logs_action;
DROP INDEX IF EXISTS auth.idx_audit_logs_actor_id;
DROP INDEX IF EXISTS auth.idx_users_is_active;
DROP INDEX IF EXISTS auth.idx_user_roles_role_id;
DROP INDEX IF EXISTS auth.idx_user_roles_user_id;

-- Drop tables
DROP TABLE IF EXISTS auth.audit_logs;

-- Remove columns from users table
ALTER TABLE auth.users
    DROP COLUMN IF EXISTS deactivated_by,
    DROP COLUMN IF EXISTS deactivated_at,
    DROP COLUMN IF EXISTS is_active;

-- Drop junction table and roles table
DROP TABLE IF EXISTS auth.user_roles;
DROP TABLE IF EXISTS auth.roles;
