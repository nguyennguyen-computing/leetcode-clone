package repository

import (
	"database/sql"
	"strings"

	"leetcode-clone-backend/pkg/models"
	"github.com/lib/pq"
)

// userRepository implements UserRepository interface
type userRepository struct {
	db *sql.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

// Create creates a new user
func (r *userRepository) Create(user *models.User) (*models.User, error) {
	query := `
		INSERT INTO users (username, email, password_hash, is_admin)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, password_hash, is_admin, created_at, updated_at`

	var created models.User
	err := r.db.QueryRow(query, user.Username, user.Email, user.PasswordHash, user.IsAdmin).Scan(
		&created.ID,
		&created.Username,
		&created.Email,
		&created.PasswordHash,
		&created.IsAdmin,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "username") {
					return nil, NewRepositoryError("Create", ErrDuplicateKey, "username_exists")
				}
				if strings.Contains(pqErr.Detail, "email") {
					return nil, NewRepositoryError("Create", ErrDuplicateKey, "email_exists")
				}
				return nil, NewRepositoryError("Create", ErrDuplicateKey, "duplicate_key")
			}
		}
		return nil, NewRepositoryError("Create", err, "database_error")
	}

	return &created, nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(id int) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, is_admin, created_at, updated_at
		FROM users
		WHERE id = $1`

	var user models.User
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewRepositoryError("GetByID", ErrNotFound, "user_not_found")
		}
		return nil, NewRepositoryError("GetByID", err, "database_error")
	}

	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, is_admin, created_at, updated_at
		FROM users
		WHERE username = $1`

	var user models.User
	err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewRepositoryError("GetByUsername", ErrNotFound, "user_not_found")
		}
		return nil, NewRepositoryError("GetByUsername", err, "database_error")
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, is_admin, created_at, updated_at
		FROM users
		WHERE email = $1`

	var user models.User
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.IsAdmin,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewRepositoryError("GetByEmail", ErrNotFound, "user_not_found")
		}
		return nil, NewRepositoryError("GetByEmail", err, "database_error")
	}

	return &user, nil
}

// Update updates an existing user
func (r *userRepository) Update(user *models.User) (*models.User, error) {
	query := `
		UPDATE users
		SET username = $2, email = $3, password_hash = $4, is_admin = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, username, email, password_hash, is_admin, created_at, updated_at`

	var updated models.User
	err := r.db.QueryRow(query, user.ID, user.Username, user.Email, user.PasswordHash, user.IsAdmin).Scan(
		&updated.ID,
		&updated.Username,
		&updated.Email,
		&updated.PasswordHash,
		&updated.IsAdmin,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewRepositoryError("Update", ErrNotFound, "user_not_found")
		}
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "username") {
					return nil, NewRepositoryError("Update", ErrDuplicateKey, "username_exists")
				}
				if strings.Contains(pqErr.Detail, "email") {
					return nil, NewRepositoryError("Update", ErrDuplicateKey, "email_exists")
				}
				return nil, NewRepositoryError("Update", ErrDuplicateKey, "duplicate_key")
			}
		}
		return nil, NewRepositoryError("Update", err, "database_error")
	}

	return &updated, nil
}

// Delete deletes a user by ID
func (r *userRepository) Delete(id int) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return NewRepositoryError("Delete", err, "database_error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return NewRepositoryError("Delete", err, "database_error")
	}

	if rowsAffected == 0 {
		return NewRepositoryError("Delete", ErrNotFound, "user_not_found")
	}

	return nil
}

// List retrieves a list of users with pagination
func (r *userRepository) List(limit, offset int) ([]*models.User, error) {
	query := `
		SELECT id, username, email, password_hash, is_admin, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, NewRepositoryError("List", err, "database_error")
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.PasswordHash,
			&user.IsAdmin,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, NewRepositoryError("List", err, "scan_error")
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, NewRepositoryError("List", err, "rows_error")
	}

	return users, nil
}