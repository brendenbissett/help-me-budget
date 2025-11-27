-- Create roles table for RBAC
CREATE TABLE auth.roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create user_roles junction table (many-to-many)
CREATE TABLE auth.user_roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES auth.users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES auth.roles(id) ON DELETE CASCADE,
    granted_by UUID REFERENCES auth.users(id) ON DELETE SET NULL, -- Who granted this role
    granted_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, role_id) -- Prevent duplicate role assignments
);

-- Add is_active column to users table for enabling/disabling accounts
ALTER TABLE auth.users
    ADD COLUMN is_active BOOLEAN DEFAULT TRUE,
    ADD COLUMN deactivated_at TIMESTAMP WITH TIME ZONE,
    ADD COLUMN deactivated_by UUID REFERENCES auth.users(id) ON DELETE SET NULL;

-- Create audit_logs table for tracking admin actions
CREATE TABLE auth.audit_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    actor_id UUID REFERENCES auth.users(id) ON DELETE SET NULL, -- Who performed the action
    action VARCHAR(100) NOT NULL, -- 'user.disable', 'user.delete', 'session.kill', etc.
    resource_type VARCHAR(50) NOT NULL, -- 'user', 'session', 'role', etc.
    resource_id UUID, -- ID of the affected resource
    details JSONB, -- Additional context (e.g., {"reason": "spam", "ip": "1.2.3.4"})
    ip_address INET, -- IP address of the admin
    user_agent TEXT, -- Browser/client info
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Indexes for performance
CREATE INDEX idx_user_roles_user_id ON auth.user_roles(user_id);
CREATE INDEX idx_user_roles_role_id ON auth.user_roles(role_id);
CREATE INDEX idx_users_is_active ON auth.users(is_active);
CREATE INDEX idx_audit_logs_actor_id ON auth.audit_logs(actor_id);
CREATE INDEX idx_audit_logs_action ON auth.audit_logs(action);
CREATE INDEX idx_audit_logs_resource ON auth.audit_logs(resource_type, resource_id);
CREATE INDEX idx_audit_logs_created_at ON auth.audit_logs(created_at DESC);

-- Create updated_at trigger for roles
CREATE TRIGGER update_roles_updated_at
    BEFORE UPDATE ON auth.roles
    FOR EACH ROW
    EXECUTE FUNCTION auth.update_updated_at_column();

-- Insert default roles
INSERT INTO auth.roles (name, description) VALUES
    ('admin', 'Full system access - can manage users, sessions, and system settings'),
    ('user', 'Standard user access - can use the application features'),
    ('moderator', 'Limited admin access - can view audit logs and manage sessions');

-- Add comment explaining the RBAC design
COMMENT ON TABLE auth.roles IS 'System roles for role-based access control';
COMMENT ON TABLE auth.user_roles IS 'Junction table mapping users to their roles (many-to-many)';
COMMENT ON TABLE auth.audit_logs IS 'Audit trail for all administrative actions';
