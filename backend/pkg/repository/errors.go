package repository

import (
	"errors"
	"fmt"
)

// Common repository errors
var (
	ErrNotFound      = errors.New("record not found")
	ErrDuplicateKey  = errors.New("duplicate key violation")
	ErrInvalidInput  = errors.New("invalid input")
	ErrDatabase      = errors.New("database error")
	ErrTransaction   = errors.New("transaction error")
)

// RepositoryError represents a repository-specific error
type RepositoryError struct {
	Op   string // Operation that failed
	Err  error  // Underlying error
	Code string // Error code for client handling
}

func (e *RepositoryError) Error() string {
	if e.Op == "" {
		return e.Err.Error()
	}
	return fmt.Sprintf("%s: %v", e.Op, e.Err)
}

func (e *RepositoryError) Unwrap() error {
	return e.Err
}

// NewRepositoryError creates a new repository error
func NewRepositoryError(op string, err error, code string) *RepositoryError {
	return &RepositoryError{
		Op:   op,
		Err:  err,
		Code: code,
	}
}

// IsNotFound checks if the error is a "not found" error
func IsNotFound(err error) bool {
	var repoErr *RepositoryError
	if errors.As(err, &repoErr) {
		return errors.Is(repoErr.Err, ErrNotFound)
	}
	return errors.Is(err, ErrNotFound)
}

// IsDuplicateKey checks if the error is a duplicate key error
func IsDuplicateKey(err error) bool {
	var repoErr *RepositoryError
	if errors.As(err, &repoErr) {
		return errors.Is(repoErr.Err, ErrDuplicateKey)
	}
	return errors.Is(err, ErrDuplicateKey)
}

// IsInvalidInput checks if the error is an invalid input error
func IsInvalidInput(err error) bool {
	var repoErr *RepositoryError
	if errors.As(err, &repoErr) {
		return errors.Is(repoErr.Err, ErrInvalidInput)
	}
	return errors.Is(err, ErrInvalidInput)
}