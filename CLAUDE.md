# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Help-Me-Budget is a budgeting application using a modern 3-tier architecture:
- **Backend**: Go + Fiber web framework (RESTful API)
- **Frontend**: SvelteKit (not yet scaffolded - directory empty)
- **Database**: PostgreSQL (not yet set up - directory empty)

The project is in early development with basic OAuth authentication (Google & Facebook) implemented.

## Architecture & Directory Structure

### `/api` - Backend (Go)
- **`cmd/server/main.go`**: Entry point for the API server
  - Fiber web server listening on port 3000
  - CORS middleware configured for localhost:5173
  - Initializes OAuth providers and sets up authentication routes
  - Minimal bootstrapping - delegates to auth module
- **`internal/auth/oauth.go`**: OAuth authentication module
  - `InitializeOAuthProviders()` - Loads environment variables and configures Google/Facebook OAuth
  - `SetupAuthRoutes()` - Registers OAuth flow routes
  - `handleAuthStart()` - Initiates OAuth flow with state verification
  - `handleAuthCallback()` - Completes OAuth flow and returns user info
  - `NewSessionStore()` - Creates Fiber session store with secure defaults
- **`go.mod` / `go.sum`**: Go dependency management
- **`.env`**: Environment variables (OAuth keys and secrets)

### `/frontend` - Frontend (SvelteKit)
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

### `/database` - PostgreSQL
- Currently empty - schema and migrations to be added

## Development Environment Setup

### Required Environment Variables
Set these in `api/.env`:
- `GOOGLE_KEY` - OAuth application ID from Google Console
- `GOOGLE_SECRET` - OAuth application secret from Google Console
- `FACEBOOK_KEY` - OAuth application ID from Facebook Developers
- `FACEBOOK_SECRET` - OAuth application secret from Facebook Developers

**Note**: The `.env` file currently contains real credentials. Ensure it's in `.gitignore` before committing any credentials.

### Local Development Prerequisites
- Go 1.25.4 or later
- Node.js/npm (for frontend development)
- PostgreSQL (when database development begins)

## Common Development Commands

### Backend (API)

**Build the API server:**
```bash
cd api && go build -o server ./cmd/server
```

**Run the API server (requires .env in api/ directory):**
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
   - Go API completes OAuth exchange and returns user data
   - Stores user data in secure HTTP-only cookie
   - Redirects user back to home page

6. **Check Authentication**:
   - Frontend calls `/api/auth/me` to check if user is logged in
   - SvelteKit backend returns user data from cookie
   - Frontend displays user email and provider

### Security Features

- Frontend never communicates directly with Go API (proxied through SvelteKit)
- User data stored in secure HTTP-only cookies (cannot be accessed by JavaScript)
- Session cookies are encrypted by Fiber session middleware
- CSRF protection via state parameter verification
- 24-hour cookie expiration

## Key Dependencies

**Backend (Go)**:
- `gofiber/fiber` - Web framework with built-in session management
- `markbates/goth` - OAuth provider abstraction
- `joho/godotenv` - Environment variable loading

## Important Notes for Future Development

- **Database**: Postgres schema and Go database integration needed for user persistence
- **Token Security**: Currently returning access tokens in responses - should be removed in production
- **Callback URLs**: OAuth callback URLs can be configured via environment variables:
  - `GOOGLE_CALLBACK_URL` - defaults to `http://localhost:5173/api/auth/callback/google`
  - `FACEBOOK_CALLBACK_URL` - defaults to `http://localhost:5173/api/auth/callback/facebook`
- **User Storage**: No user persistence yet - OAuth flow returns data but doesn't store it in database
- **Error Handling**: Consider standardizing error responses to JSON format
- **Production Deployment**:
  - Change `secure: false` to `secure: true` in cookie settings (requires HTTPS)
  - Set `sameSite: 'strict'` instead of `'lax'` for stricter CSRF protection
  - Update CORS origins to production domain
  - Store sensitive data like OAuth secrets in secure environment variables, not `.env` files

## Git Workflow

The repository uses a simple main-branch workflow. Review commit history:
```bash
git log --oneline
```

Recent commits show progression from project setup → tech stack documentation → OAuth implementation.
