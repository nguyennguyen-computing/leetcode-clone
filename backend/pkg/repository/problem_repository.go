package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"leetcode-clone-backend/pkg/models"
	"github.com/lib/pq"
)

// problemRepository implements ProblemRepository interface
type problemRepository struct {
	db *sql.DB
}

// NewProblemRepository creates a new problem repository
func NewProblemRepository(db *sql.DB) ProblemRepository {
	return &problemRepository{db: db}
}

// Create creates a new problem
func (r *problemRepository) Create(problem *models.Problem) (*models.Problem, error) {
	query := `
		INSERT INTO problems (title, slug, description, difficulty, tags, examples, constraints, template_code)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, slug, description, difficulty, tags, examples, constraints, template_code, created_at, updated_at`

	var created models.Problem
	err := r.db.QueryRow(
		query,
		problem.Title,
		problem.Slug,
		problem.Description,
		problem.Difficulty,
		problem.Tags,
		problem.Examples,
		problem.Constraints,
		problem.TemplateCode,
	).Scan(
		&created.ID,
		&created.Title,
		&created.Slug,
		&created.Description,
		&created.Difficulty,
		&created.Tags,
		&created.Examples,
		&created.Constraints,
		&created.TemplateCode,
		&created.CreatedAt,
		&created.UpdatedAt,
	)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "slug") {
					return nil, NewRepositoryError("Create", ErrDuplicateKey, "slug_exists")
				}
				return nil, NewRepositoryError("Create", ErrDuplicateKey, "duplicate_key")
			case "23514": // check_violation
				return nil, NewRepositoryError("Create", ErrInvalidInput, "invalid_difficulty")
			}
		}
		return nil, NewRepositoryError("Create", err, "database_error")
	}

	return &created, nil
}

// GetByID retrieves a problem by ID
func (r *problemRepository) GetByID(id int) (*models.Problem, error) {
	query := `
		SELECT id, title, slug, description, difficulty, tags, examples, constraints, template_code, created_at, updated_at
		FROM problems
		WHERE id = $1`

	var problem models.Problem
	err := r.db.QueryRow(query, id).Scan(
		&problem.ID,
		&problem.Title,
		&problem.Slug,
		&problem.Description,
		&problem.Difficulty,
		&problem.Tags,
		&problem.Examples,
		&problem.Constraints,
		&problem.TemplateCode,
		&problem.CreatedAt,
		&problem.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewRepositoryError("GetByID", ErrNotFound, "problem_not_found")
		}
		return nil, NewRepositoryError("GetByID", err, "database_error")
	}

	return &problem, nil
}

// GetBySlug retrieves a problem by slug
func (r *problemRepository) GetBySlug(slug string) (*models.Problem, error) {
	query := `
		SELECT id, title, slug, description, difficulty, tags, examples, constraints, template_code, created_at, updated_at
		FROM problems
		WHERE slug = $1`

	var problem models.Problem
	err := r.db.QueryRow(query, slug).Scan(
		&problem.ID,
		&problem.Title,
		&problem.Slug,
		&problem.Description,
		&problem.Difficulty,
		&problem.Tags,
		&problem.Examples,
		&problem.Constraints,
		&problem.TemplateCode,
		&problem.CreatedAt,
		&problem.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewRepositoryError("GetBySlug", ErrNotFound, "problem_not_found")
		}
		return nil, NewRepositoryError("GetBySlug", err, "database_error")
	}

	return &problem, nil
}

// Update updates an existing problem
func (r *problemRepository) Update(problem *models.Problem) (*models.Problem, error) {
	query := `
		UPDATE problems
		SET title = $2, slug = $3, description = $4, difficulty = $5, tags = $6, 
		    examples = $7, constraints = $8, template_code = $9, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, title, slug, description, difficulty, tags, examples, constraints, template_code, created_at, updated_at`

	var updated models.Problem
	err := r.db.QueryRow(
		query,
		problem.ID,
		problem.Title,
		problem.Slug,
		problem.Description,
		problem.Difficulty,
		problem.Tags,
		problem.Examples,
		problem.Constraints,
		problem.TemplateCode,
	).Scan(
		&updated.ID,
		&updated.Title,
		&updated.Slug,
		&updated.Description,
		&updated.Difficulty,
		&updated.Tags,
		&updated.Examples,
		&updated.Constraints,
		&updated.TemplateCode,
		&updated.CreatedAt,
		&updated.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewRepositoryError("Update", ErrNotFound, "problem_not_found")
		}
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code {
			case "23505": // unique_violation
				if strings.Contains(pqErr.Detail, "slug") {
					return nil, NewRepositoryError("Update", ErrDuplicateKey, "slug_exists")
				}
				return nil, NewRepositoryError("Update", ErrDuplicateKey, "duplicate_key")
			case "23514": // check_violation
				return nil, NewRepositoryError("Update", ErrInvalidInput, "invalid_difficulty")
			}
		}
		return nil, NewRepositoryError("Update", err, "database_error")
	}

	return &updated, nil
}

