# help-me-budget

## 🧱 Technology Stack

Help-Me-Budget is built using a modern, reliable, and high-performance technology stack chosen to balance developer productivity, long-term maintainability, and great user experience.

## 🗄️ Backend API — Go + Fiber

The backend is implemented in Go, chosen for its simplicity, speed, and robust standard library. Go’s concurrency model and minimal runtime make it an ideal choice for building highly reliable web services.

The API layer is built with the Fiber web framework, providing:
- Fast, lightweight HTTP routing
- Simple middleware patterns
- An Express.js-like developer experience
- Great performance under load

The backend exposes a clean set of RESTful endpoints that handle:
- Budget categories and accounts
- Transactions and recurring expenses
- User profiles and authentication
- Data aggregation for dashboards and analytics

Go + Fiber keeps the backend easy to reason about while remaining scalable for future features.

## 🗃️ Database — PostgreSQL + Redis

The application uses PostgreSQL as its main data store with Redis for session management.

**PostgreSQL** provides:
- Strong consistency guarantees for financial data
- JSONB support for flexible data modeling (transaction matching rules)
- Powerful indexing and query capabilities
- Separate schemas for auth and budget data

**Redis** handles:
- Session storage for OAuth flows
- Fast, ephemeral data caching
- High-performance key-value operations

This combination ensures data integrity for critical financial records while maintaining fast session management.

## 🎨 Frontend — SvelteKit

The frontend is built with SvelteKit, chosen for its simplicity, excellent performance, and intuitive development experience.

Key benefits:
- Fast, minimal runtime thanks to Svelte’s compiler-based approach
- Built-in routing and server-side rendering
- Easy integration with Go APIs
- Great DX with reactive components and clean syntax

SvelteKit allows the UI to stay highly responsive while keeping the codebase clean and maintainable.
It powers the app’s dashboards, charts, category management views, and the overall budgeting workflow.

## 🧩 Overall Architecture

Together, this stack offers:
- A fast, strongly typed backend
- A reliable, production-grade database
- A lightweight, reactive frontend
- A clear separation of concerns

The system is designed for long-term maintainability, easy feature expansion, and a smooth user experience.


## 🚀 Quick Start

### Prerequisites

- **Docker & Docker Compose** (for PostgreSQL and Redis)
- **Go 1.25.4+** (for backend API)
- **Node.js & npm** (for frontend)
- **golang-migrate** CLI (for database migrations)
  ```bash
  brew install golang-migrate
  # or
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  ```

### 1. Start Database Services

```bash
cd database
chmod 755 init && chmod 644 init/01_create_schemas.sql  # Ensure correct permissions
make up          # Starts PostgreSQL and Redis
make migrate-up  # Runs database migrations
```

PostgreSQL will be available at `localhost:5432`, Redis at `localhost:6379`.

### 2. Configure Environment Variables

Copy the example file and add your OAuth credentials:

```bash
cd api
cp .env.example .env
# Edit .env and add your OAuth keys
```

**OAuth Configuration (required):**
- `GOOGLE_KEY` - OAuth client ID from Google Console
- `GOOGLE_SECRET` - OAuth client secret from Google Console
- `GOOGLE_CALLBACK_URL` - OAuth callback URL (default: `http://localhost:5173/api/auth/callback/google`)
- `FACEBOOK_KEY` - Facebook App ID from Facebook Developers
- `FACEBOOK_SECRET` - Facebook App secret from Facebook Developers
- `FACEBOOK_CALLBACK_URL` - OAuth callback URL (default: `http://localhost:5173/api/auth/callback/facebook`)

**Database Configuration (defaults work with Docker setup):**

Option 1 - Use full connection string:
- `DATABASE_URL` - Full PostgreSQL connection string (e.g., `postgres://budgetuser:budgetpass@localhost:5432/help_me_budget?sslmode=disable`)

Option 2 - Use individual variables:
- `DB_HOST` - PostgreSQL host (default: `localhost`)
- `DB_PORT` - PostgreSQL port (default: `5432`)
- `DB_USER` - PostgreSQL username (default: `budgetuser`)
- `DB_PASSWORD` - PostgreSQL password (default: `budgetpass`)
- `DB_NAME` - Database name (default: `help_me_budget`)
- `DB_SSLMODE` - SSL mode (default: `disable`)

**Redis Configuration (defaults work with Docker setup):**

Option 1 - Use full connection string:
- `REDIS_URL` - Full Redis connection string (e.g., `redis://localhost:6379`)

Option 2 - Use individual variables:
- `REDIS_ADDR` - Redis address (default: `localhost:6379`)
- `REDIS_PASSWORD` - Redis password (leave empty for local development)

**Application Environment:**
- `APP_ENV` - Environment mode: `development` or `production` (default: `development`)

### 3. Start the Backend

```bash
cd api
go run ./cmd/server
```

API server runs at `http://localhost:3000`

### 4. Start the Frontend

```bash
cd frontend/help-me-budget
npm install
npm run dev
```

Frontend runs at `http://localhost:5173`

### 5. Test OAuth Login

Visit `http://localhost:5173` and click "Login with Google" or "Login with Facebook" to test the authentication flow.

## 📊 Database Schema

The application uses separate PostgreSQL schemas for logical separation:

### Auth Schema (`auth`)
- **users** - User accounts (email, name, avatar)
- **user_oauth_providers** - OAuth provider links (supports multiple providers per user)

### Budget Schema (`budget`)
- **budgets** - Budget plans/scenarios (users can have multiple)
- **categories** - Income/expense categories with hierarchy support
- **budget_entries** - Planned recurring transactions
  - Frequencies: once-off, daily, weekly, fortnightly, monthly, annually
  - JSONB `matching_rules` for auto-matching imported transactions
- **accounts** - Bank accounts, credit cards, cash
- **transactions** - Actual transactions
  - Links to budget entries when matched
  - Confidence levels: manual, auto_high, auto_low, unmatched

See `database/README.md` for detailed schema documentation.

## 📚 Documentation

- **CLAUDE.md** - Comprehensive codebase documentation for AI assistants
- **database/README.md** - Database setup and migration guide
- **api/.env.example** - Environment variable template

## 🛠️ Development

For detailed development commands and workflow, see `CLAUDE.md`.