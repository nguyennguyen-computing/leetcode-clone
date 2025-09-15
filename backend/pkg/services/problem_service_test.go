package services

import (
	"testing"

	"leetcode-clone-backend/pkg/models"
	"leetcode-clone-backend/pkg/repository"
)

// Mock repositories for testing
type mockProblemRepository struct {
	problems map[int]*models.Problem
	nextID   int
}

func newMockProblemRepository() *mockProblemRepository {
	return &mockProblemRepository{
		problems: make(map[int]*models.Problem),
		nextID:   1,
	}
}

func (m *mockProblemRepository) Create(problem *models.Problem) (*models.Problem, error) {
	problem.ID = m.nextID
	m.nextID++
	m.problems[problem.ID] = problem
	return problem, nil
}

func (m *mockProblemRepository) GetByID(id int) (*models.Problem, error) {
	if problem, exists := m.problems[id]; exists {
		return problem, nil
	}
	return nil, repository.NewRepositoryError("GetByID", repository.ErrNotFound, "problem_not_found")
}

func (m *mockProblemRepository) GetBySlug(slug string) (*models.Problem, error) {
	for _, problem := range m.problems {
		if problem.Slug == slug {
			return problem, nil
		}
	}
	return nil, repository.NewRepositoryError("GetBySlug", repository.ErrNotFound, "problem_not_found")
}

func (m *mockProblemRepository) Update(problem *models.Problem) (*models.Problem, error) {
	if _, exists := m.problems[problem.ID]; !exists {
		return nil, repository.NewRepositoryError("Update", repository.ErrNotFound, "problem_not_found")
	}
	m.problems[problem.ID] = problem
	return problem, nil
}

func (m *mockProblemRepository) Delete(id int) error {
	if _, exists := m.problems[id]; !exists {
		return repository.NewRepositoryError("Delete", repository.ErrNotFound, "problem_not_found")
	}
	delete(m.problems, id)
	return nil
}

func (m *mockProblemRepository) List(filters repository.ProblemFilters) ([]*models.Problem, error) {
	var result []*models.Problem
	for _, problem := range m.problems {
		result = append(result, problem)
	}
	return result, nil
}

func (m *mockProblemRepository) Search(query string, filters repository.ProblemFilters) ([]*models.Problem, error) {
	var result []*models.Problem
	for _, problem := range m.problems {
		result = append(result, problem)
	}
	return result, nil
}

type mockTestCaseRepository struct {
	testCases map[int]*models.TestCase
	nextID    int
}

func newMockTestCaseRepository() *mockTestCaseRepository {
	return &mockTestCaseRepository{
		testCases: make(map[int]*models.TestCase),
		nextID:    1,
	}
}

func (m *mockTestCaseRepository) Create(testCase *models.TestCase) (*models.TestCase, error) {
	testCase.ID = m.nextID
	m.nextID++
	m.testCases[testCase.ID] = testCase
	return testCase, nil
}

func (m *mockTestCaseRepository) GetByID(id int) (*models.TestCase, error) {
	if testCase, exists := m.testCases[id]; exists {
		return testCase, nil
	}
	return nil, repository.NewRepositoryError("GetByID", repository.ErrNotFound, "testcase_not_found")
}

func (m *mockTestCaseRepository) GetByProblemID(problemID int) ([]*models.TestCase, error) {
	var result []*models.TestCase
	for _, testCase := range m.testCases {
		if testCase.ProblemID == problemID {
			result = append(result, testCase)
		}
	}
	return result, nil
}

func (m *mockTestCaseRepository) GetPublicByProblemID(problemID int) ([]*models.TestCase, error) {
	var result []*models.TestCase
	for _, testCase := range m.testCases {
		if testCase.ProblemID == problemID && !testCase.IsHidden {
			result = append(result, testCase)
		}
	}
	return result, nil
}

func (m *mockTestCaseRepository) Update(testCase *models.TestCase) (*models.TestCase, error) {
	if _, exists := m.testCases[testCase.ID]; !exists {
		return nil, repository.NewRepositoryError("Update", repository.ErrNotFound, "testcase_not_found")
	}
	m.testCases[testCase.ID] = testCase
	return testCase, nil
}

func (m *mockTestCaseRepository) Delete(id int) error {
	if _, exists := m.testCases[id]; !exists {
		return repository.NewRepositoryError("Delete", repository.ErrNotFound, "testcase_not_found")
	}
	delete(m.testCases, id)
	return nil
}

func (m *mockTestCaseRepository) DeleteByProblemID(problemID int) error {
	for id, testCase := range m.testCases {
		if testCase.ProblemID == problemID {
			delete(m.testCases, id)
		}
	}
	return nil
}

func TestProblemService_CreateProblem(t *testing.T) {
	problemRepo := newMockProblemRepository()
	testCaseRepo := newMockTestCaseRepository()
	service := NewProblemService(problemRepo, testCaseRepo)

	problem := &models.Problem{
		Title:       "Two Sum",
		Description: "Find two numbers that add up to target",
		Difficulty:  models.DifficultyEasy,
		Examples: models.Examples{
			{Input: "[2,7,11,15], 9", Output: "[0,1]"},
		},
		TemplateCode: models.TemplateCode{
			models.LanguageJavaScript: "function twoSum(nums, target) {\n    // Your code here\n}",
		},
	}

	created, err := service.CreateProblem(problem)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if created.ID == 0 {
		t.Error("Expected problem ID to be set")
	}

	if created.Slug == "" {
		t.Error("Expected slug to be generated")
	}

	if created.Slug != "two-sum" {
		t.Errorf("Expected slug 'two-sum', got '%s'", created.Slug)
	}
}

