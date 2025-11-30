# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Help-Me-Budget is a budgeting application using a modern 3-tier architecture:
- **Backend**: Go + Fiber web framework (RESTful API)
- **Frontend**: SvelteKit
- **Database**: PostgreSQL (separate schemas for auth and budget data)
- **Authentication**: Supabase (handles OAuth with Google & Facebook)
- **Cache Store**: Redis (for future admin dashboard metrics and caching)

The project is in early development with Supabase authentication and database persistence implemented.

**Security**: The Go API uses shared secret authentication to ensure only the SvelteKit backend can make requests. All API calls require a valid `X-API-Key` header (except the health check endpoint).

## Architecture & Directory Structure

### `/api` - Backend (Go)
- **`cmd/server/main.go`**: Entry point for the API server
  - Fiber web server listening on port 3000
  - CORS middleware configured for localhost:5173
  - API key authentication middleware (validates X-API-Key header)
  - Initializes database and Redis
  - Sets up authentication and admin routes
- **`internal/middleware/api_auth.go`**: API authentication middleware
  - `ValidateAPIKey()` - Ensures requests come from SvelteKit backend
  - Whitelists health check endpoint (GET /) for monitoring
  - Returns 401 Unauthorized for missing/invalid API keys
- **`internal/auth/oauth.go`**: Authentication route setup
  - `SetupAuthRoutes()` - Registers authentication routes
  - Note: OAuth is handled by Supabase, not by this API
- **`internal/auth/supabase_sync_handler.go`**: Supabase user sync
  - `HandleSupabaseUserSync()` - Syncs Supabase authenticated users to local PostgreSQL
  - Creates or updates user records with OAuth provider information
- **`internal/auth/user_repository.go`**: User database operations
  - `UpsertUserWithOAuth()` - Creates or updates user and links OAuth provider
  - `GetUserByID()` - Retrieves user by ID
  - `GetUserByEmail()` - Retrieves user by email
- **`internal/database/database.go`**: PostgreSQL connection management
  - `InitDatabase()` - Initializes connection pool with pgx
  - `Close()` - Closes database connections
- **`internal/database/redis.go`**: Redis connection management
  - `InitRedis()` - Initializes Redis client
  - `CloseRedis()` - Closes Redis connection
  - Note: Redis is used for caching and future admin dashboard metrics
- **`internal/admin/handlers.go`**: Admin dashboard API endpoints
  - User management and analytics
  - System health monitoring
  - Audit logs and security features
- **`go.mod` / `go.sum`**: Go dependency management
- **`.env`**: Environment variables (Supabase keys, database, and Redis configuration)

### `/frontend` - Frontend (SvelteKit)
- **`src/lib/server/api-client.ts`**: Authenticated fetch helpers
  - `authenticatedFetch()` - Wrapper that automatically includes API key header
  - `authenticatedFetchWithUser()` - Includes both API key and user ID headers
  - Centralizes API URL configuration
- **`src/lib/server/auth-helpers.ts`**: Authentication helper functions
  - `getLocalUserId()` - Bridges Supabase auth with local PostgreSQL user ID
  - Syncs Supabase users to local database on first access
- **`src/routes/+page.svelte`**: Main landing page
  - Supabase authentication UI components
  - Shows authenticated user info when logged in
  - Client-side state management with reactive variables
  - Tailwind CSS styling with responsive design
- **`package.json`**: Node.js dependencies
- **`tsconfig.json`**: TypeScript configuration
- **`svelte.config.js`**: SvelteKit configuration with Tailwind CSS

### `/database` - PostgreSQL & Redis
- **`docker-compose.yml`**: Docker configuration for PostgreSQL and Redis
  - PostgreSQL 16 with persistent storage
  - Redis 7 with AOF persistence
  - Health checks for both services
- **`init/01_create_schemas.sql`**: Initial schema setup
  - Creates `auth` schema for user and OAuth data
  - Creates `budget` schema for budget planning and transactions
  - Sets up permissions and search paths
