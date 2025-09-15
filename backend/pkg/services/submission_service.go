package services

import (
	"fmt"
	"time"

	"leetcode-clone-backend/pkg/execution"
	"leetcode-clone-backend/pkg/models"
	"leetcode-clone-backend/pkg/repository"
)

// SubmissionServiceInterface defines the interface for submission service
type SubmissionServiceInterface interface {
	ProcessSubmission(req *SubmissionRequest) (*SubmissionResponse, error)
	GetSubmissionByID(id int) (*models.Submission, error)
	GetUserSubmissions(userID, page, pageSize int) (*SubmissionListResponse, error)
	GetProblemSubmissions(problemID, page, pageSize int) (*SubmissionListResponse, error)
	GetUserProblemSubmissions(userID, problemID, page, pageSize int) (*SubmissionListResponse, error)
	GetUserSubmissionStats(userID int) (map[string]interface{}, error)
}

// SubmissionService handles business logic for code submissions
type SubmissionService struct {
	submissionRepo   repository.SubmissionRepository
	testCaseRepo     repository.TestCaseRepository
	userProgressRepo repository.UserProgressRepository
	executionService execution.ExecutionServiceInterface
}

// NewSubmissionService creates a new submission service
func NewSubmissionService(
	submissionRepo repository.SubmissionRepository,
	testCaseRepo repository.TestCaseRepository,
	userProgressRepo repository.UserProgressRepository,
	executionService execution.ExecutionServiceInterface,
) *SubmissionService {
	return &SubmissionService{
		submissionRepo:   submissionRepo,
		testCaseRepo:     testCaseRepo,
		userProgressRepo: userProgressRepo,
		executionService: executionService,
	}
}

// SubmissionRequest represents a code submission request
type SubmissionRequest struct {
	UserID    int    `json:"user_id"`
	ProblemID int    `json:"problem_id"`
	Language  string `json:"language"`
	Code      string `json:"code"`
}

// SubmissionResponse represents the response after processing a submission
type SubmissionResponse struct {
	ID              int                    `json:"id"`
	Status          string                 `json:"status"`
	RuntimeMs       *int                   `json:"runtime_ms"`
	MemoryKb        *int                   `json:"memory_kb"`
	TestCasesPassed int                    `json:"test_cases_passed"`
	TotalTestCases  int                    `json:"total_test_cases"`
	ErrorMessage    *string                `json:"error_message"`
	SubmittedAt     time.Time              `json:"submitted_at"`
	TestResults     []execution.TestResult `json:"test_results,omitempty"`
}

// SubmissionListResponse represents a paginated list of submissions
type SubmissionListResponse struct {
	Submissions []*models.Submission `json:"submissions"`
	Total       int                  `json:"total"`
	Page        int                  `json:"page"`
	PageSize    int                  `json:"page_size"`
	HasNext     bool                 `json:"has_next"`
}

