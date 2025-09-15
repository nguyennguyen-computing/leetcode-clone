package repository

import (
	"leetcode-clone-backend/pkg/models"
)

// UserRepository defines the interface for user data operations
type UserRepository interface {
	Create(user *models.User) (*models.User, error)
	GetByID(id int) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Update(user *models.User) (*models.User, error)
	Delete(id int) error
	List(limit, offset int) ([]*models.User, error)
}

// ProblemRepository defines the interface for problem data operations
type ProblemRepository interface {
	Create(problem *models.Problem) (*models.Problem, error)
	GetByID(id int) (*models.Problem, error)
	GetBySlug(slug string) (*models.Problem, error)
	Update(problem *models.Problem) (*models.Problem, error)
	Delete(id int) error
	List(filters ProblemFilters) ([]*models.Problem, error)
	Search(query string, filters ProblemFilters) ([]*models.Problem, error)
}

// TestCaseRepository defines the interface for test case data operations
type TestCaseRepository interface {
	Create(testCase *models.TestCase) (*models.TestCase, error)
	GetByID(id int) (*models.TestCase, error)
	GetByProblemID(problemID int) ([]*models.TestCase, error)
	GetPublicByProblemID(problemID int) ([]*models.TestCase, error)
	Update(testCase *models.TestCase) (*models.TestCase, error)
	Delete(id int) error
	DeleteByProblemID(problemID int) error
}

// SubmissionRepository defines the interface for submission data operations
type SubmissionRepository interface {
	Create(submission *models.Submission) (*models.Submission, error)
	GetByID(id int) (*models.Submission, error)
	GetByUserID(userID int, limit, offset int) ([]*models.Submission, error)
	GetByProblemID(problemID int, limit, offset int) ([]*models.Submission, error)
	GetByUserAndProblem(userID, problemID int, limit, offset int) ([]*models.Submission, error)
	Update(submission *models.Submission) (*models.Submission, error)
	Delete(id int) error
	GetLatestByUserAndProblem(userID, problemID int) (*models.Submission, error)
}

// UserProgressRepository defines the interface for user progress data operations
type UserProgressRepository interface {
	Create(progress *models.UserProgress) (*models.UserProgress, error)
	GetByUserAndProblem(userID, problemID int) (*models.UserProgress, error)
	GetByUserID(userID int) ([]*models.UserProgress, error)
	Update(progress *models.UserProgress) (*models.UserProgress, error)
	Delete(userID, problemID int) error
	GetSolvedCount(userID int) (int, error)
	GetSolvedCountByDifficulty(userID int) (map[string]int, error)
}

// ProblemFilters represents filters for problem queries
type ProblemFilters struct {
	Difficulty []string
	Tags       []string
	Limit      int
	Offset     int
	SortBy     string // "title", "difficulty", "created_at"
	SortOrder  string // "asc", "desc"
}

// Repository aggregates all repository interfaces
type Repository struct {
	User         UserRepository
	Problem      ProblemRepository
	TestCase     TestCaseRepository
	Submission   SubmissionRepository
	UserProgress UserProgressRepository
}