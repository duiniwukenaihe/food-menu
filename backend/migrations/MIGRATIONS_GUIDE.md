# Migrations Guide

## Quick Start

The database migrations are automatically run when you start the application. No manual intervention is needed for normal operation.

## What Happens on Startup

1. Application connects to the database using `DATABASE_URL` environment variable
2. Migration system automatically detects database type (PostgreSQL/MySQL)
3. Runs all pending migrations in order (000001, 000002, etc.)
4. Seed data is applied through migration 000002
5. Application continues with normal startup

## Environment Variables

```bash
# Required
DATABASE_URL="postgres://postgres:password@localhost:5432/appdb?sslmode=disable"

# Optional
DATABASE_DRIVER="postgres"  # Auto-detected if not set
MIGRATIONS_PATH="./migrations"  # Default location
```

## Database Support

### PostgreSQL (Recommended)
```bash
DATABASE_URL="postgres://postgres:password@localhost:5432/appdb?sslmode=disable"
```

### MySQL
```bash
DATABASE_URL="mysql://user:password@tcp(localhost:3306)/appdb?parseTime=true"
```

## Migration Files

Located in `/backend/migrations/`:

1. **000001_initial_schema** - All tables, indexes, constraints
2. **000002_seed_data** - Initial data (users, categories, sample dishes)

## Default Seed Data

### Users
- `admin` / admin@example.com (role: admin)
- `chef_demo` / chef@example.com (role: user)
- **Password**: `password123` (change in production!)

### Dish Categories
- Appetizers, Main Course, Desserts, Beverages, Breakfast, Salads

### Content Categories
- Recipes, Cooking Tips, Nutrition, Kitchen Tools

### Sample Dishes
5 fully-featured dishes with ingredients, steps, images, etc.

## Development Workflow

### First Time Setup
```bash
# Start database
docker-compose up -d postgres

# Start application (migrations run automatically)
cd backend && go run main.go
```

### Reset Database
```bash
# Using Make
make db-reset

# Manual
docker-compose exec postgres psql -U postgres -d appdb -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
cd backend && go run main.go
```

## Creating New Migrations

While the application auto-runs migrations, you can create new ones:

1. **Create migration files:**
```bash
timestamp=$(date +%Y%m%d%H%M%S)
touch backend/migrations/${timestamp}_your_change.up.sql
touch backend/migrations/${timestamp}_your_change.down.sql
```

2. **Write up migration** (e.g., add column):
```sql
-- 000003_add_dish_rating.up.sql
ALTER TABLE dishes ADD COLUMN rating DECIMAL(3,2) DEFAULT 0.0;
CREATE INDEX idx_dishes_rating ON dishes(rating);
```

3. **Write down migration** (rollback):
```sql
-- 000003_add_dish_rating.down.sql
DROP INDEX IF EXISTS idx_dishes_rating;
ALTER TABLE dishes DROP COLUMN IF EXISTS rating;
```

4. **Restart application** - Migration runs automatically

## Repository Pattern

The application uses repositories for database access:

```go
// Get repository store
store := repository.NewStore(db)

// Use repositories
ctx := context.Background()
user, err := store.Users().WithContext(ctx).FindByID(1)
dish, err := store.Dishes().WithContext(ctx).FindByID(1)
```

### Available Repositories
- **UserRepository** - User management
- **DishRepository** - Dishes/recipes with advanced filtering
- **ContentRepository** - CMS content

### Features
- Context-aware operations
- Transaction support
- Soft delete support
- Filter/search capabilities
- Preloaded relationships

## Schema Documentation

Complete schema documentation: `/docs/schema.md`

Includes:
- Table structures
- Relationships
- Indexes
- Constraints
- Query examples
- Best practices

## Troubleshooting

### Migration Fails
Check logs for specific error. Common issues:
- Database connection string incorrect
- Permissions issues
- Syntax errors in migration SQL

### Can't Connect to Database
```bash
# Verify database is running
docker-compose ps

# Check logs
docker-compose logs postgres

# Test connection
psql "postgres://postgres:password@localhost:5432/appdb"
```

### Start Fresh
```bash
make db-reset
```

This drops and recreates all tables, then re-applies migrations.

## Production Deployment

1. **Backup database** before deploying
2. Ensure `DATABASE_URL` is set correctly
3. Migrations run automatically on first startup
4. **Never modify existing migrations** - create new ones
5. Test migrations in staging first

## Best Practices

1. ✅ Let migrations run automatically
2. ✅ Create both up and down migrations
3. ✅ Test in development first
4. ✅ Use `IF EXISTS` / `IF NOT EXISTS` clauses
5. ✅ Keep migrations small and focused
6. ✅ Backup before production migrations
7. ❌ Never modify existing migration files
8. ❌ Don't skip migrations
9. ❌ Don't manually edit migrated databases

## Support

For issues or questions:
- Check `/docs/schema.md` for schema details
- Check `/backend/migrations/README.md` for advanced usage
- Review `/backend/internal/repository/examples.go` for code examples
