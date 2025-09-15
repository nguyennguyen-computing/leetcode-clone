package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"leetcode-clone-backend/pkg/auth"
	"leetcode-clone-backend/pkg/models"
	"leetcode-clone-backend/pkg/repository"
	"leetcode-clone-backend/pkg/services"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock submission service
type MockSubmissionService struct {
	mock.Mock
}

func (m *MockSubmissionService) ProcessSubmission(req *services.SubmissionRequest) (*services.SubmissionResponse, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.SubmissionResponse), args.Error(1)
}

func (m *MockSubmissionService) GetSubmissionByID(id int) (*models.Submission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Submission), args.Error(1)
}

func (m *MockSubmissionService) GetUserSubmissions(userID, page, pageSize int) (*services.SubmissionListResponse, error) {
	args := m.Called(userID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.SubmissionListResponse), args.Error(1)
}

func (m *MockSubmissionService) GetProblemSubmissions(problemID, page, pageSize int) (*services.SubmissionListResponse, error) {
	args := m.Called(problemID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.SubmissionListResponse), args.Error(1)
}

func (m *MockSubmissionService) GetUserProblemSubmissions(userID, problemID, page, pageSize int) (*services.SubmissionListResponse, error) {
	args := m.Called(userID, problemID, page, pageSize)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*services.SubmissionListResponse), args.Error(1)
}

func (m *MockSubmissionService) GetUserSubmissionStats(userID int) (map[string]interface{}, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]interface{}), args.Error(1)
}

func setupSubmissionTestRouter() (*gin.Engine, *MockSubmissionService) {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	mockService := new(MockSubmissionService)
	handler := NewSubmissionHandlers(mockService)

	// Add auth middleware mock
	router.Use(func(c *gin.Context) {
		// Set a mock user in context
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
	api.GET("/submissions/user/:userId", handler.GetUserSubmissions)
	api.GET("/submissions/stats/me", handler.GetUserSubmissionStats)
	api.GET("/submissions/stats/:userId", handler.GetUserSubmissionStats)
	api.GET("/problems/:problemId/submissions", handler.GetProblemSubmissions)

	return router, mockService
}

func TestSubmissionHandlers_CreateSubmission(t *testing.T) {
	router, mockService := setupSubmissionTestRouter()

	t.Run("successful submission", func(t *testing.T) {
		requestBody := SubmitCodeRequest{
			ProblemID: 1,
			Language:  models.LanguageJavaScript,
			Code:      "function solution(input) { return 'test'; }",
		}

		expectedResponse := &services.SubmissionResponse{
			ID:              1,
			Status:          models.StatusAccepted,
			TestCasesPassed: 1,
			TotalTestCases:  1,
			SubmittedAt:     time.Now(),
		}

		mockService.On("ProcessSubmission", mock.AnythingOfType("*services.SubmissionRequest")).Return(expectedResponse, nil)

		body, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/api/v1/submissions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)

		var response services.SubmissionResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedResponse.ID, response.ID)
		assert.Equal(t, expectedResponse.Status, response.Status)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid request payload", func(t *testing.T) {
		requestBody := map[string]interface{}{
			"problem_id": "invalid", // Should be int
			"language":   models.LanguageJavaScript,
			"code":       "function solution() {}",
		}

		body, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/api/v1/submissions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)

		var response map[string]string
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Invalid request payload", response["error"])
	})

	t.Run("service error", func(t *testing.T) {
		router2, mockService2 := setupSubmissionTestRouter()

		requestBody := SubmitCodeRequest{
			ProblemID: 1,
			Language:  models.LanguageJavaScript,
			Code:      "function solution() {}",
		}

		mockService2.On("ProcessSubmission", mock.AnythingOfType("*services.SubmissionRequest")).Return(nil, errors.New("service error"))

		body, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "/api/v1/submissions", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		router2.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockService2.AssertExpectations(t)
	})
}