// Delete deletes a problem by ID
func (r *problemRepository) Delete(id int) error {
	query := `DELETE FROM problems WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return NewRepositoryError("Delete", err, "database_error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return NewRepositoryError("Delete", err, "database_error")
	}

	if rowsAffected == 0 {
		return NewRepositoryError("Delete", ErrNotFound, "problem_not_found")
	}

	return nil
}

// List retrieves problems with filters
func (r *problemRepository) List(filters ProblemFilters) ([]*models.Problem, error) {
	query := `
		SELECT id, title, slug, description, difficulty, tags, examples, constraints, template_code, created_at, updated_at
		FROM problems`
	
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Apply difficulty filters
	if len(filters.Difficulty) > 0 {
		placeholders := make([]string, len(filters.Difficulty))
		for i, difficulty := range filters.Difficulty {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, difficulty)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("difficulty IN (%s)", strings.Join(placeholders, ",")))
	}

	// Apply tag filters
	if len(filters.Tags) > 0 {
		tagConditions := make([]string, len(filters.Tags))
		for i, tag := range filters.Tags {
			tagConditions[i] = fmt.Sprintf("$%d = ANY(tags)", argIndex)
			args = append(args, tag)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(tagConditions, " OR ")))
	}

	// Add WHERE clause if conditions exist
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add sorting
	sortBy := "created_at"
	if filters.SortBy != "" {
		switch filters.SortBy {
		case "title", "difficulty", "created_at":
			sortBy = filters.SortBy
		}
	}

	sortOrder := "DESC"
	if filters.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	query += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Add pagination
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filters.Limit)
		argIndex++
	}

	if filters.Offset > 0 {
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filters.Offset)
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, NewRepositoryError("List", err, "database_error")
	}
	defer rows.Close()

	var problems []*models.Problem
	for rows.Next() {
		var problem models.Problem
		err := rows.Scan(
			&problem.ID,
			&problem.Title,
			&problem.Slug,
			&problem.Description,
			&problem.Difficulty,
			&problem.Tags,
			&problem.Examples,
			&problem.Constraints,
			&problem.TemplateCode,
			&problem.CreatedAt,
			&problem.UpdatedAt,
		)
		if err != nil {
			return nil, NewRepositoryError("List", err, "scan_error")
		}
		problems = append(problems, &problem)
	}

	if err = rows.Err(); err != nil {
		return nil, NewRepositoryError("List", err, "rows_error")
	}

	return problems, nil
}

// Search searches problems by title or description
func (r *problemRepository) Search(query string, filters ProblemFilters) ([]*models.Problem, error) {
	sqlQuery := `
		SELECT id, title, slug, description, difficulty, tags, examples, constraints, template_code, created_at, updated_at
		FROM problems
		WHERE (title ILIKE $1 OR description ILIKE $1)`
	
	var conditions []string
	var args []interface{}
	args = append(args, "%"+query+"%")
	argIndex := 2

	// Apply difficulty filters
	if len(filters.Difficulty) > 0 {
		placeholders := make([]string, len(filters.Difficulty))
		for i, difficulty := range filters.Difficulty {
			placeholders[i] = fmt.Sprintf("$%d", argIndex)
			args = append(args, difficulty)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("difficulty IN (%s)", strings.Join(placeholders, ",")))
	}

	// Apply tag filters
	if len(filters.Tags) > 0 {
		tagConditions := make([]string, len(filters.Tags))
		for i, tag := range filters.Tags {
			tagConditions[i] = fmt.Sprintf("$%d = ANY(tags)", argIndex)
			args = append(args, tag)
			argIndex++
		}
		conditions = append(conditions, fmt.Sprintf("(%s)", strings.Join(tagConditions, " OR ")))
	}

	// Add additional conditions
	if len(conditions) > 0 {
		sqlQuery += " AND " + strings.Join(conditions, " AND ")
	}

	// Add sorting
	sortBy := "created_at"
	if filters.SortBy != "" {
		switch filters.SortBy {
		case "title", "difficulty", "created_at":
			sortBy = filters.SortBy
		}
	}

	sortOrder := "DESC"
	if filters.SortOrder == "asc" {
		sortOrder = "ASC"
	}

	sqlQuery += fmt.Sprintf(" ORDER BY %s %s", sortBy, sortOrder)

	// Add pagination
	if filters.Limit > 0 {
		sqlQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filters.Limit)
		argIndex++
	}

	if filters.Offset > 0 {
		sqlQuery += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, filters.Offset)
	}

	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, NewRepositoryError("Search", err, "database_error")
	}
	defer rows.Close()

	var problems []*models.Problem
	for rows.Next() {
		var problem models.Problem
		err := rows.Scan(
			&problem.ID,
			&problem.Title,
			&problem.Slug,
			&problem.Description,
			&problem.Difficulty,
			&problem.Tags,
			&problem.Examples,
			&problem.Constraints,
			&problem.TemplateCode,
			&problem.CreatedAt,
			&problem.UpdatedAt,
		)
		if err != nil {
			return nil, NewRepositoryError("Search", err, "scan_error")
		}
		problems = append(problems, &problem)
	}

	if err = rows.Err(); err != nil {
		return nil, NewRepositoryError("Search", err, "rows_error")
	}

	return problems, nil
}