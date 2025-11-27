# Admin Panel Setup Guide

This guide will help you set up the admin panel with RBAC (Role-Based Access Control).

## Overview

The admin panel allows privileged users to:
- **User Management**: View, deactivate, reactivate, and delete users
- **Role Management**: Grant and revoke roles
- **Session Management**: View active sessions and kill specific sessions
- **Audit Logs**: View all administrative actions

## Roles

Three roles are automatically created:
- **admin** - Full system access (can manage users, roles, sessions)
- **moderator** - Limited admin access (can view audit logs and manage sessions)
- **user** - Standard user access (default for all OAuth users)

## Setup Steps

### 1. Run the RBAC Migration

The migration creates the necessary tables and default roles:

```bash
cd database
make migrate-up
```

This creates:
- `auth.roles` - System roles
- `auth.user_roles` - User-to-role mappings
- `auth.audit_logs` - Audit trail for admin actions
- Adds `is_active`, `deactivated_at`, `deactivated_by` to `auth.users`

### 2. Log in with OAuth

First, log in to the application using Google or Facebook OAuth. This creates your user record in the database.

### 3. Grant Yourself Admin Role

After logging in, run the SQL script to make yourself an admin:

```bash
# Edit the script to use your email
nano scripts/make-admin.sql  # Change 'your-email@example.com' to your actual email

# Run the script
psql postgres://budgetuser:budgetpass@localhost:5432/help_me_budget < scripts/make-admin.sql
```

Or connect to the database directly:

```bash
psql postgres://budgetuser:budgetpass@localhost:5432/help_me_budget

-- Replace with your email
INSERT INTO auth.user_roles (user_id, role_id, granted_by)
SELECT u.id, r.id, u.id
FROM auth.users u, auth.roles r
WHERE u.email = 'your-email@example.com' AND r.name = 'admin';

-- Verify
SELECT u.email, r.name as role
FROM auth.users u
INNER JOIN auth.user_roles ur ON ur.user_id = u.id
INNER JOIN auth.roles r ON r.id = ur.role_id
WHERE u.email = 'your-email@example.com';
```

## API Endpoints

All admin endpoints require authentication and the `admin` role (except audit logs which also allow `moderator`).

### User Management

```
GET    /admin/users                    - List all users (paginated)
POST   /admin/users/:id/deactivate     - Deactivate a user account
POST   /admin/users/:id/reactivate     - Reactivate a user account
DELETE /admin/users/:id                - Permanently delete a user
POST   /admin/users/:id/roles/grant    - Grant a role to a user
POST   /admin/users/:id/roles/revoke   - Revoke a role from a user
```

### Session Management

```
GET    /admin/sessions          - List all active Redis sessions
DELETE /admin/sessions/:key     - Kill a specific session
```

### Audit Logs

```
GET    /admin/audit-logs        - List audit logs (paginated, admin or moderator)
```

## Authentication Flow

The admin panel uses the existing OAuth authentication:

1. User logs in via Google/Facebook (creates user in database)
2. SvelteKit stores user data in HTTP-only cookie
3. SvelteKit admin routes extract user ID from cookie
4. SvelteKit proxies requests to Go API with `X-User-ID` header
5. Go API middleware validates user ID and checks admin role
6. If user has admin role, request proceeds

## Example API Calls

**List all users:**
```bash
curl -X GET http://localhost:3000/admin/users \
  -H "X-User-ID: <your-user-id-uuid>"
```

**Deactivate a user:**
```bash
curl -X POST http://localhost:3000/admin/users/<user-id>/deactivate \
  -H "X-User-ID: <your-user-id-uuid>" \
  -H "Content-Type: application/json" \
  -d '{"reason": "Spam account"}'
```

**Grant admin role:**
```bash
curl -X POST http://localhost:3000/admin/users/<user-id>/roles/grant \
  -H "X-User-ID: <your-user-id-uuid>" \
  -H "Content-Type: application/json" \
  -d '{"role_name": "admin"}'
```

**Kill a session:**
```bash
curl -X DELETE http://localhost:3000/admin/sessions/session:abc123 \
  -H "X-User-ID: <your-user-id-uuid>"
```

## Audit Logging

All admin actions are automatically logged to `auth.audit_logs` with:
- Actor (who performed the action)
- Action type (e.g., `user.deactivate`, `session.kill`)
- Resource type and ID
- Additional details (JSONB)
- IP address and user agent
- Timestamp

## Security Considerations

- Admin role grants full system access - be careful who you promote
- All admin actions are logged and cannot be deleted
- Deactivated users cannot log in but their data is preserved
- Deleted users are permanently removed (cascade delete)
- Session killing immediately logs out the user

## Frontend Integration (Coming Next)

The SvelteKit frontend will have:
- `/admin` - Admin dashboard
- `/admin/users` - User management interface
- `/admin/sessions` - Active sessions viewer
- `/admin/audit-logs` - Audit log viewer

The frontend will:
1. Check if current user has admin role
2. Render admin navigation if authorized
3. Proxy all requests to Go API with `X-User-ID` header
4. Handle 401/403 errors gracefully
