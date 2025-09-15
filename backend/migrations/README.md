# Database Migrations

This directory contains SQL migration files for the LeetCode clone database schema.

## Migration System

The migration system automatically runs when the application starts and applies any pending migrations to the database.

### Migration Files

Migration files follow the naming convention: `{version}_{description}.sql`

- `version`: A zero-padded number (e.g., 001, 002, 003)
- `description`: A descriptive name using underscores

Example: `001_initial_schema.sql`

### How It Works

1. On application startup, the migration system:
   - Creates a `schema_migrations` table to track applied migrations
   - Scans the `migrations/` directory for `.sql` files
   - Applies any migrations that haven't been run yet
   - Records each successful migration in the tracking table

2. Migrations are applied in order based on their version number
3. Each migration runs in a transaction - if it fails, it's rolled back
4. Once a migration is applied, it won't be run again

### Current Schema

The initial migration (`001_initial_schema.sql`) creates:

#### Tables
- `users` - User accounts and authentication
- `problems` - Coding problems with descriptions and metadata
- `test_cases` - Input/output test cases for problems
- `submissions` - User code submissions and results
- `user_progress` - Tracking of user problem-solving progress
- `schema_migrations` - Migration tracking (created automatically)

#### Indexes
- Performance indexes on frequently queried columns
- GIN index on problem tags for efficient tag-based filtering
- Composite indexes for common query patterns

#### Features
- Foreign key constraints for data integrity
- Check constraints for data validation
- Automatic timestamp updates with triggers
- Proper CASCADE deletion rules

### Adding New Migrations

To add a new migration:

1. Create a new file with the next version number: `002_add_new_feature.sql`
2. Write your SQL DDL statements
3. Test the migration on a development database
4. The migration will be applied automatically on next application start

### Environment Variables

The database connection uses these environment variables:

- `DB_HOST` - Database host (default: localhost)
- `DB_PORT` - Database port (default: 5432)
- `DB_USER` - Database user (default: postgres)
- `DB_PASSWORD` - Database password (default: password)
- `DB_NAME` - Database name (default: leetcode_clone)
- `DB_SSLMODE` - SSL mode (default: disable)

### Validation

Use the validation script to check schema status:

```bash
go run scripts/validate_schema.go
```

This will verify that all tables and indexes are properly created.