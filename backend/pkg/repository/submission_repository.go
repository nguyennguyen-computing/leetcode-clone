package repository

import (
	"database/sql"

	"leetcode-clone-backend/pkg/models"
)

// submissionRepository implements SubmissionRepository interface
type submissionRepository struct {
	db *sql.DB
}

// NewSubmissionRepository creates a new submission repository
func NewSubmissionRepository(db *sql.DB) SubmissionRepository {
	return &submissionRepository{db: db}
}

// Create creates a new submission
func (r *submissionRepository) Create(submission *models.Submission) (*models.Submission, error) {
	query := `
		INSERT INTO submissions (user_id, problem_id, language, code, status, runtime_ms, memory_kb, 
		                        test_cases_passed, total_test_cases, error_message)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, user_id, problem_id, language, code, status, runtime_ms, memory_kb, 
		          test_cases_passed, total_test_cases, error_message, submitted_at`

	var created models.Submission
	err := r.db.QueryRow(
		query,
		submission.UserID,
		submission.ProblemID,
		submission.Language,
		submission.Code,
		submission.Status,
		submission.RuntimeMs,
		submission.MemoryKb,
		submission.TestCasesPassed,
		submission.TotalTestCases,
		submission.ErrorMessage,
	).Scan(
		&created.ID,
		&created.UserID,
		&created.ProblemID,
		&created.Language,
		&created.Code,
		&created.Status,
		&created.RuntimeMs,
		&created.MemoryKb,
		&created.TestCasesPassed,
		&created.TotalTestCases,
		&created.ErrorMessage,
		&created.SubmittedAt,
	)

	if err != nil {
		return nil, NewRepositoryError("Create", err, "database_error")
	}

	return &created, nil
}

// GetByID retrieves a submission by ID
func (r *submissionRepository) GetByID(id int) (*models.Submission, error) {
	query := `
		SELECT id, user_id, problem_id, language, code, status, runtime_ms, memory_kb, 
		       test_cases_passed, total_test_cases, error_message, submitted_at
		FROM submissions
		WHERE id = $1`

	var submission models.Submission
	err := r.db.QueryRow(query, id).Scan(
		&submission.ID,
		&submission.UserID,
		&submission.ProblemID,
		&submission.Language,
		&submission.Code,
		&submission.Status,
		&submission.RuntimeMs,
		&submission.MemoryKb,
		&submission.TestCasesPassed,
		&submission.TotalTestCases,
		&submission.ErrorMessage,
		&submission.SubmittedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewRepositoryError("GetByID", ErrNotFound, "submission_not_found")
		}
		return nil, NewRepositoryError("GetByID", err, "database_error")
	}

	return &submission, nil
}

// GetByUserID retrieves submissions by user ID with pagination
func (r *submissionRepository) GetByUserID(userID int, limit, offset int) ([]*models.Submission, error) {
	query := `
		SELECT id, user_id, problem_id, language, code, status, runtime_ms, memory_kb, 
		       test_cases_passed, total_test_cases, error_message, submitted_at
		FROM submissions
		WHERE user_id = $1
		ORDER BY submitted_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, NewRepositoryError("GetByUserID", err, "database_error")
	}
	defer rows.Close()

	var submissions []*models.Submission
	for rows.Next() {
		var submission models.Submission
		err := rows.Scan(
			&submission.ID,
			&submission.UserID,
			&submission.ProblemID,
			&submission.Language,
			&submission.Code,
			&submission.Status,
			&submission.RuntimeMs,
			&submission.MemoryKb,
			&submission.TestCasesPassed,
			&submission.TotalTestCases,
			&submission.ErrorMessage,
			&submission.SubmittedAt,
		)
		if err != nil {
			return nil, NewRepositoryError("GetByUserID", err, "scan_error")
		}
		submissions = append(submissions, &submission)
	}

	if err = rows.Err(); err != nil {
		return nil, NewRepositoryError("GetByUserID", err, "rows_error")
	}

	return submissions, nil
}

// GetByProblemID retrieves submissions by problem ID with pagination
func (r *submissionRepository) GetByProblemID(problemID int, limit, offset int) ([]*models.Submission, error) {
	query := `
		SELECT id, user_id, problem_id, language, code, status, runtime_ms, memory_kb, 
		       test_cases_passed, total_test_cases, error_message, submitted_at
		FROM submissions
		WHERE problem_id = $1
		ORDER BY submitted_at DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, problemID, limit, offset)
	if err != nil {
		return nil, NewRepositoryError("GetByProblemID", err, "database_error")
	}
	defer rows.Close()

	var submissions []*models.Submission
	for rows.Next() {
		var submission models.Submission
		err := rows.Scan(
			&submission.ID,
			&submission.UserID,
			&submission.ProblemID,
			&submission.Language,
			&submission.Code,
			&submission.Status,
			&submission.RuntimeMs,
			&submission.MemoryKb,
			&submission.TestCasesPassed,
			&submission.TotalTestCases,
			&submission.ErrorMessage,
			&submission.SubmittedAt,
		)
		if err != nil {
			return nil, NewRepositoryError("GetByProblemID", err, "scan_error")
		}
		submissions = append(submissions, &submission)
	}

	if err = rows.Err(); err != nil {
		return nil, NewRepositoryError("GetByProblemID", err, "rows_error")
	}

	return submissions, nil
}

