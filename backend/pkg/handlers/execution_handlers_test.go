package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"leetcode-clone-backend/pkg/execution"
	"leetcode-clone-backend/pkg/models"

	"github.com/gin-gonic/gin"
)

// MockTestCaseRepository for testing
type MockTestCaseRepository struct {
	testCases []*models.TestCase
}

func (m *MockTestCaseRepository) GetByProblemID(problemID int) ([]*models.TestCase, error) {
	return m.testCases, nil
}

func (m *MockTestCaseRepository) GetPublicByProblemID(problemID int) ([]*models.TestCase, error) {
	var publicTestCases []*models.TestCase
	for _, tc := range m.testCases {
		if !tc.IsHidden {
			publicTestCases = append(publicTestCases, tc)
		}
	}
	return publicTestCases, nil
}

func (m *MockTestCaseRepository) Create(testCase *models.TestCase) (*models.TestCase, error) {
	return testCase, nil
}

func (m *MockTestCaseRepository) Update(testCase *models.TestCase) (*models.TestCase, error) {
	return testCase, nil
}

func (m *MockTestCaseRepository) Delete(id int) error {
	return nil
}

func (m *MockTestCaseRepository) DeleteByProblemID(problemID int) error {
	return nil
}

func (m *MockTestCaseRepository) GetByID(id int) (*models.TestCase, error) {
	return nil, nil
}

func TestExecutionHandlers_ValidateCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	executionService := execution.NewExecutionService()
	mockRepo := &MockTestCaseRepository{}
	handler := NewExecutionHandlers(executionService, mockRepo)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedValid  bool
	}{
		{
			name: "Valid JavaScript code",
			requestBody: map[string]interface{}{
				"code":     "function solution(input) { return input; }",
				"language": models.LanguageJavaScript,
			},
			expectedStatus: http.StatusOK,
			expectedValid:  true,
		},
		{
			name: "Invalid code with dangerous pattern",
			requestBody: map[string]interface{}{
				"code":     "import os; function solution(input) { return input; }",
				"language": models.LanguageJavaScript,
			},
			expectedStatus: http.StatusBadRequest,
			expectedValid:  false,
		},
		{
			name: "Missing required fields",
			requestBody: map[string]interface{}{
				"code": "function solution(input) { return input; }",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/validate", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call handler
			handler.ValidateCode(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Check response body for successful validations
			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				json.Unmarshal(w.Body.Bytes(), &response)

				if valid, ok := response["valid"].(bool); !ok || valid != tt.expectedValid {
					t.Errorf("Expected valid=%v, got %v", tt.expectedValid, valid)
				}
			}
		})
	}
}

func TestExecutionHandlers_GetSupportedLanguages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	executionService := execution.NewExecutionService()
	mockRepo := &MockTestCaseRepository{}
	handler := NewExecutionHandlers(executionService, mockRepo)

	// Create request
	req, _ := http.NewRequest("GET", "/languages", nil)
	w := httptest.NewRecorder()

	// Create Gin context
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Call handler
	handler.GetSupportedLanguages(c)

	// Check status code
	if w.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
	}

	// Check response body
	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	languages, ok := response["languages"].([]interface{})
	if !ok {
		t.Error("Expected languages array in response")
		return
	}

	// Should have 3 supported languages
	if len(languages) != 3 {
		t.Errorf("Expected 3 languages, got %d", len(languages))
	}

	// Check that each language has required fields
	for i, lang := range languages {
		langMap, ok := lang.(map[string]interface{})
		if !ok {
			t.Errorf("Language %d is not a map", i)
			continue
		}

		requiredFields := []string{"id", "name", "extension", "template"}
		for _, field := range requiredFields {
			if _, exists := langMap[field]; !exists {
				t.Errorf("Language %d missing required field: %s", i, field)
			}
		}
	}
}
func TestExecutionHandlers_RunCode(t *testing.T) {
	gin.SetMode(gin.TestMode)

	executionService := execution.NewExecutionService()

	// Mock test cases
	mockRepo := &MockTestCaseRepository{
		testCases: []*models.TestCase{
			{
				ID:             1,
				ProblemID:      1,
				Input:          "test input",
				ExpectedOutput: "test output",
				IsHidden:       false,
			},
		},
	}

	handler := NewExecutionHandlers(executionService, mockRepo)

	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Valid run request",
			requestBody: map[string]interface{}{
				"code":       "function solution(input) { return 'test output'; }",
				"language":   models.LanguageJavaScript,
				"problem_id": 1,
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Missing required fields",
			requestBody: map[string]interface{}{
				"code":     "function solution(input) { return input; }",
				"language": models.LanguageJavaScript,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			jsonBody, _ := json.Marshal(tt.requestBody)
			req, _ := http.NewRequest("POST", "/run", bytes.NewBuffer(jsonBody))
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Create Gin context
			c, _ := gin.CreateTestContext(w)
			c.Request = req

			// Call handler
			handler.RunCode(c)

			// Check status code
			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}
		})
	}
}
