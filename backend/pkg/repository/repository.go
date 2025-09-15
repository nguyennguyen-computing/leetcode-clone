package repository

import (
	"database/sql"
)

// NewRepository creates a new repository instance with all sub-repositories
func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:         NewUserRepository(db),
		Problem:      NewProblemRepository(db),
		TestCase:     NewTestCaseRepository(db),
		Submission:   NewSubmissionRepository(db),
		UserProgress: NewUserProgressRepository(db),
	}
}