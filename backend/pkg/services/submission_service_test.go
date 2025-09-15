package services

import (
	"testing"
	"time"

	"leetcode-clone-backend/pkg/execution"
	"leetcode-clone-backend/pkg/models"
	"leetcode-clone-backend/pkg/repository"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type MockSubmissionRepository struct {
	mock.Mock
}

func (m *MockSubmissionRepository) Create(submission *models.Submission) (*models.Submission, error) {
	args := m.Called(submission)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Submission), args.Error(1)
}

func (m *MockSubmissionRepository) GetByID(id int) (*models.Submission, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Submission), args.Error(1)
}

func (m *MockSubmissionRepository) GetByUserID(userID int, limit, offset int) ([]*models.Submission, error) {
	args := m.Called(userID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Submission), args.Error(1)
}

func (m *MockSubmissionRepository) GetByProblemID(problemID int, limit, offset int) ([]*models.Submission, error) {
	args := m.Called(problemID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Submission), args.Error(1)
}

func (m *MockSubmissionRepository) GetByUserAndProblem(userID, problemID int, limit, offset int) ([]*models.Submission, error) {
	args := m.Called(userID, problemID, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.Submission), args.Error(1)
}

func (m *MockSubmissionRepository) Update(submission *models.Submission) (*models.Submission, error) {
	args := m.Called(submission)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Submission), args.Error(1)
}

func (m *MockSubmissionRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockSubmissionRepository) GetLatestByUserAndProblem(userID, problemID int) (*models.Submission, error) {
	args := m.Called(userID, problemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Submission), args.Error(1)
}

type MockTestCaseRepository struct {
	mock.Mock
}

func (m *MockTestCaseRepository) Create(testCase *models.TestCase) (*models.TestCase, error) {
	args := m.Called(testCase)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TestCase), args.Error(1)
}

func (m *MockTestCaseRepository) GetByID(id int) (*models.TestCase, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TestCase), args.Error(1)
}

func (m *MockTestCaseRepository) GetByProblemID(problemID int) ([]*models.TestCase, error) {
	args := m.Called(problemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.TestCase), args.Error(1)
}

func (m *MockTestCaseRepository) GetPublicByProblemID(problemID int) ([]*models.TestCase, error) {
	args := m.Called(problemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.TestCase), args.Error(1)
}

func (m *MockTestCaseRepository) Update(testCase *models.TestCase) (*models.TestCase, error) {
	args := m.Called(testCase)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.TestCase), args.Error(1)
}

func (m *MockTestCaseRepository) Delete(id int) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockTestCaseRepository) DeleteByProblemID(problemID int) error {
	args := m.Called(problemID)
	return args.Error(0)
}

type MockUserProgressRepository struct {
	mock.Mock
}

func (m *MockUserProgressRepository) Create(progress *models.UserProgress) (*models.UserProgress, error) {
	args := m.Called(progress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserProgress), args.Error(1)
}

func (m *MockUserProgressRepository) GetByUserAndProblem(userID, problemID int) (*models.UserProgress, error) {
	args := m.Called(userID, problemID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserProgress), args.Error(1)
}

func (m *MockUserProgressRepository) GetByUserID(userID int) ([]*models.UserProgress, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*models.UserProgress), args.Error(1)
}

func (m *MockUserProgressRepository) Update(progress *models.UserProgress) (*models.UserProgress, error) {
	args := m.Called(progress)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.UserProgress), args.Error(1)
}

func (m *MockUserProgressRepository) Delete(userID, problemID int) error {
	args := m.Called(userID, problemID)
	return args.Error(0)
}

func (m *MockUserProgressRepository) GetSolvedCount(userID int) (int, error) {
	args := m.Called(userID)
	return args.Int(0), args.Error(1)
}

func (m *MockUserProgressRepository) GetSolvedCountByDifficulty(userID int) (map[string]int, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}

type MockExecutionService struct {
	mock.Mock
}

