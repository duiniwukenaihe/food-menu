# Database Migrations

This directory contains SQL migration files managed by [golang-migrate](https://github.com/golang-migrate/migrate).

## File Naming Convention

Migrations follow the naming pattern:
```
{version}_{description}.{direction}.sql
```

Example:
```
000001_initial_schema.up.sql
000001_initial_schema.down.sql
```

- **version**: Sequential version number (e.g., 000001, 000002)
- **description**: Brief description of the migration
- **direction**: Either `up` (apply) or `down` (rollback)

## Available Migrations

### 000001_initial_schema
Creates all base tables including:
- users
- categories (content categories)
- contents
- dish_categories
- dishes
- dish_pairings
- media
- recommendations

### 000002_seed_data
Seeds initial data:
- Default users (admin, chef_demo)
- Dish categories
- Content categories
- Sample dishes with full details
- Sample dish pairings

## Running Migrations

Migrations are automatically run when the application starts via `InitDB()`.

To manually manage migrations, you can use the golang-migrate CLI:

### Install CLI
```bash
# macOS
brew install golang-migrate

# Linux
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/

# Windows
scoop install migrate
```

### Apply Migrations
```bash
migrate -path ./migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" up
```

### Rollback Migrations
```bash
migrate -path ./migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" down 1
```

### Check Migration Version
```bash
migrate -path ./migrations -database "postgres://user:password@localhost:5432/dbname?sslmode=disable" version
```

## Creating New Migrations

1. Generate migration files:
```bash
timestamp=$(date +%Y%m%d%H%M%S)
touch migrations/${timestamp}_your_migration_name.up.sql
touch migrations/${timestamp}_your_migration_name.down.sql
```

2. Edit the `.up.sql` file with your schema changes
3. Edit the `.down.sql` file with the rollback logic
4. Test both directions before committing

## Best Practices

1. **Always create both up and down migrations** - Ensure rollback capability
2. **Test migrations locally first** - Verify both up and down work correctly
3. **Make migrations idempotent** - Use `IF EXISTS`, `IF NOT EXISTS` clauses
4. **Keep migrations focused** - One logical change per migration
5. **Never modify existing migrations** - Create new migrations for changes
6. **Use transactions** - Wrap complex changes in BEGIN/COMMIT
7. **Document complex migrations** - Add comments explaining non-obvious changes
8. **Backup before migrating** - Always backup production data before applying migrations

## Migration Examples

### Adding a Column
```sql
-- up
ALTER TABLE dishes ADD COLUMN calories INTEGER;

-- down
ALTER TABLE dishes DROP COLUMN IF EXISTS calories;
```

### Creating an Index
```sql
-- up
CREATE INDEX IF NOT EXISTS idx_dishes_calories ON dishes(calories);

-- down
DROP INDEX IF EXISTS idx_dishes_calories;
```

### Adding a Foreign Key
```sql
-- up
ALTER TABLE dishes 
ADD CONSTRAINT fk_dishes_category 
FOREIGN KEY (category_id) REFERENCES dish_categories(id);

-- down
ALTER TABLE dishes DROP CONSTRAINT IF EXISTS fk_dishes_category;
```

## Troubleshooting

### Dirty Database State
If a migration fails partway through, the database may be in a "dirty" state:

```bash
# Check current version
migrate -path ./migrations -database "$DATABASE_URL" version

# Force to a specific version (use with caution)
migrate -path ./migrations -database "$DATABASE_URL" force 1
```

### Failed Migration
1. Check the error message in logs
2. Review the migration SQL
3. Fix the issue
4. Roll back if needed: `migrate down 1`
5. Re-apply: `migrate up`

## Environment Variables

The application uses these environment variables for migrations:

- `DATABASE_URL` - PostgreSQL/MySQL connection string
- `DATABASE_DRIVER` - Database driver (postgres/mysql), auto-detected if not set
- `MIGRATIONS_PATH` - Path to migrations directory (defaults to ./migrations)

Example:
```bash
export DATABASE_URL="postgres://postgres:password@localhost:5432/appdb?sslmode=disable"
export DATABASE_DRIVER="postgres"
export MIGRATIONS_PATH="./migrations"
```

## Database Support

Both PostgreSQL and MySQL are supported. The migration system auto-detects the database type from the connection string.

### PostgreSQL
```
postgres://user:password@host:port/database?sslmode=disable
```

### MySQL
```
mysql://user:password@tcp(host:port)/database?parseTime=true
```

## Integration with Application

The application automatically runs migrations on startup via `database.InitDB()`:

```go
import "example.com/app/internal/database"

db, err := database.InitDB()  // Automatically runs migrations
```

To disable auto-migration, you'll need to modify `internal/database/database.go`.

## Schema Documentation

For detailed schema documentation, see `/docs/schema.md`.
