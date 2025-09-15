package repository

import (
	"database/sql"

	"leetcode-clone-backend/pkg/models"
)

// testCaseRepository implements TestCaseRepository interface
type testCaseRepository struct {
	db *sql.DB
}

// NewTestCaseRepository creates a new test case repository
func NewTestCaseRepository(db *sql.DB) TestCaseRepository {
	return &testCaseRepository{db: db}
}

// Create creates a new test case
func (r *testCaseRepository) Create(testCase *models.TestCase) (*models.TestCase, error) {
	query := `
		INSERT INTO test_cases (problem_id, input, expected_output, is_hidden)
		VALUES ($1, $2, $3, $4)
		RETURNING id, problem_id, input, expected_output, is_hidden, created_at`

	var created models.TestCase
	err := r.db.QueryRow(
		query,
		testCase.ProblemID,
		testCase.Input,
		testCase.ExpectedOutput,
		testCase.IsHidden,
	).Scan(
		&created.ID,
		&created.ProblemID,
		&created.Input,
		&created.ExpectedOutput,
		&created.IsHidden,
		&created.CreatedAt,
	)

	if err != nil {
		return nil, NewRepositoryError("Create", err, "database_error")
	}

	return &created, nil
}

// GetByID retrieves a test case by ID
func (r *testCaseRepository) GetByID(id int) (*models.TestCase, error) {
	query := `
		SELECT id, problem_id, input, expected_output, is_hidden, created_at
		FROM test_cases
		WHERE id = $1`

	var testCase models.TestCase
	err := r.db.QueryRow(query, id).Scan(
		&testCase.ID,
		&testCase.ProblemID,
		&testCase.Input,
		&testCase.ExpectedOutput,
		&testCase.IsHidden,
		&testCase.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewRepositoryError("GetByID", ErrNotFound, "testcase_not_found")
		}
		return nil, NewRepositoryError("GetByID", err, "database_error")
	}

	return &testCase, nil
}

// GetByProblemID retrieves all test cases for a problem
func (r *testCaseRepository) GetByProblemID(problemID int) ([]*models.TestCase, error) {
	query := `
		SELECT id, problem_id, input, expected_output, is_hidden, created_at
		FROM test_cases
		WHERE problem_id = $1
		ORDER BY id ASC`

	rows, err := r.db.Query(query, problemID)
	if err != nil {
		return nil, NewRepositoryError("GetByProblemID", err, "database_error")
	}
	defer rows.Close()

	var testCases []*models.TestCase
	for rows.Next() {
		var testCase models.TestCase
		err := rows.Scan(
			&testCase.ID,
			&testCase.ProblemID,
			&testCase.Input,
			&testCase.ExpectedOutput,
			&testCase.IsHidden,
			&testCase.CreatedAt,
		)
		if err != nil {
			return nil, NewRepositoryError("GetByProblemID", err, "scan_error")
		}
		testCases = append(testCases, &testCase)
	}

	if err = rows.Err(); err != nil {
		return nil, NewRepositoryError("GetByProblemID", err, "rows_error")
	}

	return testCases, nil
}

// GetPublicByProblemID retrieves only public (non-hidden) test cases for a problem
func (r *testCaseRepository) GetPublicByProblemID(problemID int) ([]*models.TestCase, error) {
	query := `
		SELECT id, problem_id, input, expected_output, is_hidden, created_at
		FROM test_cases
		WHERE problem_id = $1 AND is_hidden = false
		ORDER BY id ASC`

	rows, err := r.db.Query(query, problemID)
	if err != nil {
		return nil, NewRepositoryError("GetPublicByProblemID", err, "database_error")
	}
	defer rows.Close()

	var testCases []*models.TestCase
	for rows.Next() {
		var testCase models.TestCase
		err := rows.Scan(
			&testCase.ID,
			&testCase.ProblemID,
			&testCase.Input,
			&testCase.ExpectedOutput,
			&testCase.IsHidden,
			&testCase.CreatedAt,
		)
		if err != nil {
			return nil, NewRepositoryError("GetPublicByProblemID", err, "scan_error")
		}
		testCases = append(testCases, &testCase)
	}

	if err = rows.Err(); err != nil {
		return nil, NewRepositoryError("GetPublicByProblemID", err, "rows_error")
	}

	return testCases, nil
}

// Update updates an existing test case
func (r *testCaseRepository) Update(testCase *models.TestCase) (*models.TestCase, error) {
	query := `
		UPDATE test_cases
		SET problem_id = $2, input = $3, expected_output = $4, is_hidden = $5
		WHERE id = $1
		RETURNING id, problem_id, input, expected_output, is_hidden, created_at`

	var updated models.TestCase
	err := r.db.QueryRow(
		query,
		testCase.ID,
		testCase.ProblemID,
		testCase.Input,
		testCase.ExpectedOutput,
		testCase.IsHidden,
	).Scan(
		&updated.ID,
		&updated.ProblemID,
		&updated.Input,
		&updated.ExpectedOutput,
		&updated.IsHidden,
		&updated.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, NewRepositoryError("Update", ErrNotFound, "testcase_not_found")
		}
		return nil, NewRepositoryError("Update", err, "database_error")
	}

	return &updated, nil
}

// Delete deletes a test case by ID
func (r *testCaseRepository) Delete(id int) error {
	query := `DELETE FROM test_cases WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return NewRepositoryError("Delete", err, "database_error")
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return NewRepositoryError("Delete", err, "database_error")
	}

	if rowsAffected == 0 {
		return NewRepositoryError("Delete", ErrNotFound, "testcase_not_found")
	}

	return nil
}

// DeleteByProblemID deletes all test cases for a problem
func (r *testCaseRepository) DeleteByProblemID(problemID int) error {
	query := `DELETE FROM test_cases WHERE problem_id = $1`

	_, err := r.db.Exec(query, problemID)
	if err != nil {
		return NewRepositoryError("DeleteByProblemID", err, "database_error")
	}

	return nil
}