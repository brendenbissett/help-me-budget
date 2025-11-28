# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Help-Me-Budget is a budgeting application using a modern 3-tier architecture:
- **Backend**: Go + Fiber web framework (RESTful API)
- **Frontend**: SvelteKit
- **Database**: PostgreSQL (separate schemas for auth and budget data)
- **Session Store**: Redis

The project is in early development with OAuth authentication (Google & Facebook) and database persistence implemented.

**Security**: The Go API uses shared secret authentication to ensure only the SvelteKit backend can make requests. All API calls require a valid `X-API-Key` header (except the health check endpoint).

## Architecture & Directory Structure

### `/api` - Backend (Go)
- **`cmd/server/main.go`**: Entry point for the API server
  - Fiber web server listening on port 3000
  - CORS middleware configured for localhost:5173
  - API key authentication middleware (validates X-API-Key header)
  - Initializes database, Redis, and OAuth providers
  - Sets up authentication routes
- **`internal/middleware/api_auth.go`**: API authentication middleware
  - `ValidateAPIKey()` - Ensures requests come from SvelteKit backend
  - Whitelists health check endpoint (GET /) for monitoring
  - Returns 401 Unauthorized for missing/invalid API keys
- **`internal/auth/oauth.go`**: OAuth authentication module
  - `InitializeOAuthProviders()` - Loads environment variables and configures Google/Facebook OAuth
  - `SetupAuthRoutes()` - Registers OAuth flow routes
  - `handleAuthStart()` - Initiates OAuth flow with state verification
  - `handleAuthCallback()` - Completes OAuth flow, persists user to database, and returns user info
  - `NewSessionStore()` - Creates Fiber session store with Redis backend
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
- **`go.mod` / `go.sum`**: Go dependency management
- **`.env`**: Environment variables (OAuth keys, database, and Redis configuration)

### `/frontend` - Frontend (SvelteKit)
- **`src/lib/server/api-client.ts`**: Authenticated fetch helpers
  - `authenticatedFetch()` - Wrapper that automatically includes API key header
  - `authenticatedFetchWithUser()` - Includes both API key and user ID headers
  - Centralizes API URL configuration
- **`src/lib/server/auth-helpers.ts`**: Authentication helper functions
  - `getLocalUserId()` - Bridges Supabase auth with local PostgreSQL user ID
- **`src/routes/+page.svelte`**: Main landing page with OAuth login
  - Displays login buttons for Google and Facebook
  - Shows authenticated user info (email and provider) when logged in
  - Client-side state management with reactive variables
  - Tailwind CSS styling with responsive design
- **`src/routes/api/auth/login/[provider]/+server.ts`**: Login initiation endpoint
  - Proxies login requests to Go API at `http://localhost:3000`
  - Handles session cookie management from OAuth provider
  - Returns auth URL for frontend redirect
- **`src/routes/api/auth/callback/[provider]/+server.ts`**: OAuth callback handler
  - Receives callback from OAuth provider
  - Communicates with Go API to complete auth flow
  - Stores user data in secure HTTP-only cookie
  - Redirects to home page on success
- **`src/routes/api/auth/me/+server.ts`**: Current user endpoint
  - Returns authenticated user info from stored cookie
  - Returns 401 if not authenticated
- **`src/routes/api/auth/logout/+server.ts`**: Logout endpoint
  - Clears user data and session cookies
  - Returns success response
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
Set these in `api/.env` (see `api/.env.example` for template):

**OAuth Configuration:**
- `GOOGLE_KEY` - OAuth application ID from Google Console
- `GOOGLE_SECRET` - OAuth application secret from Google Console
- `GOOGLE_CALLBACK_URL` - OAuth callback URL (default: `http://localhost:5173/api/auth/callback/google`)
- `FACEBOOK_KEY` - OAuth application ID from Facebook Developers
- `FACEBOOK_SECRET` - OAuth application secret from Facebook Developers
- `FACEBOOK_CALLBACK_URL` - OAuth callback URL (default: `http://localhost:5173/api/auth/callback/facebook`)

**Database Configuration (choose one):**
- `DATABASE_URL` - Full PostgreSQL connection string
  - OR individual variables: `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`, `DB_SSLMODE`

