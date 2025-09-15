package migrations

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// Migration represents a single database migration
type Migration struct {
	Version int
	Name    string
	SQL     string
}

// Migrator handles database migrations
type Migrator struct {
	db            *sql.DB
	migrationsDir string
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB, migrationsDir string) *Migrator {
	return &Migrator{
		db:            db,
		migrationsDir: migrationsDir,
	}
}

// CreateMigrationsTable creates the migrations tracking table
func (m *Migrator) CreateMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := m.db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}
	return nil
}

// GetAppliedMigrations returns a list of applied migration versions
func (m *Migrator) GetAppliedMigrations() (map[int]bool, error) {
	applied := make(map[int]bool)
	
	rows, err := m.db.Query("SELECT version FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}
	defer rows.Close()
	
	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return nil, fmt.Errorf("failed to scan migration version: %w", err)
		}
		applied[version] = true
	}
	
	return applied, nil
}

// LoadMigrations loads all migration files from the migrations directory
func (m *Migrator) LoadMigrations() ([]Migration, error) {
	files, err := ioutil.ReadDir(m.migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}
	
	var migrations []Migration
	
	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".sql") {
			continue
		}
		
		// Extract version from filename (e.g., "001_initial_schema.sql" -> 1)
		parts := strings.Split(file.Name(), "_")
		if len(parts) < 2 {
			log.Printf("Skipping migration file with invalid name format: %s", file.Name())
			continue
		}
		
		var version int
		if _, err := fmt.Sscanf(parts[0], "%d", &version); err != nil {
			log.Printf("Skipping migration file with invalid version: %s", file.Name())
			continue
		}
		
		// Read migration content
		content, err := ioutil.ReadFile(filepath.Join(m.migrationsDir, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}
		
		// Extract name from filename (remove version and extension)
		name := strings.TrimSuffix(strings.Join(parts[1:], "_"), ".sql")
		
		migrations = append(migrations, Migration{
			Version: version,
			Name:    name,
			SQL:     string(content),
		})
	}
	
	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})
	
	return migrations, nil
}

// ApplyMigration applies a single migration
func (m *Migrator) ApplyMigration(migration Migration) error {
	// Start transaction
	tx, err := m.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}
	defer tx.Rollback()
	
	// Execute migration SQL
	if _, err := tx.Exec(migration.SQL); err != nil {
		return fmt.Errorf("failed to execute migration %d (%s): %w", migration.Version, migration.Name, err)
	}
	
	// Record migration as applied
	if _, err := tx.Exec(
		"INSERT INTO schema_migrations (version, name) VALUES ($1, $2)",
		migration.Version, migration.Name,
	); err != nil {
		return fmt.Errorf("failed to record migration %d: %w", migration.Version, err)
	}
	
	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration %d: %w", migration.Version, err)
	}
	
	log.Printf("Applied migration %d: %s", migration.Version, migration.Name)
	return nil
}

// Migrate runs all pending migrations
func (m *Migrator) Migrate() error {
	// Create migrations table if it doesn't exist
	if err := m.CreateMigrationsTable(); err != nil {
		return err
	}
	
	// Get applied migrations
	applied, err := m.GetAppliedMigrations()
	if err != nil {
		return err
	}
	
	// Load all migrations
	migrations, err := m.LoadMigrations()
	if err != nil {
		return err
	}
	
	// Apply pending migrations
	var appliedCount int
	for _, migration := range migrations {
		if applied[migration.Version] {
			log.Printf("Migration %d (%s) already applied, skipping", migration.Version, migration.Name)
			continue
		}
		
		if err := m.ApplyMigration(migration); err != nil {
			return err
		}
		appliedCount++
	}
	
	if appliedCount == 0 {
		log.Println("No pending migrations to apply")
	} else {
		log.Printf("Successfully applied %d migrations", appliedCount)
	}
	
	return nil
}