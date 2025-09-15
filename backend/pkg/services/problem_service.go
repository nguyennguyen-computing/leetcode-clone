package services

import (
	"fmt"
	"regexp"
	"strings"

	"leetcode-clone-backend/pkg/models"
	"leetcode-clone-backend/pkg/repository"
)

// ProblemService handles business logic for problems
type ProblemService struct {
	problemRepo  repository.ProblemRepository
	testCaseRepo repository.TestCaseRepository
}

// NewProblemService creates a new problem service
func NewProblemService(problemRepo repository.ProblemRepository, testCaseRepo repository.TestCaseRepository) *ProblemService {
	return &ProblemService{
		problemRepo:  problemRepo,
		testCaseRepo: testCaseRepo,
	}
}

// CreateProblem creates a new problem with validation
func (s *ProblemService) CreateProblem(problem *models.Problem) (*models.Problem, error) {
	// Validate problem data
	if err := s.validateProblem(problem); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Generate slug if not provided
	if problem.Slug == "" {
		problem.Slug = s.generateSlug(problem.Title)
	}

	// Create the problem
	created, err := s.problemRepo.Create(problem)
	if err != nil {
		return nil, fmt.Errorf("failed to create problem: %w", err)
	}

	return created, nil
}

// GetProblem retrieves a problem by ID
func (s *ProblemService) GetProblem(id int) (*models.Problem, error) {
	problem, err := s.problemRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get problem: %w", err)
	}

	return problem, nil
}

// GetProblemBySlug retrieves a problem by slug
func (s *ProblemService) GetProblemBySlug(slug string) (*models.Problem, error) {
	problem, err := s.problemRepo.GetBySlug(slug)
	if err != nil {
		return nil, fmt.Errorf("failed to get problem by slug: %w", err)
	}

	return problem, nil
}

// UpdateProblem updates an existing problem
func (s *ProblemService) UpdateProblem(problem *models.Problem) (*models.Problem, error) {
	// Validate problem data
	if err := s.validateProblem(problem); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Update the problem
	updated, err := s.problemRepo.Update(problem)
	if err != nil {
		return nil, fmt.Errorf("failed to update problem: %w", err)
	}

	return updated, nil
}

// DeleteProblem deletes a problem by ID
func (s *ProblemService) DeleteProblem(id int) error {
	// First delete all test cases for this problem
	if err := s.testCaseRepo.DeleteByProblemID(id); err != nil {
		return fmt.Errorf("failed to delete test cases: %w", err)
	}

	// Then delete the problem
	if err := s.problemRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete problem: %w", err)
	}

	return nil
}

// ListProblems retrieves problems with filters
func (s *ProblemService) ListProblems(filters repository.ProblemFilters) ([]*models.Problem, error) {
	// Validate filters
	if err := s.validateFilters(&filters); err != nil {
		return nil, fmt.Errorf("invalid filters: %w", err)
	}

	problems, err := s.problemRepo.List(filters)
	if err != nil {
		return nil, fmt.Errorf("failed to list problems: %w", err)
	}

	return problems, nil
}

// SearchProblems searches problems by title or description
func (s *ProblemService) SearchProblems(query string, filters repository.ProblemFilters) ([]*models.Problem, error) {
	// Validate search query
	if strings.TrimSpace(query) == "" {
		return nil, fmt.Errorf("search query cannot be empty")
	}

	// Validate filters
	if err := s.validateFilters(&filters); err != nil {
		return nil, fmt.Errorf("invalid filters: %w", err)
	}

	problems, err := s.problemRepo.Search(query, filters)
	if err != nil {
		return nil, fmt.Errorf("failed to search problems: %w", err)
	}

	return problems, nil
}

// CreateTestCase creates a new test case for a problem
func (s *ProblemService) CreateTestCase(testCase *models.TestCase) (*models.TestCase, error) {
	// Validate test case
	if err := s.validateTestCase(testCase); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Verify problem exists
	_, err := s.problemRepo.GetByID(testCase.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("problem not found: %w", err)
	}

	created, err := s.testCaseRepo.Create(testCase)
	if err != nil {
		return nil, fmt.Errorf("failed to create test case: %w", err)
	}

	return created, nil
}

// GetTestCases retrieves all test cases for a problem
func (s *ProblemService) GetTestCases(problemID int) ([]*models.TestCase, error) {
	testCases, err := s.testCaseRepo.GetByProblemID(problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get test cases: %w", err)
	}

	return testCases, nil
}

