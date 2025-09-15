package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"leetcode-clone-backend/pkg/models"
	"leetcode-clone-backend/pkg/repository"
	"leetcode-clone-backend/pkg/services"
)

// Mock repositories for testing handlers
type mockProblemRepo struct {
	problems map[int]*models.Problem
	nextID   int
}

func newMockProblemRepo() *mockProblemRepo {
	return &mockProblemRepo{
		problems: make(map[int]*models.Problem),
		nextID:   1,
	}
}

func (m *mockProblemRepo) Create(problem *models.Problem) (*models.Problem, error) {
	problem.ID = m.nextID
	m.nextID++
	m.problems[problem.ID] = problem
	return problem, nil
}

func (m *mockProblemRepo) GetByID(id int) (*models.Problem, error) {
	if problem, exists := m.problems[id]; exists {
		return problem, nil
	}
	return nil, &mockError{message: "problem not found"}
}

func (m *mockProblemRepo) GetBySlug(slug string) (*models.Problem, error) {
	for _, problem := range m.problems {
		if problem.Slug == slug {
			return problem, nil
		}
	}
	return nil, &mockError{message: "problem not found"}
}

func (m *mockProblemRepo) Update(problem *models.Problem) (*models.Problem, error) {
	if _, exists := m.problems[problem.ID]; !exists {
		return nil, &mockError{message: "problem not found"}
	}
	m.problems[problem.ID] = problem
	return problem, nil
}

func (m *mockProblemRepo) Delete(id int) error {
	if _, exists := m.problems[id]; !exists {
		return &mockError{message: "problem not found"}
	}
	delete(m.problems, id)
	return nil
}

func (m *mockProblemRepo) List(filters repository.ProblemFilters) ([]*models.Problem, error) {
	var result []*models.Problem
	for _, problem := range m.problems {
		result = append(result, problem)
	}
	return result, nil
}

func (m *mockProblemRepo) Search(query string, filters repository.ProblemFilters) ([]*models.Problem, error) {
	var result []*models.Problem
	for _, problem := range m.problems {
		result = append(result, problem)
	}
	return result, nil
}

type mockTestCaseRepo struct {
	testCases map[int]*models.TestCase
	nextID    int
}

func newMockTestCaseRepo() *mockTestCaseRepo {
	return &mockTestCaseRepo{
		testCases: make(map[int]*models.TestCase),
		nextID:    1,
	}
}

func (m *mockTestCaseRepo) Create(testCase *models.TestCase) (*models.TestCase, error) {
	testCase.ID = m.nextID
	m.nextID++
	m.testCases[testCase.ID] = testCase
	return testCase, nil
}

func (m *mockTestCaseRepo) GetByID(id int) (*models.TestCase, error) {
	if testCase, exists := m.testCases[id]; exists {
		return testCase, nil
	}
	return nil, &mockError{message: "testcase not found"}
}

func (m *mockTestCaseRepo) GetByProblemID(problemID int) ([]*models.TestCase, error) {
	var result []*models.TestCase
	for _, testCase := range m.testCases {
		if testCase.ProblemID == problemID {
			result = append(result, testCase)
		}
	}
	return result, nil
}

func (m *mockTestCaseRepo) GetPublicByProblemID(problemID int) ([]*models.TestCase, error) {
	var result []*models.TestCase
	for _, testCase := range m.testCases {
		if testCase.ProblemID == problemID && !testCase.IsHidden {
			result = append(result, testCase)
		}
	}
	return result, nil
}

func (m *mockTestCaseRepo) Update(testCase *models.TestCase) (*models.TestCase, error) {
	if _, exists := m.testCases[testCase.ID]; !exists {
		return nil, &mockError{message: "testcase not found"}
	}
	m.testCases[testCase.ID] = testCase
	return testCase, nil
}

func (m *mockTestCaseRepo) Delete(id int) error {
	if _, exists := m.testCases[id]; !exists {
		return &mockError{message: "testcase not found"}
	}
	delete(m.testCases, id)
	return nil
}

func (m *mockTestCaseRepo) DeleteByProblemID(problemID int) error {
	for id, testCase := range m.testCases {
		if testCase.ProblemID == problemID {
			delete(m.testCases, id)
		}
	}
	return nil
}

type mockError struct {
	message string
}

func (e *mockError) Error() string {
	return e.message
}

func setupTestRouter() (*gin.Engine, *ProblemHandlers) {
	gin.SetMode(gin.TestMode)
	
	problemRepo := newMockProblemRepo()
	testCaseRepo := newMockTestCaseRepo()
	problemService := services.NewProblemService(problemRepo, testCaseRepo)
	problemHandler := NewProblemHandlers(problemService)
	
	router := gin.New()
	
	return router, problemHandler
}

func TestProblemHandlers_ListProblems(t *testing.T) {
	router, handler := setupTestRouter()
	router.GET("/problems", handler.ListProblems)

	req, _ := http.NewRequest("GET", "/problems", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if _, exists := response["problems"]; !exists {
		t.Error("Expected 'problems' field in response")
	}

	if _, exists := response["count"]; !exists {
		t.Error("Expected 'count' field in response")
	}
}

func TestProblemHandlers_CreateProblem(t *testing.T) {
	router, handler := setupTestRouter()
	router.POST("/problems", handler.CreateProblem)

	problem := models.Problem{
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

	jsonData, _ := json.Marshal(problem)
	req, _ := http.NewRequest("POST", "/problems", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		t.Logf("Response body: %s", w.Body.String())
	}

	var response models.Problem
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.ID == 0 {
		t.Error("Expected problem ID to be set")
	}

	if response.Title != problem.Title {
		t.Errorf("Expected title '%s', got '%s'", problem.Title, response.Title)
	}
}

func TestProblemHandlers_GetProblem(t *testing.T) {
	router, handler := setupTestRouter()
	router.POST("/problems", handler.CreateProblem)
	router.GET("/problems/:id", handler.GetProblem)

	// First create a problem
	problem := models.Problem{
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

	jsonData, _ := json.Marshal(problem)
	req, _ := http.NewRequest("POST", "/problems", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	var created models.Problem
	json.Unmarshal(w.Body.Bytes(), &created)

	// Now get the problem
	req, _ = http.NewRequest("GET", "/problems/1", nil)
	w = httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response models.Problem
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Title != problem.Title {
		t.Errorf("Expected title '%s', got '%s'", problem.Title, response.Title)
	}
}

func TestProblemHandlers_SearchProblems(t *testing.T) {
	router, handler := setupTestRouter()
	router.GET("/problems/search", handler.SearchProblems)

	req, _ := http.NewRequest("GET", "/problems/search?q=test", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	var response map[string]interface{}
	err := json.Unmarshal(w.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if _, exists := response["problems"]; !exists {
		t.Error("Expected 'problems' field in response")
	}

	if _, exists := response["query"]; !exists {
		t.Error("Expected 'query' field in response")
	}
}

func TestProblemHandlers_SearchProblems_EmptyQuery(t *testing.T) {
	router, handler := setupTestRouter()
	router.GET("/problems/search", handler.SearchProblems)

	req, _ := http.NewRequest("GET", "/problems/search", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestProblemHandlers_GetProblem_InvalidID(t *testing.T) {
	router, handler := setupTestRouter()
	router.GET("/problems/:id", handler.GetProblem)

	req, _ := http.NewRequest("GET", "/problems/invalid", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}