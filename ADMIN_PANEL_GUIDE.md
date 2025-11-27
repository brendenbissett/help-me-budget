# Admin Panel User Guide

Complete guide for using the Help-Me-Budget admin panel with RBAC (Role-Based Access Control).

## 🚀 Quick Start

### 1. Setup (First Time Only)

```bash
# 1. Ensure database and Redis are running
cd database
make up
make migrate-up

# 2. Start the API server
cd ../api
go run ./cmd/server

# 3. Start the frontend
cd ../frontend/help-me-budget
npm run dev

# 4. Log in with OAuth (Google or Facebook)
# Visit http://localhost:5173 and log in

# 5. Make yourself an admin
cd ../../database
# Edit scripts/make-admin.sql with your email
psql postgres://budgetuser:budgetpass@localhost:5432/help_me_budget < scripts/make-admin.sql
```

### 2. Access the Admin Panel

Visit: **http://localhost:5173/admin**

The admin panel will:
- Check if you're logged in
- Verify you have the `admin` role
- Redirect to home if not authorized

## 📊 Features

### Dashboard (`/admin`)

The main dashboard provides:
- **Total Users** - Count of all registered users
- **Active Sessions** - Number of currently logged-in users
- **Recent Audit Logs** - Count of recent administrative actions
- **Quick Actions** - Links to all admin functions

### User Management (`/admin/users`)

**View All Users:**
- See all registered users with their roles and status
- View user avatars, names, and emails
- Check last login timestamps

**User Actions:**

1. **Deactivate User**
   - Prevents user from logging in
   - Requires a reason (logged in audit trail)
   - User data is preserved
   - Can be reversed with "Reactivate"

2. **Reactivate User**
   - Re-enables a deactivated account
   - Logged in audit trail

3. **Delete User**
   - ⚠️ Permanently deletes user and all their data
   - Requires double confirmation
   - Requires a reason
   - **Cannot be undone**
   - Cascade deletes: budgets, transactions, OAuth links, roles

### Session Management (`/admin/sessions`)

**View Active Sessions:**
- See all active Redis sessions
- Each session shows the session key and encrypted data

**Kill Session:**
- Immediately logs out the user
- Session is deleted from Redis
- Action is logged in audit trail
- User must log in again

**Use Cases:**
- Force logout of suspicious accounts
- Clear stuck sessions
- Security incident response

### Audit Logs (`/admin/audit-logs`)

**View All Administrative Actions:**
- Every admin action is automatically logged
- Cannot be deleted or modified
- Paginated view (50 logs per page)

**Log Details Include:**
- **Timestamp** - When the action occurred
- **Action** - What was done (e.g., `user.deactivate`, `session.kill`)
- **Resource** - What was affected (user ID, session key)
- **Actor** - Who performed the action (admin user ID)
- **Details** - Additional context (reasons, parameters)
- **IP Address** - Where the action came from
- **User Agent** - Browser/client information

**Action Types:**
- `user.deactivate` - User account disabled
- `user.reactivate` - User account re-enabled
- `user.delete` - User permanently deleted
- `role.grant` - Role assigned to user
- `role.revoke` - Role removed from user
- `session.kill` - User session terminated

## 🔐 Roles & Permissions

### Admin Role
**Full system access:**
- ✅ View all users
- ✅ Deactivate/reactivate/delete users
- ✅ Grant and revoke roles
- ✅ View and kill sessions
- ✅ View audit logs

### Moderator Role
**Limited admin access:**
- ✅ View and kill sessions
- ✅ View audit logs
- ❌ Cannot manage users
- ❌ Cannot manage roles

### User Role
**Standard access:**
- ✅ Use the application
- ❌ No admin panel access

## 🛡️ Security Features

### Authentication Flow
1. User logs in via OAuth (Google/Facebook)
2. SvelteKit stores user data in HTTP-only cookie
3. Admin routes check user's roles in database
4. Only users with `admin` role can access `/admin/*`
5. All API calls include `X-User-ID` header for verification

### Audit Trail
- **Immutable** - Logs cannot be deleted or modified
- **Comprehensive** - Every action is logged automatically
- **Forensic** - IP address and user agent tracked
- **Compliance** - Supports security audits and investigations