**Redis Configuration (choose one):**
- `REDIS_URL` - Full Redis connection string
  - OR individual variables: `REDIS_ADDR`, `REDIS_PASSWORD`

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

## OAuth Flow Overview

The OAuth authentication uses a 3-layer architecture for security:

**Frontend → SvelteKit Backend → Go API → OAuth Provider**

### Flow Steps

1. **Frontend Login Initiation**:
   - User clicks "Login with Google/Facebook" button on landing page (`src/routes/+page.svelte`)
   - Frontend calls `/api/auth/login/{provider}` on SvelteKit backend

2. **SvelteKit Backend - Login Route** (`src/routes/api/auth/login/[provider]/+server.ts`):
   - Proxies request to Go API (`http://localhost:3000/auth/:provider`)
   - Captures session cookie from Go API response
   - Returns OAuth provider's auth URL to frontend

3. **Frontend Redirect**:
   - Frontend redirects user to OAuth provider's consent screen

4. **OAuth Provider Callback**:
   - Provider redirects back to SvelteKit backend at `/api/auth/callback/{provider}`

5. **SvelteKit Backend - Callback Route** (`src/routes/api/auth/callback/[provider]/+server.ts`):
   - Forwards callback to Go API with query parameters and session cookie
   - Go API completes OAuth exchange, persists user to PostgreSQL, and returns user data
   - Stores user data in secure HTTP-only cookie
   - Redirects user back to home page

6. **Check Authentication**:
   - Frontend calls `/api/auth/me` to check if user is logged in
   - SvelteKit backend returns user data from cookie
   - Frontend displays user email and provider

### Security Features

- **API Key Authentication**: Shared secret (`X-API-Key` header) ensures only SvelteKit backend can access Go API
  - All endpoints require valid API key (except GET / health check)
  - Returns 401 Unauthorized for missing/invalid keys
  - Prevents direct API access from unauthorized clients
- Frontend never communicates directly with Go API (proxied through SvelteKit)
- User data stored in secure HTTP-only cookies (cannot be accessed by JavaScript)
- Session data stored in Redis (not in-memory)
- OAuth sessions managed in Redis with 24-hour expiration
- User records persisted to PostgreSQL with OAuth provider links
- CSRF protection via state parameter verification
- Supports multiple OAuth providers per user account

## Key Dependencies

**Backend (Go)**:
- `gofiber/fiber/v2` - Web framework with middleware support
- `markbates/goth` - OAuth provider abstraction (Google, Facebook)
- `jackc/pgx/v5` - PostgreSQL driver and connection pooling
- `redis/go-redis/v9` - Redis client
- `gofiber/storage/redis/v3` - Redis storage adapter for Fiber sessions
- `joho/godotenv` - Environment variable loading

**Database Tools**:
- `golang-migrate` - Database migration management
- PostgreSQL 16 - Database server
- Redis 7 - Session and cache store

## Important Notes for Future Development

- **Transaction Matching**: The `matching_rules` JSONB field in `budget_entries` enables auto-matching of imported transactions to planned entries
- **Budget Projections**: Use `budget_entries` with frequency rules to project future cash flow
- **OAuth Tokens**: Access/refresh tokens are stored in `user_oauth_providers` but not currently used for API calls
- **Error Handling**: Consider standardizing error responses to JSON format across all endpoints
- **Production Deployment**:
  - Change `secure: false` to `secure: true` in cookie settings (requires HTTPS)
  - Set `sameSite: 'strict'` instead of `'lax'` for stricter CSRF protection
  - Update CORS origins to production domain
  - Use managed PostgreSQL and Redis services instead of Docker
  - Store sensitive data in secure secret management (AWS Secrets Manager, HashiCorp Vault, etc.)
  - Enable SSL for PostgreSQL connections
- **Future Features**:
  - Bank account import/sync integration
  - Transaction matching algorithm implementation
  - Budget vs actual comparison dashboards
  - Recurring transaction automation

## Git Workflow

The repository uses a simple main-branch workflow. Review commit history:
```bash
git log --oneline
```

Recent commits show progression from project setup → tech stack documentation → OAuth implementation.