func TestProblemService_CreateProblem_ValidationError(t *testing.T) {
	problemRepo := newMockProblemRepository()
	testCaseRepo := newMockTestCaseRepository()
	service := NewProblemService(problemRepo, testCaseRepo)

	// Test with empty title
	problem := &models.Problem{
		Title:       "",
		Description: "Test description",
		Difficulty:  models.DifficultyEasy,
	}

	_, err := service.CreateProblem(problem)
	if err == nil {
		t.Error("Expected validation error for empty title")
	}
}

func TestProblemService_GetProblem(t *testing.T) {
	problemRepo := newMockProblemRepository()
	testCaseRepo := newMockTestCaseRepository()
	service := NewProblemService(problemRepo, testCaseRepo)

	// Create a problem first
	problem := &models.Problem{
		Title:       "Test Problem",
		Description: "Test description",
		Difficulty:  models.DifficultyEasy,
		Examples: models.Examples{
			{Input: "test", Output: "test"},
		},
		TemplateCode: models.TemplateCode{
			models.LanguageJavaScript: "function test() {}",
		},
	}

	created, err := service.CreateProblem(problem)
	if err != nil {
		t.Fatalf("Failed to create problem: %v", err)
	}

	// Get the problem
	retrieved, err := service.GetProblem(created.ID)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrieved.ID != created.ID {
		t.Errorf("Expected ID %d, got %d", created.ID, retrieved.ID)
	}

	if retrieved.Title != created.Title {
		t.Errorf("Expected title '%s', got '%s'", created.Title, retrieved.Title)
	}
}

func TestProblemService_CreateTestCase(t *testing.T) {
	problemRepo := newMockProblemRepository()
	testCaseRepo := newMockTestCaseRepository()
	service := NewProblemService(problemRepo, testCaseRepo)

	// Create a problem first
	problem := &models.Problem{
		Title:       "Test Problem",
		Description: "Test description",
		Difficulty:  models.DifficultyEasy,
		Examples: models.Examples{
			{Input: "test", Output: "test"},
		},
		TemplateCode: models.TemplateCode{
			models.LanguageJavaScript: "function test() {}",
		},
	}

	created, err := service.CreateProblem(problem)
	if err != nil {
		t.Fatalf("Failed to create problem: %v", err)
	}

	// Create a test case
	testCase := &models.TestCase{
		ProblemID:      created.ID,
		Input:          "[1,2,3]",
		ExpectedOutput: "6",
		IsHidden:       false,
	}

	createdTestCase, err := service.CreateTestCase(testCase)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if createdTestCase.ID == 0 {
		t.Error("Expected test case ID to be set")
	}

	if createdTestCase.ProblemID != created.ID {
		t.Errorf("Expected problem ID %d, got %d", created.ID, createdTestCase.ProblemID)
	}
}

func TestProblemService_ValidateFilters(t *testing.T) {
	problemRepo := newMockProblemRepository()
	testCaseRepo := newMockTestCaseRepository()
	service := NewProblemService(problemRepo, testCaseRepo)

	// Test valid filters
	filters := repository.ProblemFilters{
		Difficulty: []string{models.DifficultyEasy, models.DifficultyMedium},
		Limit:      10,
		Offset:     0,
		SortBy:     "title",
		SortOrder:  "asc",
	}

	err := service.validateFilters(&filters)
	if err != nil {
		t.Errorf("Expected no error for valid filters, got %v", err)
	}

	// Test invalid difficulty
	invalidFilters := repository.ProblemFilters{
		Difficulty: []string{"Invalid"},
	}

	err = service.validateFilters(&invalidFilters)
	if err == nil {
		t.Error("Expected error for invalid difficulty")
	}

	// Test limit bounds
	limitFilters := repository.ProblemFilters{
		Limit: 200, // Should be capped at 100
	}

	err = service.validateFilters(&limitFilters)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if limitFilters.Limit != 100 {
		t.Errorf("Expected limit to be capped at 100, got %d", limitFilters.Limit)
	}
}

func TestProblemService_GenerateSlug(t *testing.T) {
	problemRepo := newMockProblemRepository()
	testCaseRepo := newMockTestCaseRepository()
	service := NewProblemService(problemRepo, testCaseRepo)

	tests := []struct {
		title    string
		expected string
	}{
		{"Two Sum", "two-sum"},
		{"Add Two Numbers", "add-two-numbers"},
		{"Longest Substring Without Repeating Characters", "longest-substring-without-repeating-characters"},
		{"3Sum", "3sum"},
		{"Valid Parentheses", "valid-parentheses"},
	}

	for _, test := range tests {
		result := service.generateSlug(test.title)
		if result != test.expected {
			t.Errorf("For title '%s', expected slug '%s', got '%s'", test.title, test.expected, result)
		}
	}
}