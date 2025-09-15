package models

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"

	"github.com/lib/pq"
)

// User represents a user in the system
type User struct {
	ID           int       `json:"id" db:"id"`
	Username     string    `json:"username" db:"username"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"` // Don't include in JSON responses
	IsAdmin      bool      `json:"is_admin" db:"is_admin"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// Example represents a problem example with input and output
type Example struct {
	Input       string `json:"input"`
	Output      string `json:"output"`
	Explanation string `json:"explanation,omitempty"`
}

// Examples is a slice of Example that implements sql.Scanner and driver.Valuer
type Examples []Example

// Scan implements the sql.Scanner interface for Examples
func (e *Examples) Scan(value interface{}) error {
	if value == nil {
		*e = Examples{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into Examples", value)
	}

	return json.Unmarshal(bytes, e)
}

// Value implements the driver.Valuer interface for Examples
func (e Examples) Value() (driver.Value, error) {
	if len(e) == 0 {
		return json.Marshal([]Example{})
	}
	return json.Marshal(e)
}

// TemplateCode represents code templates for different programming languages
type TemplateCode map[string]string

// Scan implements the sql.Scanner interface for TemplateCode
func (tc *TemplateCode) Scan(value interface{}) error {
	if value == nil {
		*tc = TemplateCode{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("cannot scan %T into TemplateCode", value)
	}

	return json.Unmarshal(bytes, tc)
}

// Value implements the driver.Valuer interface for TemplateCode
func (tc TemplateCode) Value() (driver.Value, error) {
	if len(tc) == 0 {
		return json.Marshal(map[string]string{})
	}
	return json.Marshal(tc)
}

// Problem represents a coding problem
type Problem struct {
	ID           int          `json:"id" db:"id"`
	Title        string       `json:"title" db:"title"`
	Slug         string       `json:"slug" db:"slug"`
	Description  string       `json:"description" db:"description"`
	Difficulty   string       `json:"difficulty" db:"difficulty"`
	Tags         pq.StringArray `json:"tags" db:"tags"`
	Examples     Examples     `json:"examples" db:"examples"`
	Constraints  string       `json:"constraints" db:"constraints"`
	TemplateCode TemplateCode `json:"template_code" db:"template_code"`
	CreatedAt    time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at" db:"updated_at"`
}

// TestCase represents a test case for a problem
type TestCase struct {
	ID             int       `json:"id" db:"id"`
	ProblemID      int       `json:"problem_id" db:"problem_id"`
	Input          string    `json:"input" db:"input"`
	ExpectedOutput string    `json:"expected_output" db:"expected_output"`
	IsHidden       bool      `json:"is_hidden" db:"is_hidden"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// Submission represents a code submission
type Submission struct {
	ID               int        `json:"id" db:"id"`
	UserID           int        `json:"user_id" db:"user_id"`
	ProblemID        int        `json:"problem_id" db:"problem_id"`
	Language         string     `json:"language" db:"language"`
	Code             string     `json:"code" db:"code"`
	Status           string     `json:"status" db:"status"`
	RuntimeMs        *int       `json:"runtime_ms" db:"runtime_ms"`
	MemoryKb         *int       `json:"memory_kb" db:"memory_kb"`
	TestCasesPassed  int        `json:"test_cases_passed" db:"test_cases_passed"`
	TotalTestCases   int        `json:"total_test_cases" db:"total_test_cases"`
	ErrorMessage     *string    `json:"error_message" db:"error_message"`
	SubmittedAt      time.Time  `json:"submitted_at" db:"submitted_at"`
}

// UserProgress represents a user's progress on a specific problem
type UserProgress struct {
	UserID           int        `json:"user_id" db:"user_id"`
	ProblemID        int        `json:"problem_id" db:"problem_id"`
	IsSolved         bool       `json:"is_solved" db:"is_solved"`
	BestSubmissionID *int       `json:"best_submission_id" db:"best_submission_id"`
	Attempts         int        `json:"attempts" db:"attempts"`
	FirstSolvedAt    *time.Time `json:"first_solved_at" db:"first_solved_at"`
}

// Submission status constants
const (
	StatusAccepted           = "Accepted"
	StatusWrongAnswer        = "Wrong Answer"
	StatusTimeLimitExceeded  = "Time Limit Exceeded"
	StatusMemoryLimitExceeded = "Memory Limit Exceeded"
	StatusRuntimeError       = "Runtime Error"
	StatusCompileError       = "Compile Error"
	StatusInternalError      = "Internal Error"
)

// Problem difficulty constants
const (
	DifficultyEasy   = "Easy"
	DifficultyMedium = "Medium"
	DifficultyHard   = "Hard"
)

// Supported programming languages
const (
	LanguageJavaScript = "javascript"
	LanguagePython     = "python"
	LanguageJava       = "java"
)