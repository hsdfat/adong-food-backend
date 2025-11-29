# Database Auto-Migration

This package handles automatic database schema initialization when connecting to a new database.

## How It Works

1. **Automatic Detection**: On application startup, the migration system checks if the database is already initialized by looking for the `master_users` table.

2. **Schema Creation**: If the database is empty, it automatically runs the following migrations in order:
   - `db.sql` - Creates all tables, constraints, and indexes
   - `init_admin_user.sql` - Creates the default admin user

3. **Embedded Files**: SQL migration files are embedded into the compiled binary using Go's `embed` package, so no external files are needed at runtime.

## Migration Files

All SQL migration files are stored in the `migrate/sql/` directory and are automatically embedded into the binary during compilation.

Current migrations:
- `db.sql` - Complete database schema
- `init_admin_user.sql` - Default admin user (username: admin, password: admin@adong)

## Usage

The migration runs automatically when the application starts. No manual intervention is required.

```go
// In cmd/main.go
if err := migrate.AutoMigrate(store.DB.GormClient); err != nil {
    log.Fatal("Failed to auto-migrate database:", err)
}
```

## Adding New Migrations

To add new migration files:

1. Create a new `.sql` file in the `migrate/sql/` directory
2. The files will be automatically discovered and embedded
3. Use the `RunMigrations()` function instead of `AutoMigrate()` for more flexible migration handling

## Configuration

The migration system uses the same database connection configured in `cmd/main.go`:
- Environment variable: `DATABASE_URL`
- Default: `host=localhost user=adong password=adong123 dbname=adongfood port=5432 sslmode=disable`

## Features

- **Idempotent**: Safe to run multiple times - skips if database is already initialized
- **Embedded**: SQL files are compiled into the binary (no external dependencies)
- **Logging**: Detailed logs for each migration step
- **Error Handling**: Proper error reporting and rollback on failure
- **Health Checks**: Built-in database health verification

## Docker Support

The Dockerfile has been updated to include the migration files during the build process:

```dockerfile
# Copy entire project for build (needed for embedded SQL files)
COPY . .

# Build the application
RUN go build -o /tmp/main cmd/main.go
```

The SQL files are embedded in the binary, so they're available at runtime even in the minimal scratch container.