- **`migrations/`**: Database migrations using golang-migrate
  - `000001_create_users_table` - Users and OAuth providers tables
  - `000002_create_budget_schema` - Budgets, categories, accounts, transactions, and budget entries
- **`Makefile`**: Database management commands
  - `make up` - Start PostgreSQL and Redis containers
  - `make down` - Stop containers
  - `make migrate-up` - Run pending migrations
  - `make migrate-down` - Rollback last migration
  - `make migrate-create NAME=<name>` - Create new migration files

#### Database Schema Design

**Auth Schema (`auth`):**
- `users` - Core user records (email, name, avatar, timestamps)
- `user_oauth_providers` - OAuth provider links (supports multiple providers per user)

**Budget Schema (`budget`):**
- `budgets` - Budget plans (users can have multiple budgets/scenarios)
- `categories` - Income/expense categories with hierarchy support
- `budget_entries` - Planned recurring income/expenses with frequency rules
  - Supports: once-off, daily, weekly, fortnightly, monthly, annually
  - `matching_rules` JSONB field for auto-matching imported transactions
- `accounts` - Bank accounts, credit cards, cash accounts
- `transactions` - Actual transactions (can be linked to budget entries)
  - `match_confidence` field: manual, auto_high, auto_low, unmatched

## Development Environment Setup

### Required Environment Variables
Set these in `api/.env` and `frontend/help-me-budget/.env` (see example files for templates):

**Supabase Configuration (Frontend):**
- `PUBLIC_SUPABASE_URL` - Your Supabase project URL
- `PUBLIC_SUPABASE_ANON_KEY` - Your Supabase anonymous/public API key

**Database Configuration (choose one):**
- `DATABASE_URL` - Full PostgreSQL connection string
  - OR individual variables: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`

**Redis Configuration (Backend - choose one):**
- `REDIS_URL` - Full Redis connection string
  - OR individual variables: `REDIS_ADDR`, `REDIS_PASSWORD`
- Note: Redis is used for caching and admin dashboard metrics, not authentication

**Application:**
- `APP_ENV` - Environment (development/production)

**API Authentication:**
- `API_SECRET_KEY` - Shared secret between SvelteKit and Go API (required for both frontend and backend)
  - Generate with: `node -e "console.log(require('crypto').randomBytes(32).toString('hex'))"`
  - Must be identical in both `api/.env` and `frontend/help-me-budget/.env`

**Note**: The `.env` files contain sensitive credentials. They're already in `.gitignore`.

### Local Development Prerequisites
- Go 1.25.4 or later
- Docker & Docker Compose (for PostgreSQL and Redis)
- Node.js/npm (for frontend development)
- golang-migrate CLI (for running migrations): `brew install golang-migrate` or `go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest`

## Common Development Commands

### Database (Docker + Migrations)

**Start PostgreSQL and Redis:**
```bash
cd database && make up
```
PostgreSQL runs on `localhost:5432`, Redis on `localhost:6379`

**Stop containers:**
```bash
cd database && make down
```

**Run migrations:**
```bash
cd database && make migrate-up
```

**Rollback last migration:**
```bash
cd database && make migrate-down
```

**Create new migration:**
```bash
cd database && make migrate-create NAME=add_new_feature
```

### Backend (API)

**Build the API server:**
```bash
cd api && go build -o server ./cmd/server
```

**Run the API server (requires database and Redis to be running):**
```bash
cd api && go run ./cmd/server
```
Server runs on `http://localhost:3000`

**Run tests (when tests are added):**
```bash
cd api && go test ./...
```

**Format Go code:**
```bash
cd api && go fmt ./...
```

**Check Go code for issues:**
```bash
cd api && go vet ./...
```

**View Go module dependencies:**
```bash
cd api && go list -m all
```

### Frontend (SvelteKit)

**Run the development server:**
```bash
cd frontend/help-me-budget && npm run dev
```
Frontend runs on `http://localhost:5173`

