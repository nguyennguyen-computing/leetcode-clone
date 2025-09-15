package database

import (
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestConnect(t *testing.T) {
	// Skip test if no database is available
	if os.Getenv("TEST_DB_HOST") == "" {
		t.Skip("Skipping database test - no TEST_DB_HOST environment variable set")
	}

	config := &Config{
		Host:     os.Getenv("TEST_DB_HOST"),
		Port:     getEnv("TEST_DB_PORT", "5432"),
		User:     getEnv("TEST_DB_USER", "postgres"),
		Password: getEnv("TEST_DB_PASSWORD", "password"),
		DBName:   getEnv("TEST_DB_NAME", "leetcode_clone_test"),
		SSLMode:  "disable",
	}

	db, err := Connect(config)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	// Test that we can ping the database
	if err := db.Ping(); err != nil {
		t.Fatalf("Failed to ping test database: %v", err)
	}
}

func TestRunMigrations(t *testing.T) {
	// Skip test if no database is available
	if os.Getenv("TEST_DB_HOST") == "" {
		t.Skip("Skipping migration test - no TEST_DB_HOST environment variable set")
	}

	config := &Config{
		Host:     os.Getenv("TEST_DB_HOST"),
		Port:     getEnv("TEST_DB_PORT", "5432"),
		User:     getEnv("TEST_DB_USER", "postgres"),
		Password: getEnv("TEST_DB_PASSWORD", "password"),
		DBName:   getEnv("TEST_DB_NAME", "leetcode_clone_test"),
		SSLMode:  "disable",
	}

	db, err := Connect(config)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	defer db.Close()

	// Clean up any existing schema_migrations table
	db.Exec("DROP TABLE IF EXISTS schema_migrations CASCADE")
	db.Exec("DROP TABLE IF EXISTS user_progress CASCADE")
	db.Exec("DROP TABLE IF EXISTS submissions CASCADE")
	db.Exec("DROP TABLE IF EXISTS test_cases CASCADE")
	db.Exec("DROP TABLE IF EXISTS problems CASCADE")
	db.Exec("DROP TABLE IF EXISTS users CASCADE")

	// Test running migrations
	if err := RunMigrations(db); err != nil {
		t.Fatalf("Failed to run migrations: %v", err)
	}

	// Verify that tables were created
	tables := []string{"users", "problems", "test_cases", "submissions", "user_progress", "schema_migrations"}
	for _, table := range tables {
		var exists bool
		err := db.QueryRow(`
			SELECT EXISTS (
				SELECT FROM information_schema.tables 
				WHERE table_schema = 'public' 
				AND table_name = $1
			)
		`, table).Scan(&exists)

		if err != nil {
			t.Fatalf("Failed to check if table %s exists: %v", table, err)
		}

		if !exists {
			t.Errorf("Table %s was not created", table)
		}
	}

	// Verify that migration was recorded
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM schema_migrations").Scan(&count)
	if err != nil {
		t.Fatalf("Failed to count migrations: %v", err)
	}

	if count == 0 {
		t.Error("No migrations were recorded in schema_migrations table")
	}
}