func TestSubmissionHandlers_GetSubmission(t *testing.T) {
	router, mockService := setupSubmissionTestRouter()

	t.Run("successful retrieval", func(t *testing.T) {
		expectedSubmission := &models.Submission{
			ID:        1,
			UserID:    1, // Same as mock user
			ProblemID: 1,
			Language:  models.LanguageJavaScript,
			Code:      "test code",
			Status:    models.StatusAccepted,
		}

		mockService.On("GetSubmissionByID", 1).Return(expectedSubmission, nil)

		req, _ := http.NewRequest("GET", "/api/v1/submissions/1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Submission
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedSubmission.ID, response.ID)

		mockService.AssertExpectations(t)
	})

	t.Run("submission not found", func(t *testing.T) {
		mockService.On("GetSubmissionByID", 999).Return(nil, repository.NewRepositoryError("GetByID", repository.ErrNotFound, "not_found"))

		req, _ := http.NewRequest("GET", "/api/v1/submissions/999", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("access denied - different user", func(t *testing.T) {
		router2, mockService2 := setupSubmissionTestRouter()

		expectedSubmission := &models.Submission{
			ID:        1,
			UserID:    2, // Different from mock user (1)
			ProblemID: 1,
			Language:  models.LanguageJavaScript,
			Code:      "test code",
			Status:    models.StatusAccepted,
		}

		mockService2.On("GetSubmissionByID", 1).Return(expectedSubmission, nil)

		req, _ := http.NewRequest("GET", "/api/v1/submissions/1", nil)
		w := httptest.NewRecorder()
		router2.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)

		mockService2.AssertExpectations(t)
	})

	t.Run("invalid submission ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/submissions/invalid", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestSubmissionHandlers_GetUserSubmissions(t *testing.T) {
	router, mockService := setupSubmissionTestRouter()

	t.Run("get current user submissions", func(t *testing.T) {
		expectedResponse := &services.SubmissionListResponse{
			Submissions: []*models.Submission{
				{ID: 1, UserID: 1, Status: models.StatusAccepted},
				{ID: 2, UserID: 1, Status: models.StatusWrongAnswer},
			},
			Total:    2,
			Page:     1,
			PageSize: 20,
			HasNext:  false,
		}

		mockService.On("GetUserSubmissions", 1, 1, 20).Return(expectedResponse, nil)

		req, _ := http.NewRequest("GET", "/api/v1/submissions/me", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response services.SubmissionListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Submissions, 2)

		mockService.AssertExpectations(t)
	})

	t.Run("get user submissions with problem filter", func(t *testing.T) {
		expectedResponse := &services.SubmissionListResponse{
			Submissions: []*models.Submission{
				{ID: 1, UserID: 1, ProblemID: 1, Status: models.StatusAccepted},
			},
			Total:    1,
			Page:     1,
			PageSize: 20,
			HasNext:  false,
		}

		mockService.On("GetUserProblemSubmissions", 1, 1, 1, 20).Return(expectedResponse, nil)

		req, _ := http.NewRequest("GET", "/api/v1/submissions/me?problem_id=1", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		mockService.AssertExpectations(t)
	})

	t.Run("access denied - different user", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/submissions/user/2", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusForbidden, w.Code)
	})
}

func TestSubmissionHandlers_GetUserSubmissionStats(t *testing.T) {
	router, mockService := setupSubmissionTestRouter()

	t.Run("get current user stats", func(t *testing.T) {
		expectedStats := map[string]interface{}{
			"total_submissions": 10,
			"accepted":          7,
			"wrong_answer":      3,
			"acceptance_rate":   70.0,
		}

		mockService.On("GetUserSubmissionStats", 1).Return(expectedStats, nil)

		req, _ := http.NewRequest("GET", "/api/v1/submissions/stats/me", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]interface{}
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(10), response["total_submissions"])
		assert.Equal(t, float64(70), response["acceptance_rate"])

		mockService.AssertExpectations(t)
	})

	t.Run("service error", func(t *testing.T) {
		router2, mockService2 := setupSubmissionTestRouter()

		mockService2.On("GetUserSubmissionStats", 1).Return(nil, errors.New("service error"))

		req, _ := http.NewRequest("GET", "/api/v1/submissions/stats/me", nil)
		w := httptest.NewRecorder()
		router2.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)

		mockService2.AssertExpectations(t)
	})
}

func TestSubmissionHandlers_GetProblemSubmissions(t *testing.T) {
	router, mockService := setupSubmissionTestRouter()

	t.Run("successful retrieval", func(t *testing.T) {
		expectedResponse := &services.SubmissionListResponse{
			Submissions: []*models.Submission{
				{ID: 1, UserID: 1, ProblemID: 1, Status: models.StatusAccepted},
				{ID: 2, UserID: 2, ProblemID: 1, Status: models.StatusWrongAnswer},
			},
			Total:    2,
			Page:     1,
			PageSize: 20,
			HasNext:  false,
		}

		mockService.On("GetProblemSubmissions", 1, 1, 20).Return(expectedResponse, nil)

		req, _ := http.NewRequest("GET", "/api/v1/problems/1/submissions", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response services.SubmissionListResponse
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Len(t, response.Submissions, 2)

		mockService.AssertExpectations(t)
	})

	t.Run("invalid problem ID", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/problems/invalid/submissions", nil)
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