// ProcessSubmission processes a code submission by executing it and storing the result
func (ss *SubmissionService) ProcessSubmission(req *SubmissionRequest) (*SubmissionResponse, error) {
	// Validate the submission request
	if err := ss.validateSubmissionRequest(req); err != nil {
		return nil, fmt.Errorf("invalid submission request: %w", err)
	}

	// Get all test cases for the problem
	testCases, err := ss.testCaseRepo.GetByProblemID(req.ProblemID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve test cases: %w", err)
	}

	if len(testCases) == 0 {
		return nil, fmt.Errorf("no test cases available for problem %d", req.ProblemID)
	}

	// Convert []*models.TestCase to []models.TestCase for execution service
	allTestCases := make([]models.TestCase, len(testCases))
	for i, tc := range testCases {
		allTestCases[i] = *tc
	}

	// Execute the code against all test cases
	executionResult, err := ss.executionService.ExecuteCode(req.Code, req.Language, allTestCases)
	if err != nil {
		return nil, fmt.Errorf("code execution failed: %w", err)
	}

	// Create submission record
	submission := &models.Submission{
		UserID:          req.UserID,
		ProblemID:       req.ProblemID,
		Language:        req.Language,
		Code:            req.Code,
		Status:          executionResult.Status,
		TestCasesPassed: executionResult.TestCasesPassed,
		TotalTestCases:  executionResult.TotalTestCases,
	}

	// Set performance metrics if available
	if executionResult.RuntimeMs > 0 {
		submission.RuntimeMs = &executionResult.RuntimeMs
	}
	if executionResult.MemoryKb > 0 {
		submission.MemoryKb = &executionResult.MemoryKb
	}

	// Set error message if execution failed
	if executionResult.ErrorMessage != "" {
		submission.ErrorMessage = &executionResult.ErrorMessage
	}

	// Store the submission
	createdSubmission, err := ss.submissionRepo.Create(submission)
	if err != nil {
		return nil, fmt.Errorf("failed to store submission: %w", err)
	}

	// Update user progress if submission was accepted
	if executionResult.Status == models.StatusAccepted {
		if err := ss.updateUserProgress(req.UserID, req.ProblemID, createdSubmission.ID); err != nil {
			// Log error but don't fail the submission
			// In a production system, you might want to use a message queue for this
			fmt.Printf("Warning: failed to update user progress: %v\n", err)
		}
	}

	// Prepare response with filtered test results (only public test cases)
	response := &SubmissionResponse{
		ID:              createdSubmission.ID,
		Status:          createdSubmission.Status,
		RuntimeMs:       createdSubmission.RuntimeMs,
		MemoryKb:        createdSubmission.MemoryKb,
		TestCasesPassed: createdSubmission.TestCasesPassed,
		TotalTestCases:  createdSubmission.TotalTestCases,
		ErrorMessage:    createdSubmission.ErrorMessage,
		SubmittedAt:     createdSubmission.SubmittedAt,
		TestResults:     make([]execution.TestResult, 0),
	}

	// Only include public test case results in the response
	for i, testResult := range executionResult.TestResults {
		if i < len(testCases) && !testCases[i].IsHidden {
			response.TestResults = append(response.TestResults, testResult)
		}
	}

	return response, nil
}

// GetSubmissionByID retrieves a submission by its ID
func (ss *SubmissionService) GetSubmissionByID(id int) (*models.Submission, error) {
	submission, err := ss.submissionRepo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve submission: %w", err)
	}
	return submission, nil
}

// GetUserSubmissions retrieves submissions for a specific user with pagination
func (ss *SubmissionService) GetUserSubmissions(userID, page, pageSize int) (*SubmissionListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // Default page size
	}

	offset := (page - 1) * pageSize

	submissions, err := ss.submissionRepo.GetByUserID(userID, pageSize+1, offset) // Get one extra to check if there's a next page
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user submissions: %w", err)
	}

	hasNext := len(submissions) > pageSize
	if hasNext {
		submissions = submissions[:pageSize] // Remove the extra submission
	}

	return &SubmissionListResponse{
		Submissions: submissions,
		Total:       len(submissions), // Note: This is not the total count, just current page count
		Page:        page,
		PageSize:    pageSize,
		HasNext:     hasNext,
	}, nil
}

// GetProblemSubmissions retrieves submissions for a specific problem with pagination
func (ss *SubmissionService) GetProblemSubmissions(problemID, page, pageSize int) (*SubmissionListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // Default page size
	}

	offset := (page - 1) * pageSize

	submissions, err := ss.submissionRepo.GetByProblemID(problemID, pageSize+1, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve problem submissions: %w", err)
	}

	hasNext := len(submissions) > pageSize
	if hasNext {
		submissions = submissions[:pageSize]
	}

	return &SubmissionListResponse{
		Submissions: submissions,
		Total:       len(submissions),
		Page:        page,
		PageSize:    pageSize,
		HasNext:     hasNext,
	}, nil
}

// GetUserProblemSubmissions retrieves submissions for a specific user and problem with pagination
func (ss *SubmissionService) GetUserProblemSubmissions(userID, problemID, page, pageSize int) (*SubmissionListResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20 // Default page size
	}

	offset := (page - 1) * pageSize

	submissions, err := ss.submissionRepo.GetByUserAndProblem(userID, problemID, pageSize+1, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user problem submissions: %w", err)
	}

	hasNext := len(submissions) > pageSize
	if hasNext {
		submissions = submissions[:pageSize]
	}

	return &SubmissionListResponse{
		Submissions: submissions,
		Total:       len(submissions),
		Page:        page,
		PageSize:    pageSize,
		HasNext:     hasNext,
	}, nil
}

