# Database Setup

This directory contains the database configuration, migrations, and Docker setup for Help-Me-Budget.

## Quick Start

1. **Ensure correct permissions:**
   ```bash
   chmod 755 init
   chmod 644 init/01_create_schemas.sql
   ```

2. **Start PostgreSQL and Redis:**
   ```bash
   make up
   ```

3. **Run migrations:**
   ```bash
   make migrate-up
   ```

4. **Verify setup:**
   - PostgreSQL: `psql postgres://budgetuser:budgetpass@localhost:5432/help_me_budget`
   - Redis: `redis-cli ping`

## Database Schema

### Auth Schema
- **users**: Core user records (email, name, avatar)
- **user_oauth_providers**: OAuth provider links (Google, Facebook, etc.)

### Budget Schema
- **budgets**: Budget plans/scenarios
- **categories**: Income/expense categories with hierarchy
- **budget_entries**: Planned recurring transactions with frequency rules
- **accounts**: Bank accounts, credit cards, etc.
- **transactions**: Actual transactions with matching support

## Available Commands

```bash
make up              # Start PostgreSQL and Redis containers
make down            # Stop containers
make migrate-up      # Run pending migrations
make migrate-down    # Rollback last migration
make migrate-create NAME=<name>  # Create new migration
```

## Connection Details

**PostgreSQL:**
- Host: `localhost:5432`
- Database: `help_me_budget`
- User: `budgetuser`
- Password: `budgetpass`
- Schemas: `auth`, `budget`

**Redis:**
- Host: `localhost:6379`
- No password (local development)

## Creating Migrations

```bash
make migrate-create NAME=add_new_table
```

This creates two files:
- `migrations/NNNNNN_add_new_table.up.sql` - Apply migration
- `migrations/NNNNNN_add_new_table.down.sql` - Rollback migration

## Troubleshooting

**"Permission denied" on init directory:**
```bash
chmod 755 init
chmod 644 init/01_create_schemas.sql
make down && make up
```

**"Connection refused" when running migrations:**
- Ensure containers are running: `docker ps`
- Check database logs: `docker logs help-me-budget-db`
- Wait a few seconds for PostgreSQL to fully initialize
- Verify connection: `psql postgres://budgetuser:budgetpass@localhost:5432/help_me_budget`

**Migrations fail with "database not found":**
- Check that init script ran: `docker logs help-me-budget-db | grep schema`
- If not, fix permissions and restart containers

## Notes

- Docker volumes persist data between restarts
- Initial schemas are created via `init/01_create_schemas.sql`
- Search path is set to `auth, budget, public`
- All timestamps use `TIMESTAMP WITH TIME ZONE`
- Currency amounts use `NUMERIC(20, 2)`