func (m *MockExecutionService) ExecuteCode(code, language string, testCases []models.TestCase) (*execution.ExecutionResult, error) {
	args := m.Called(code, language, testCases)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*execution.ExecutionResult), args.Error(1)
}

func (m *MockExecutionService) ValidateCode(code, language string) error {
	args := m.Called(code, language)
	return args.Error(0)
}

func TestSubmissionService_ProcessSubmission(t *testing.T) {
	// Setup mocks
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockTestCaseRepo := new(MockTestCaseRepository)
	mockUserProgressRepo := new(MockUserProgressRepository)
	mockExecutionService := new(MockExecutionService)

	service := NewSubmissionService(mockSubmissionRepo, mockTestCaseRepo, mockUserProgressRepo, mockExecutionService)

	t.Run("successful submission", func(t *testing.T) {
		// Setup test data
		req := &SubmissionRequest{
			UserID:    1,
			ProblemID: 1,
			Language:  models.LanguageJavaScript,
			Code:      "function solution(input) { return 'test'; }",
		}

		testCases := []*models.TestCase{
			{
				ID:             1,
				ProblemID:      1,
				Input:          "test input",
				ExpectedOutput: "test",
				IsHidden:       false,
			},
		}

		executionResult := &execution.ExecutionResult{
			Status:          models.StatusAccepted,
			RuntimeMs:       100,
			MemoryKb:        1024,
			TestCasesPassed: 1,
			TotalTestCases:  1,
			TestResults: []execution.TestResult{
				{
					Input:          "test input",
					ExpectedOutput: "test",
					ActualOutput:   "test",
					Passed:         true,
					RuntimeMs:      100,
					MemoryKb:       1024,
				},
			},
		}

		createdSubmission := &models.Submission{
			ID:              1,
			UserID:          1,
			ProblemID:       1,
			Language:        models.LanguageJavaScript,
			Code:            req.Code,
			Status:          models.StatusAccepted,
			RuntimeMs:       &executionResult.RuntimeMs,
			MemoryKb:        &executionResult.MemoryKb,
			TestCasesPassed: 1,
			TotalTestCases:  1,
			SubmittedAt:     time.Now(),
		}

		// Setup expectations
		mockTestCaseRepo.On("GetByProblemID", 1).Return(testCases, nil)
		mockExecutionService.On("ExecuteCode", req.Code, req.Language, mock.AnythingOfType("[]models.TestCase")).Return(executionResult, nil)
		mockSubmissionRepo.On("Create", mock.AnythingOfType("*models.Submission")).Return(createdSubmission, nil)
		mockUserProgressRepo.On("GetByUserAndProblem", 1, 1).Return(nil, repository.NewRepositoryError("GetByUserAndProblem", repository.ErrNotFound, "not_found"))
		mockUserProgressRepo.On("Create", mock.AnythingOfType("*models.UserProgress")).Return(&models.UserProgress{}, nil)

		// Execute
		result, err := service.ProcessSubmission(req)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, models.StatusAccepted, result.Status)
		assert.Equal(t, 1, result.TestCasesPassed)
		assert.Equal(t, 1, result.TotalTestCases)
		assert.Len(t, result.TestResults, 1)

		// Verify all expectations were met
		mockTestCaseRepo.AssertExpectations(t)
		mockExecutionService.AssertExpectations(t)
		mockSubmissionRepo.AssertExpectations(t)
		mockUserProgressRepo.AssertExpectations(t)
	})

	t.Run("invalid submission request", func(t *testing.T) {
		req := &SubmissionRequest{
			UserID:    0, // Invalid user ID
			ProblemID: 1,
			Language:  models.LanguageJavaScript,
			Code:      "function solution(input) { return 'test'; }",
		}

		result, err := service.ProcessSubmission(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid submission request")
	})

	t.Run("no test cases available", func(t *testing.T) {
		// Create new mocks for this test
		mockSubmissionRepo2 := new(MockSubmissionRepository)
		mockTestCaseRepo2 := new(MockTestCaseRepository)
		mockUserProgressRepo2 := new(MockUserProgressRepository)
		mockExecutionService2 := new(MockExecutionService)

		service2 := NewSubmissionService(mockSubmissionRepo2, mockTestCaseRepo2, mockUserProgressRepo2, mockExecutionService2)

		req := &SubmissionRequest{
			UserID:    1,
			ProblemID: 1,
			Language:  models.LanguageJavaScript,
			Code:      "function solution(input) { return 'test'; }",
		}

		mockTestCaseRepo2.On("GetByProblemID", 1).Return([]*models.TestCase{}, nil)

		result, err := service2.ProcessSubmission(req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no test cases available")

		mockTestCaseRepo2.AssertExpectations(t)
	})
}

func TestSubmissionService_GetSubmissionByID(t *testing.T) {
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockTestCaseRepo := new(MockTestCaseRepository)
	mockUserProgressRepo := new(MockUserProgressRepository)
	mockExecutionService := new(MockExecutionService)

	service := NewSubmissionService(mockSubmissionRepo, mockTestCaseRepo, mockUserProgressRepo, mockExecutionService)

	t.Run("successful retrieval", func(t *testing.T) {
		expectedSubmission := &models.Submission{
			ID:        1,
			UserID:    1,
			ProblemID: 1,
			Language:  models.LanguageJavaScript,
			Code:      "test code",
			Status:    models.StatusAccepted,
		}

		mockSubmissionRepo.On("GetByID", 1).Return(expectedSubmission, nil)

		result, err := service.GetSubmissionByID(1)

		assert.NoError(t, err)
		assert.Equal(t, expectedSubmission, result)

		mockSubmissionRepo.AssertExpectations(t)
	})

	t.Run("submission not found", func(t *testing.T) {
		mockSubmissionRepo.On("GetByID", 999).Return(nil, repository.NewRepositoryError("GetByID", repository.ErrNotFound, "not_found"))

		result, err := service.GetSubmissionByID(999)

		assert.Error(t, err)
		assert.Nil(t, result)

		mockSubmissionRepo.AssertExpectations(t)
	})
}

func TestSubmissionService_GetUserSubmissions(t *testing.T) {
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockTestCaseRepo := new(MockTestCaseRepository)
	mockUserProgressRepo := new(MockUserProgressRepository)
	mockExecutionService := new(MockExecutionService)

	service := NewSubmissionService(mockSubmissionRepo, mockTestCaseRepo, mockUserProgressRepo, mockExecutionService)

	t.Run("successful retrieval with pagination", func(t *testing.T) {
		submissions := []*models.Submission{
			{ID: 1, UserID: 1, Status: models.StatusAccepted},
			{ID: 2, UserID: 1, Status: models.StatusWrongAnswer},
		}

		mockSubmissionRepo.On("GetByUserID", 1, 21, 0).Return(submissions, nil) // 21 to check for next page

		result, err := service.GetUserSubmissions(1, 1, 20)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Len(t, result.Submissions, 2)
		assert.Equal(t, 1, result.Page)
		assert.Equal(t, 20, result.PageSize)
		assert.False(t, result.HasNext)

		mockSubmissionRepo.AssertExpectations(t)
	})
}

func TestSubmissionService_GetUserSubmissionStats(t *testing.T) {
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockTestCaseRepo := new(MockTestCaseRepository)
	mockUserProgressRepo := new(MockUserProgressRepository)
	mockExecutionService := new(MockExecutionService)

	service := NewSubmissionService(mockSubmissionRepo, mockTestCaseRepo, mockUserProgressRepo, mockExecutionService)

	t.Run("calculate stats correctly", func(t *testing.T) {
		runtime1 := 100
		runtime2 := 200
		memory1 := 1024
		memory2 := 2048

		submissions := []*models.Submission{
			{ID: 1, Status: models.StatusAccepted, RuntimeMs: &runtime1, MemoryKb: &memory1},
			{ID: 2, Status: models.StatusWrongAnswer, RuntimeMs: &runtime2, MemoryKb: &memory2},
			{ID: 3, Status: models.StatusAccepted, RuntimeMs: &runtime2, MemoryKb: &memory1},
		}

		mockSubmissionRepo.On("GetByUserID", 1, 1000, 0).Return(submissions, nil)

		stats, err := service.GetUserSubmissionStats(1)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 3, stats["total_submissions"])
		assert.Equal(t, 2, stats["accepted"])
		assert.Equal(t, 1, stats["wrong_answer"])
		assert.InDelta(t, 66.67, stats["acceptance_rate"], 0.01)
		assert.Equal(t, 166, stats["avg_runtime_ms"]) // (100+200+200)/3 = 166.67 rounded to 166
		assert.Equal(t, 1365, stats["avg_memory_kb"]) // (1024+2048+1024)/3 = 1365.33 rounded to 1365

		mockSubmissionRepo.AssertExpectations(t)
	})

	t.Run("empty submissions", func(t *testing.T) {
		mockSubmissionRepo2 := new(MockSubmissionRepository)
		mockTestCaseRepo2 := new(MockTestCaseRepository)
		mockUserProgressRepo2 := new(MockUserProgressRepository)
		mockExecutionService2 := new(MockExecutionService)

		service2 := NewSubmissionService(mockSubmissionRepo2, mockTestCaseRepo2, mockUserProgressRepo2, mockExecutionService2)

		mockSubmissionRepo2.On("GetByUserID", 1, 1000, 0).Return([]*models.Submission{}, nil)

		stats, err := service2.GetUserSubmissionStats(1)

		assert.NoError(t, err)
		assert.NotNil(t, stats)
		assert.Equal(t, 0, stats["total_submissions"])
		assert.Equal(t, 0.0, stats["acceptance_rate"])

		mockSubmissionRepo2.AssertExpectations(t)
	})
}