**Build for production:**
```bash
cd frontend/help-me-budget && npm run build
```

**Format and lint code:**
```bash
cd frontend/help-me-budget && npm run format
cd frontend/help-me-budget && npm run lint
```

**Run unit tests:**
```bash
cd frontend/help-me-budget && npm run test:unit
```

**Run end-to-end tests:**
```bash
cd frontend/help-me-budget && npm run test:e2e
```

## Authentication Flow Overview

The application uses Supabase for authentication with local PostgreSQL user sync:

**Frontend (Supabase Client) ↔ Supabase Auth ↔ SvelteKit Backend ↔ Go API ↔ PostgreSQL**

### Flow Steps

1. **Frontend Authentication**:
   - User authenticates via Supabase client (Google/Facebook OAuth)
   - Supabase handles all OAuth flows and session management
   - Frontend receives Supabase session with JWT token

2. **User Sync to Local Database**:
   - On first access, SvelteKit backend calls `getLocalUserId()` helper
   - Helper extracts user info from Supabase session
   - Calls Go API `/auth/sync` endpoint to sync user to local PostgreSQL
   - Returns local user ID for subsequent API calls

3. **Authenticated API Requests**:
   - SvelteKit backend makes API calls to Go backend
   - Includes `X-API-Key` (shared secret) and `X-User-ID` (local user ID) headers
   - Go API validates requests and performs operations

### Security Features

- **Supabase Authentication**: Industry-standard OAuth implementation with session management
- **API Key Authentication**: Shared secret (`X-API-Key` header) ensures only SvelteKit backend can access Go API
  - All endpoints require valid API key (except GET / health check)
  - Returns 401 Unauthorized for missing/invalid keys
  - Prevents direct API access from unauthorized clients
- **User Context Middleware**: Extracts user ID from headers for authorization
- Frontend never communicates directly with Go API (proxied through SvelteKit)
- User records persisted to PostgreSQL with OAuth provider links
- Supports multiple OAuth providers per user account

## Key Dependencies

**Backend (Go)**:
- `gofiber/fiber/v2` - Web framework with middleware support
- `jackc/pgx/v5` - PostgreSQL driver and connection pooling
- `redis/go-redis/v9` - Redis client for caching and metrics
- `google/uuid` - UUID generation for unique identifiers
- `joho/godotenv` - Environment variable loading

**Frontend (SvelteKit)**:
- `@supabase/supabase-js` - Supabase client for authentication
- SvelteKit - Full-stack web framework
- Tailwind CSS - Utility-first CSS framework

**Database Tools**:
- `golang-migrate` - Database migration management
- PostgreSQL 16 - Database server
- Redis 7 - Cache store for metrics and future features

## Important Notes for Future Development

- **Authentication**: Supabase handles all OAuth flows. Local PostgreSQL syncs user records for relational data.
- **Redis Usage**: Currently initialized but not heavily used. Reserved for:
  - Admin dashboard metrics caching
  - Performance optimization for frequently accessed data
  - Future real-time features
- **Transaction Matching**: The `matching_rules` JSONB field in `budget_entries` enables auto-matching of imported transactions to planned entries
- **Budget Projections**: Use `budget_entries` with frequency rules to project future cash flow
- **Error Handling**: Consider standardizing error responses to JSON format across all endpoints
- **Production Deployment**:
  - Update CORS origins to production domain
  - Use managed PostgreSQL and Redis services instead of Docker
  - Store sensitive data in secure secret management (AWS Secrets Manager, HashiCorp Vault, etc.)
  - Enable SSL for PostgreSQL connections
  - Configure Supabase production environment
- **Future Features**:
  - Bank account import/sync integration
  - Transaction matching algorithm implementation
  - Budget vs actual comparison dashboards
  - Recurring transaction automation
  - Admin dashboard with user analytics and system health monitoring

## Git Workflow

The repository uses a simple main-branch workflow. Review commit history:
```bash
git log --oneline
```

Recent commits show progression from project setup → tech stack documentation → OAuth implementation.
