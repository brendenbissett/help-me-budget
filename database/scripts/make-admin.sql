-- Script to grant admin role to a user
-- Replace 'your-email@example.com' with the email of the user you want to make an admin

-- Step 1: Find the user by email and the admin role
DO $$
DECLARE
    v_user_id UUID;
    v_role_id UUID;
    v_user_email TEXT := 'your-email@example.com'; -- CHANGE THIS TO YOUR EMAIL
BEGIN
    -- Get user ID
    SELECT id INTO v_user_id
    FROM auth.users
    WHERE email = v_user_email;

    IF v_user_id IS NULL THEN
        RAISE EXCEPTION 'User with email % not found. Please log in first via OAuth.', v_user_email;
    END IF;

    -- Get admin role ID
    SELECT id INTO v_role_id
    FROM auth.roles
    WHERE name = 'admin';

    IF v_role_id IS NULL THEN
        RAISE EXCEPTION 'Admin role not found';
    END IF;

    -- Grant admin role to user
    INSERT INTO auth.user_roles (user_id, role_id, granted_by)
    VALUES (v_user_id, v_role_id, v_user_id) -- Self-granted for first admin
    ON CONFLICT (user_id, role_id) DO NOTHING;

    RAISE NOTICE 'Successfully granted admin role to user: %', v_user_email;
END $$;

-- Verify the role was granted
SELECT
    u.email,
    u.name,
    r.name as role,
    ur.granted_at
FROM auth.users u
INNER JOIN auth.user_roles ur ON ur.user_id = u.id
INNER JOIN auth.roles r ON r.id = ur.role_id
WHERE u.email = 'your-email@example.com'; -- CHANGE THIS TO YOUR EMAIL
