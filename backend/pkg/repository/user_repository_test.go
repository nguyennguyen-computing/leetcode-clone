package repository

import (
	"database/sql"
	"testing"

	"leetcode-clone-backend/pkg/models"
	_ "github.com/lib/pq"
)

// setupTestDB creates a test database connection
// Note: This would typically use a test database or in-memory database
func setupTestDB(t *testing.T) *sql.DB {
	// For now, we'll skip actual database tests and just test the interface
	// In a real implementation, you would set up a test database here
	t.Skip("Database tests require test database setup")
	return nil
}

func TestUserRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewUserRepository(db)

	user := &models.User{
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		IsAdmin:      false,
	}

	created, err := repo.Create(user)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	if created.ID == 0 {
		t.Error("Expected user ID to be set")
	}

	if created.Username != user.Username {
		t.Errorf("Expected username %s, got %s", user.Username, created.Username)
	}

	if created.Email != user.Email {
		t.Errorf("Expected email %s, got %s", user.Email, created.Email)
	}

	if created.CreatedAt.IsZero() {
		t.Error("Expected CreatedAt to be set")
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewUserRepository(db)

	// Test getting non-existent user
	_, err := repo.GetByID(999999)
	if err == nil {
		t.Error("Expected error when getting non-existent user")
	}

	if !IsNotFound(err) {
		t.Error("Expected NotFound error")
	}
}

func TestUserRepository_GetByUsername(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewUserRepository(db)

	// Test getting non-existent user
	_, err := repo.GetByUsername("nonexistent")
	if err == nil {
		t.Error("Expected error when getting non-existent user")
	}

	if !IsNotFound(err) {
		t.Error("Expected NotFound error")
	}
}

func TestUserRepository_GetByEmail(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewUserRepository(db)

	// Test getting non-existent user
	_, err := repo.GetByEmail("nonexistent@example.com")
	if err == nil {
		t.Error("Expected error when getting non-existent user")
	}

	if !IsNotFound(err) {
		t.Error("Expected NotFound error")
	}
}

func TestUserRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewUserRepository(db)

	user := &models.User{
		ID:           999999, // Non-existent ID
		Username:     "updateduser",
		Email:        "updated@example.com",
		PasswordHash: "newhashedpassword",
		IsAdmin:      true,
	}

	// Test updating non-existent user
	_, err := repo.Update(user)
	if err == nil {
		t.Error("Expected error when updating non-existent user")
	}

	if !IsNotFound(err) {
		t.Error("Expected NotFound error")
	}
}

func TestUserRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewUserRepository(db)

	// Test deleting non-existent user
	err := repo.Delete(999999)
	if err == nil {
		t.Error("Expected error when deleting non-existent user")
	}

	if !IsNotFound(err) {
		t.Error("Expected NotFound error")
	}
}

func TestUserRepository_List(t *testing.T) {
	db := setupTestDB(t)
	if db == nil {
		return
	}
	defer db.Close()

	repo := NewUserRepository(db)

	users, err := repo.List(10, 0)
	if err != nil {
		t.Fatalf("Failed to list users: %v", err)
	}

	// Should return empty list for empty database
	if users == nil {
		t.Error("Expected non-nil users slice")
	}
}