// GetUserSubmissionStats calculates submission statistics for a user
func (ss *SubmissionService) GetUserSubmissionStats(userID int) (map[string]interface{}, error) {
	// Get all user submissions (we'll need to implement a count method in repository for efficiency)
	submissions, err := ss.submissionRepo.GetByUserID(userID, 1000, 0) // Get a large number for stats
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve user submissions for stats: %w", err)
	}

	stats := map[string]interface{}{
		"total_submissions":     len(submissions),
		"accepted":              0,
		"wrong_answer":          0,
		"time_limit_exceeded":   0,
		"memory_limit_exceeded": 0,
		"runtime_error":         0,
		"compile_error":         0,
		"acceptance_rate":       0.0,
		"avg_runtime_ms":        0,
		"avg_memory_kb":         0,
	}

	if len(submissions) == 0 {
		return stats, nil
	}

	totalRuntime := 0
	totalMemory := 0
	runtimeCount := 0
	memoryCount := 0

	for _, submission := range submissions {
		switch submission.Status {
		case models.StatusAccepted:
			stats["accepted"] = stats["accepted"].(int) + 1
		case models.StatusWrongAnswer:
			stats["wrong_answer"] = stats["wrong_answer"].(int) + 1
		case models.StatusTimeLimitExceeded:
			stats["time_limit_exceeded"] = stats["time_limit_exceeded"].(int) + 1
		case models.StatusMemoryLimitExceeded:
			stats["memory_limit_exceeded"] = stats["memory_limit_exceeded"].(int) + 1
		case models.StatusRuntimeError:
			stats["runtime_error"] = stats["runtime_error"].(int) + 1
		case models.StatusCompileError:
			stats["compile_error"] = stats["compile_error"].(int) + 1
		}

		if submission.RuntimeMs != nil {
			totalRuntime += *submission.RuntimeMs
			runtimeCount++
		}
		if submission.MemoryKb != nil {
			totalMemory += *submission.MemoryKb
			memoryCount++
		}
	}

	// Calculate acceptance rate
	acceptedCount := stats["accepted"].(int)
	if len(submissions) > 0 {
		stats["acceptance_rate"] = float64(acceptedCount) / float64(len(submissions)) * 100
	}

	// Calculate average performance metrics
	if runtimeCount > 0 {
		stats["avg_runtime_ms"] = totalRuntime / runtimeCount
	}
	if memoryCount > 0 {
		stats["avg_memory_kb"] = totalMemory / memoryCount
	}

	return stats, nil
}

// validateSubmissionRequest validates the submission request
func (ss *SubmissionService) validateSubmissionRequest(req *SubmissionRequest) error {
	if req.UserID <= 0 {
		return fmt.Errorf("invalid user ID")
	}
	if req.ProblemID <= 0 {
		return fmt.Errorf("invalid problem ID")
	}
	if req.Code == "" {
		return fmt.Errorf("code cannot be empty")
	}
	if req.Language == "" {
		return fmt.Errorf("language cannot be empty")
	}

	// Validate language support
	supportedLanguages := []string{
		models.LanguageJavaScript,
		models.LanguagePython,
		models.LanguageJava,
	}

	languageSupported := false
	for _, lang := range supportedLanguages {
		if req.Language == lang {
			languageSupported = true
			break
		}
	}

	if !languageSupported {
		return fmt.Errorf("unsupported language: %s", req.Language)
	}

	return nil
}

// updateUserProgress updates the user's progress for a problem when they get an accepted submission
func (ss *SubmissionService) updateUserProgress(userID, problemID, submissionID int) error {
	// Get existing progress
	progress, err := ss.userProgressRepo.GetByUserAndProblem(userID, problemID)
	if err != nil {
		// If no progress exists, create new one
		if repository.IsNotFound(err) {
			newProgress := &models.UserProgress{
				UserID:           userID,
				ProblemID:        problemID,
				IsSolved:         true,
				BestSubmissionID: &submissionID,
				Attempts:         1,
				FirstSolvedAt:    &time.Time{},
			}
			now := time.Now()
			newProgress.FirstSolvedAt = &now

			_, err := ss.userProgressRepo.Create(newProgress)
			return err
		}
		return err
	}

	// Update existing progress
	progress.Attempts++
	if !progress.IsSolved {
		progress.IsSolved = true
		progress.BestSubmissionID = &submissionID
		now := time.Now()
		progress.FirstSolvedAt = &now
	} else {
		// If already solved, update best submission if this one is better
		// For now, we'll just update to the latest accepted submission
		// In a more sophisticated system, you'd compare performance metrics
		progress.BestSubmissionID = &submissionID
	}

	_, err = ss.userProgressRepo.Update(progress)
	return err
}