func TestSubmissionService_validateSubmissionRequest(t *testing.T) {
	mockSubmissionRepo := new(MockSubmissionRepository)
	mockTestCaseRepo := new(MockTestCaseRepository)
	mockUserProgressRepo := new(MockUserProgressRepository)
	mockExecutionService := new(MockExecutionService)

	service := NewSubmissionService(mockSubmissionRepo, mockTestCaseRepo, mockUserProgressRepo, mockExecutionService)

	tests := []struct {
		name    string
		req     *SubmissionRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: &SubmissionRequest{
				UserID:    1,
				ProblemID: 1,
				Language:  models.LanguageJavaScript,
				Code:      "function solution() {}",
			},
			wantErr: false,
		},
		{
			name: "invalid user ID",
			req: &SubmissionRequest{
				UserID:    0,
				ProblemID: 1,
				Language:  models.LanguageJavaScript,
				Code:      "function solution() {}",
			},
			wantErr: true,
			errMsg:  "invalid user ID",
		},
		{
			name: "invalid problem ID",
			req: &SubmissionRequest{
				UserID:    1,
				ProblemID: 0,
				Language:  models.LanguageJavaScript,
				Code:      "function solution() {}",
			},
			wantErr: true,
			errMsg:  "invalid problem ID",
		},
		{
			name: "empty code",
			req: &SubmissionRequest{
				UserID:    1,
				ProblemID: 1,
				Language:  models.LanguageJavaScript,
				Code:      "",
			},
			wantErr: true,
			errMsg:  "code cannot be empty",
		},
		{
			name: "unsupported language",
			req: &SubmissionRequest{
				UserID:    1,
				ProblemID: 1,
				Language:  "unsupported",
				Code:      "some code",
			},
			wantErr: true,
			errMsg:  "unsupported language",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.validateSubmissionRequest(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