// GetPublicTestCases retrieves only public test cases for a problem
func (s *ProblemService) GetPublicTestCases(problemID int) ([]*models.TestCase, error) {
	testCases, err := s.testCaseRepo.GetPublicByProblemID(problemID)
	if err != nil {
		return nil, fmt.Errorf("failed to get public test cases: %w", err)
	}

	return testCases, nil
}

// UpdateTestCase updates an existing test case
func (s *ProblemService) UpdateTestCase(testCase *models.TestCase) (*models.TestCase, error) {
	// Validate test case
	if err := s.validateTestCase(testCase); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	updated, err := s.testCaseRepo.Update(testCase)
	if err != nil {
		return nil, fmt.Errorf("failed to update test case: %w", err)
	}

	return updated, nil
}

// DeleteTestCase deletes a test case by ID
func (s *ProblemService) DeleteTestCase(id int) error {
	if err := s.testCaseRepo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete test case: %w", err)
	}

	return nil
}

// validateProblem validates problem data
func (s *ProblemService) validateProblem(problem *models.Problem) error {
	if strings.TrimSpace(problem.Title) == "" {
		return fmt.Errorf("title is required")
	}

	if len(problem.Title) > 200 {
		return fmt.Errorf("title must be 200 characters or less")
	}

	if strings.TrimSpace(problem.Description) == "" {
		return fmt.Errorf("description is required")
	}

	// Validate difficulty
	validDifficulties := map[string]bool{
		models.DifficultyEasy:   true,
		models.DifficultyMedium: true,
		models.DifficultyHard:   true,
	}
	if !validDifficulties[problem.Difficulty] {
		return fmt.Errorf("difficulty must be one of: Easy, Medium, Hard")
	}

	// Validate examples
	if len(problem.Examples) == 0 {
		return fmt.Errorf("at least one example is required")
	}

	for i, example := range problem.Examples {
		if strings.TrimSpace(example.Input) == "" {
			return fmt.Errorf("example %d input is required", i+1)
		}
		if strings.TrimSpace(example.Output) == "" {
			return fmt.Errorf("example %d output is required", i+1)
		}
	}

	// Validate template code
	if len(problem.TemplateCode) == 0 {
		return fmt.Errorf("at least one template code is required")
	}

	validLanguages := map[string]bool{
		models.LanguageJavaScript: true,
		models.LanguagePython:     true,
		models.LanguageJava:       true,
	}

	for lang, code := range problem.TemplateCode {
		if !validLanguages[lang] {
			return fmt.Errorf("unsupported language: %s", lang)
		}
		if strings.TrimSpace(code) == "" {
			return fmt.Errorf("template code for %s cannot be empty", lang)
		}
	}

	return nil
}

// validateTestCase validates test case data
func (s *ProblemService) validateTestCase(testCase *models.TestCase) error {
	if testCase.ProblemID <= 0 {
		return fmt.Errorf("problem ID is required")
	}

	if strings.TrimSpace(testCase.Input) == "" {
		return fmt.Errorf("input is required")
	}

	if strings.TrimSpace(testCase.ExpectedOutput) == "" {
		return fmt.Errorf("expected output is required")
	}

	return nil
}

// validateFilters validates and sets defaults for problem filters
func (s *ProblemService) validateFilters(filters *repository.ProblemFilters) error {
	// Validate difficulty filters
	validDifficulties := map[string]bool{
		models.DifficultyEasy:   true,
		models.DifficultyMedium: true,
		models.DifficultyHard:   true,
	}

	for _, difficulty := range filters.Difficulty {
		if !validDifficulties[difficulty] {
			return fmt.Errorf("invalid difficulty: %s", difficulty)
		}
	}

	// Set default limit if not provided
	if filters.Limit <= 0 {
		filters.Limit = 50 // Default limit
	}

	// Validate limit bounds
	if filters.Limit > 100 {
		filters.Limit = 100 // Maximum limit
	}

	// Validate offset
	if filters.Offset < 0 {
		filters.Offset = 0
	}

	// Validate sort by
	if filters.SortBy != "" {
		validSortFields := map[string]bool{
			"title":      true,
			"difficulty": true,
			"created_at": true,
		}
		if !validSortFields[filters.SortBy] {
			return fmt.Errorf("invalid sort field: %s", filters.SortBy)
		}
	}

	// Validate sort order
	if filters.SortOrder != "" && filters.SortOrder != "asc" && filters.SortOrder != "desc" {
		return fmt.Errorf("invalid sort order: %s (must be 'asc' or 'desc')", filters.SortOrder)
	}

	return nil
}

// generateSlug generates a URL-friendly slug from a title
func (s *ProblemService) generateSlug(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces and special characters with hyphens
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	slug = reg.ReplaceAllString(slug, "-")

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	return slug
}