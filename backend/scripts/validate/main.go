package main

import (
	"database/sql"
	"fmt"
	"log"

	"leetcode-clone-backend/pkg/database"
)

func main() {
	// Load database configuration
	dbConfig := database.LoadConfigFromEnv()
	
	// Connect to database
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	fmt.Println("=== Database Schema Validation ===")
	
	// Check if all required tables exist
	tables := []string{"users", "problems", "test_cases", "submissions", "user_progress", "schema_migrations"}
	
	fmt.Println("\nChecking tables:")
	for _, table := range tables {
		exists, err := tableExists(db, table)
		if err != nil {
			log.Printf("Error checking table %s: %v", table, err)
			continue
		}
		
		status := "❌ MISSING"
		if exists {
			status = "✅ EXISTS"
		}
		fmt.Printf("  %s: %s\n", table, status)
	}
	
	// Check indexes
	fmt.Println("\nChecking key indexes:")
	indexes := []string{
		"idx_problems_difficulty",
		"idx_problems_tags", 
		"idx_submissions_user_id",
		"idx_submissions_problem_id",
		"idx_user_progress_user_id",
	}
	
	for _, index := range indexes {
		exists, err := indexExists(db, index)
		if err != nil {
			log.Printf("Error checking index %s: %v", index, err)
			continue
		}
		
		status := "❌ MISSING"
		if exists {
			status = "✅ EXISTS"
		}
		fmt.Printf("  %s: %s\n", index, status)
	}
	
	// Check applied migrations
	fmt.Println("\nApplied migrations:")
	migrations, err := getAppliedMigrations(db)
	if err != nil {
		log.Printf("Error getting migrations: %v", err)
	} else {
		for _, migration := range migrations {
			fmt.Printf("  ✅ Migration %d: %s\n", migration.Version, migration.Name)
		}
	}
	
	fmt.Println("\n=== Validation Complete ===")
}

func tableExists(db *sql.DB, tableName string) (bool, error) {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_schema = 'public' 
			AND table_name = $1
		)
	`, tableName).Scan(&exists)
	return exists, err
}

func indexExists(db *sql.DB, indexName string) (bool, error) {
	var exists bool
	err := db.QueryRow(`
		SELECT EXISTS (
			SELECT FROM pg_indexes 
			WHERE schemaname = 'public' 
			AND indexname = $1
		)
	`, indexName).Scan(&exists)
	return exists, err
}

type Migration struct {
	Version int
	Name    string
}

func getAppliedMigrations(db *sql.DB) ([]Migration, error) {
	rows, err := db.Query("SELECT version, name FROM schema_migrations ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var migrations []Migration
	for rows.Next() {
		var m Migration
		if err := rows.Scan(&m.Version, &m.Name); err != nil {
			return nil, err
		}
		migrations = append(migrations, m)
	}
	
	return migrations, nil
}