### Session Security
- Sessions stored in Redis (not in-memory)
- 24-hour auto-expiration
- Admins can force logout via session kill
- All session kills are logged

## 📝 Common Tasks

### Promote a User to Admin

```bash
psql postgres://budgetuser:budgetpass@localhost:5432/help_me_budget

INSERT INTO auth.user_roles (user_id, role_id, granted_by)
SELECT u.id, r.id, u.id
FROM auth.users u, auth.roles r
WHERE u.email = 'user@example.com' AND r.name = 'admin';
```

### Revoke Admin Access

```bash
psql postgres://budgetuser:budgetpass@localhost:5432/help_me_budget

DELETE FROM auth.user_roles ur
USING auth.users u, auth.roles r
WHERE ur.user_id = u.id
AND ur.role_id = r.id
AND u.email = 'user@example.com'
AND r.name = 'admin';
```

### Check User's Roles

```bash
psql postgres://budgetuser:budgetpass@localhost:5432/help_me_budget

SELECT u.email, r.name as role, ur.granted_at
FROM auth.users u
INNER JOIN auth.user_roles ur ON ur.user_id = u.id
INNER JOIN auth.roles r ON r.id = ur.role_id
WHERE u.email = 'user@example.com';
```

### View Recent Admin Actions

```bash
psql postgres://budgetuser:budgetpass@localhost:5432/help_me_budget

SELECT
    al.created_at,
    al.action,
    al.resource_type,
    u.email as actor_email,
    al.details
FROM auth.audit_logs al
LEFT JOIN auth.users u ON u.id = al.actor_id
ORDER BY al.created_at DESC
LIMIT 20;
```

## 🚨 Best Practices

### User Management
- ✅ **Deactivate** instead of delete when possible (preserves data)
- ✅ Always provide clear reasons for deactivation/deletion
- ✅ Review user activity before taking action
- ✅ Use audit logs to track administrative actions
- ❌ Don't delete users without proper investigation

### Session Management
- ✅ Kill sessions for security incidents
- ✅ Review audit logs after killing sessions
- ❌ Don't kill sessions unnecessarily (disrupts user experience)

### Role Management
- ✅ Limit admin role to trusted individuals
- ✅ Use moderator role for limited admin access
- ✅ Document why users are granted admin access
- ❌ Don't grant admin to regular users

### Audit Logs
- ✅ Review logs regularly for suspicious activity
- ✅ Use logs for compliance and security audits
- ✅ Investigate unusual patterns
- ✅ Keep audit logs for regulatory requirements

## 🔧 Troubleshooting

### "Access Denied" when accessing /admin
**Problem:** You're not an admin
**Solution:** Have an existing admin grant you the admin role, or use the SQL script

### Admin panel loads but shows no data
**Problem:** API server not running or connection issue
**Solution:**
- Check if API server is running on port 3000
- Check browser console for errors
- Verify Redis and PostgreSQL are running

### Changes not reflected immediately
**Problem:** UI caching
**Solution:** Click the "Refresh" button on each page

### Session kill doesn't work
**Problem:** Session key format issue
**Solution:** Copy the exact session key from the UI

## 📚 Additional Resources

- **Database Schema:** See `database/ADMIN_SETUP.md`
- **API Documentation:** See `database/ADMIN_SETUP.md`
- **Migration Files:** `database/migrations/000004_*`
- **Backend Code:** `api/internal/admin/` and `api/internal/auth/rbac_repository.go`
- **Frontend Code:** `frontend/help-me-budget/src/routes/admin/`

## 🎯 Future Enhancements

Potential features for future development:
- [ ] Role creation/editing UI (currently requires SQL)
- [ ] User search and filtering
- [ ] Bulk user actions
- [ ] Export audit logs to CSV
- [ ] Real-time session monitoring
- [ ] Email notifications for admin actions
- [ ] Two-factor authentication for admins
- [ ] IP whitelisting for admin access

---

**Need Help?** Check the audit logs for debugging or review the API responses in browser dev tools.
