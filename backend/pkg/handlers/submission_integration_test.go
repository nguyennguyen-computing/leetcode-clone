package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"leetcode-clone-backend/pkg/auth"
	"leetcode-clone-backend/pkg/execution"
	"leetcode-clone-backend/pkg/models"
	"leetcode-clone-backend/pkg/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// Integration test to verify the submission flow works end-to-end
func TestSubmissionIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Create a mock submission service that simulates real behavior
	mockService := &MockSubmissionServiceIntegration{}
	handler := NewSubmissionHandlers(mockService)

	router := gin.New()

	// Add auth middleware mock
	router.Use(func(c *gin.Context) {
		user := &auth.Claims{
			UserID:   1,
			Username: "testuser",
		}
		c.Set("user", user)
		c.Next()
	})

	// Setup routes
	api := router.Group("/api/v1")
	api.POST("/submissions", handler.CreateSubmission)
	api.GET("/submissions/:id", handler.GetSubmission)
	api.GET("/submissions/me", handler.GetUserSubmissions)

	t.Run("complete submission workflow", func(t *testing.T) {
		// Step 1: Create a submission
		requestBody := SubmitCodeRequest{
			ProblemID: 1,
			Language:  models.LanguageJavaScript,
			Code:      "function solution(input) { return input.trim(); }",
		}

		body, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/api/v1/submissions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var createResponse services.SubmissionResponse
		err := json.Unmarshal(w.Body.Bytes(), &createResponse)
		assert.NoError(t, err)
		assert.Equal(t, models.StatusAccepted, createResponse.Status)
		assert.Equal(t, 1, createResponse.ID)

		// Step 2: Retrieve the submission by ID
		req2, _ := http.NewRequest("GET", "/api/v1/submissions/1", nil)
		w2 := httptest.NewRecorder()
		router.ServeHTTP(w2, req2)

		assert.Equal(t, http.StatusOK, w2.Code)

		var getResponse models.Submission
		err = json.Unmarshal(w2.Body.Bytes(), &getResponse)
		assert.NoError(t, err)
		assert.Equal(t, 1, getResponse.ID)
		assert.Equal(t, 1, getResponse.UserID)

		// Step 3: Get user submissions
		req3, _ := http.NewRequest("GET", "/api/v1/submissions/me", nil)
		w3 := httptest.NewRecorder()
		router.ServeHTTP(w3, req3)

		assert.Equal(t, http.StatusOK, w3.Code)

		var listResponse services.SubmissionListResponse
		err = json.Unmarshal(w3.Body.Bytes(), &listResponse)
		assert.NoError(t, err)
		assert.Len(t, listResponse.Submissions, 1)
		assert.Equal(t, 1, listResponse.Submissions[0].ID)
	})
}

// MockSubmissionServiceIntegration provides a more realistic mock for integration testing
type MockSubmissionServiceIntegration struct{}

func (m *MockSubmissionServiceIntegration) ProcessSubmission(req *services.SubmissionRequest) (*services.SubmissionResponse, error) {
	return &services.SubmissionResponse{
		ID:              1,
		Status:          models.StatusAccepted,
		RuntimeMs:       &[]int{150}[0],
		MemoryKb:        &[]int{1024}[0],
		TestCasesPassed: 2,
		TotalTestCases:  2,
		SubmittedAt:     time.Now(),
		TestResults: []execution.TestResult{
			{
				Input:          "test",
				ExpectedOutput: "test",
				ActualOutput:   "test",
				Passed:         true,
				RuntimeMs:      75,
				MemoryKb:       512,
			},
		},
	}, nil
}

func (m *MockSubmissionServiceIntegration) GetSubmissionByID(id int) (*models.Submission, error) {
	runtime := 150
	memory := 1024
	return &models.Submission{
		ID:              id,
		UserID:          1,
		ProblemID:       1,
		Language:        models.LanguageJavaScript,
		Code:            "function solution(input) { return input.trim(); }",
		Status:          models.StatusAccepted,
		RuntimeMs:       &runtime,
		MemoryKb:        &memory,
		TestCasesPassed: 2,
		TotalTestCases:  2,
	}, nil
}

func (m *MockSubmissionServiceIntegration) GetUserSubmissions(userID, page, pageSize int) (*services.SubmissionListResponse, error) {
	runtime := 150
	memory := 1024
	submissions := []*models.Submission{
		{
			ID:              1,
			UserID:          userID,
			ProblemID:       1,
			Language:        models.LanguageJavaScript,
			Code:            "function solution(input) { return input.trim(); }",
			Status:          models.StatusAccepted,
			RuntimeMs:       &runtime,
			MemoryKb:        &memory,
			TestCasesPassed: 2,
			TotalTestCases:  2,
		},
	}

	return &services.SubmissionListResponse{
		Submissions: submissions,
		Total:       1,
		Page:        page,
		PageSize:    pageSize,
		HasNext:     false,
	}, nil
}

func (m *MockSubmissionServiceIntegration) GetProblemSubmissions(problemID, page, pageSize int) (*services.SubmissionListResponse, error) {
	return &services.SubmissionListResponse{
		Submissions: []*models.Submission{},
		Total:       0,
		Page:        page,
		PageSize:    pageSize,
		HasNext:     false,
	}, nil
}

func (m *MockSubmissionServiceIntegration) GetUserProblemSubmissions(userID, problemID, page, pageSize int) (*services.SubmissionListResponse, error) {
	return &services.SubmissionListResponse{
		Submissions: []*models.Submission{},
		Total:       0,
		Page:        page,
		PageSize:    pageSize,
		HasNext:     false,
	}, nil
}

func (m *MockSubmissionServiceIntegration) GetUserSubmissionStats(userID int) (map[string]interface{}, error) {
	return map[string]interface{}{
		"total_submissions": 1,
		"accepted":          1,
		"wrong_answer":      0,
		"acceptance_rate":   100.0,
		"avg_runtime_ms":    150,
		"avg_memory_kb":     1024,
	}, nil
}
