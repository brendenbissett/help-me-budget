-- Create separate schemas for auth and budget data
CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS budget;

-- Grant usage on schemas to the database user
GRANT USAGE ON SCHEMA auth TO budgetuser;
GRANT USAGE ON SCHEMA budget TO budgetuser;

-- Grant all privileges on all tables in the schemas
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA auth TO budgetuser;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA budget TO budgetuser;

-- Grant privileges on future tables
ALTER DEFAULT PRIVILEGES IN SCHEMA auth GRANT ALL ON TABLES TO budgetuser;
ALTER DEFAULT PRIVILEGES IN SCHEMA budget GRANT ALL ON TABLES TO budgetuser;

-- Set search path to include both schemas
ALTER DATABASE help_me_budget SET search_path TO auth, budget, public;