// GetByUserAndProblem retrieves submissions by user and problem with pagination
func (r *submissionRepository) GetByUserAndProblem(userID, problemID int, limit, offset int) ([]*models.Submission, error) {
	query := `
		SELECT id, user_id, problem_id, language, code, status, runtime_ms, memory_kb, 
		       test_cases_passed, total_test_cases, error_message, submitted_at
		FROM submissions
		WHERE user_id = $1 AND problem_id = $2
		ORDER BY submitted_at DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.db.Query(query, userID, problemID, limit, offset)
	if err != nil {
		return nil, NewRepositoryError("GetByUserAndProblem", err, "database_error")
	}
	defer rows.Close()

	var submissions []*models.Submission
	for rows.Next() {
		var submission models.Submission
		err := rows.Scan(
			&submission.ID,
			&submission.UserID,
			&submission.ProblemID,
			&submission.Language,
			&submission.Code,
			&submission.Status,
			&submission.RuntimeMs,
			&submission.MemoryKb,
			&submission.TestCasesPassed,
			&submission.TotalTestCases,
			&submission.ErrorMessage,
			&submission.SubmittedAt,
		)
		if err != nil {
			return nil, NewRepositoryError("GetByUserAndProblem", err, "scan_error")
		}
		submissions = append(submissions, &submission)
	}

	if err = rows.Err(); err != nil {
		return nil, NewRepositoryError("GetByUserAndProblem", err, "rows_error")
	}

	return submissions, nil
}

// Update updates an existing submission
func (r *submissionRepository) Update(submission *models.Submission) (*models.Submission, error) {
	query := `
		UPDATE submissions
		SET user_id = $2, problem_id = $3, language = $4, code = $5, status = $6, 
		    runtime_ms = $7, memory_kb = $8, test_cases_passed = $9, total_test_cases = $10, 
		    error_message = $11
		WHERE id = $1
		RETURNING id, user_id, problem_id, language, code, status, runtime_ms, memory_kb, 
		          test_cases_passed, total_test_cases, error_message, submitted_at`

	var updated models.Submission
	err := r.db.QueryRow(
		query,
		submission.ID,
		submission.UserID,
		submission.ProblemID,
		submission.Language,
		submission.Code,
		submission.Status,
		submission.RuntimeMs,
		submission.MemoryKb,
		submission.TestCasesPassed,
		submission.TotalTestCases,
		submission.ErrorMessage,
	).Scan(
		&updated.ID,
		&updated.UserID,
		&updated.ProblemID,
		&updated.Language,
		&updated.Code,
		&updated.Status,
		&updated.RuntimeMs,
		&updated.MemoryKb,
		&updated.TestCasesPassed,
		&updated.TotalTestCases,
		&updated.ErrorMessage,
		&updated.SubmittedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewRepositoryError("Update", ErrNotFound, "submission_not_found")
		}
		return nil, NewRepositoryError("Update", err, "database_error")
	}

	return &updated, nil
}

// Delete deletes a submission by ID
func (r *submissionRepository) Delete(id int) error {
	query := `DELETE FROM submissions WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return NewRepositoryError("Delete", err, "database_error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return NewRepositoryError("Delete", err, "database_error")
	}

	if rowsAffected == 0 {
		return NewRepositoryError("Delete", ErrNotFound, "submission_not_found")
	}

	return nil
}

// GetLatestByUserAndProblem retrieves the latest submission for a user and problem
func (r *submissionRepository) GetLatestByUserAndProblem(userID, problemID int) (*models.Submission, error) {
	query := `
		SELECT id, user_id, problem_id, language, code, status, runtime_ms, memory_kb, 
		       test_cases_passed, total_test_cases, error_message, submitted_at
		FROM submissions
		WHERE user_id = $1 AND problem_id = $2
		ORDER BY submitted_at DESC
		LIMIT 1`

	var submission models.Submission
	err := r.db.QueryRow(query, userID, problemID).Scan(
		&submission.ID,
		&submission.UserID,
		&submission.ProblemID,
		&submission.Language,
		&submission.Code,
		&submission.Status,
		&submission.RuntimeMs,
		&submission.MemoryKb,
		&submission.TestCasesPassed,
		&submission.TotalTestCases,
		&submission.ErrorMessage,
		&submission.SubmittedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewRepositoryError("GetLatestByUserAndProblem", ErrNotFound, "submission_not_found")
		}
		return nil, NewRepositoryError("GetLatestByUserAndProblem", err, "database_error")
	}

	return &submission, nil
}