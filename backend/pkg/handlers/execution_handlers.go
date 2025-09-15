package handlers

import (
	"net/http"

	"leetcode-clone-backend/pkg/execution"
	"leetcode-clone-backend/pkg/models"
	"leetcode-clone-backend/pkg/repository"

	"github.com/gin-gonic/gin"
)

// ExecutionHandlers handles code execution related HTTP requests
type ExecutionHandlers struct {
	executionService *execution.ExecutionService
	testCaseRepo     repository.TestCaseRepository
}

// NewExecutionHandlers creates a new execution handlers instance
func NewExecutionHandlers(executionService *execution.ExecutionService, testCaseRepo repository.TestCaseRepository) *ExecutionHandlers {
	return &ExecutionHandlers{
		executionService: executionService,
		testCaseRepo:     testCaseRepo,
	}
}

// ExecuteCodeRequest represents the request payload for code execution
type ExecuteCodeRequest struct {
	Code      string `json:"code" binding:"required"`
	Language  string `json:"language" binding:"required"`
	ProblemID int    `json:"problem_id" binding:"required"`
}

// RunCode executes code against public test cases (for testing during development)
func (eh *ExecutionHandlers) RunCode(c *gin.Context) {
	var req ExecuteCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Get public test cases for the problem
	testCases, err := eh.testCaseRepo.GetByProblemID(req.ProblemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve test cases"})
		return
	}

	// Filter to only public test cases for run (not submit)
	publicTestCases := make([]models.TestCase, 0)
	for _, tc := range testCases {
		if !tc.IsHidden {
			publicTestCases = append(publicTestCases, *tc)
		}
	}

	if len(publicTestCases) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No public test cases available for this problem"})
		return
	}

	// Execute code
	result, err := eh.executionService.ExecuteCode(req.Code, req.Language, publicTestCases)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Code execution failed"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// SubmitCode executes code against all test cases (including hidden ones) for submission
func (eh *ExecutionHandlers) SubmitCode(c *gin.Context) {
	var req ExecuteCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Get all test cases for the problem (including hidden ones)
	testCases, err := eh.testCaseRepo.GetByProblemID(req.ProblemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve test cases"})
		return
	}

	if len(testCases) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No test cases available for this problem"})
		return
	}

	// Convert []*models.TestCase to []models.TestCase
	allTestCases := make([]models.TestCase, len(testCases))
	for i, tc := range testCases {
		allTestCases[i] = *tc
	}

	// Execute code against all test cases
	result, err := eh.executionService.ExecuteCode(req.Code, req.Language, allTestCases)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Code execution failed"})
		return
	}

	// For submissions, we don't return detailed test results for hidden test cases
	// Only return the overall result and public test case details
	submissionResult := &execution.ExecutionResult{
		Status:          result.Status,
		Output:          result.Output,
		ErrorMessage:    result.ErrorMessage,
		RuntimeMs:       result.RuntimeMs,
		MemoryKb:        result.MemoryKb,
		TestCasesPassed: result.TestCasesPassed,
		TotalTestCases:  result.TotalTestCases,
		TestResults:     make([]execution.TestResult, 0),
	}

	// Only include public test case results in the response
	for i, testResult := range result.TestResults {
		if i < len(testCases) && !testCases[i].IsHidden {
			submissionResult.TestResults = append(submissionResult.TestResults, testResult)
		}
	}

	c.JSON(http.StatusOK, submissionResult)
}

// ValidateCode validates code without executing it (for syntax checking)
func (eh *ExecutionHandlers) ValidateCode(c *gin.Context) {
	var req struct {
		Code     string `json:"code" binding:"required"`
		Language string `json:"language" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Create a temporary execution service instance for validation
	tempService := execution.NewExecutionService()

	// Validate the code (this checks for dangerous patterns and length)
	if err := tempService.ValidateCode(req.Code, req.Language); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"valid": false,
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid": true,
	})
}

// GetSupportedLanguages returns the list of supported programming languages
func (eh *ExecutionHandlers) GetSupportedLanguages(c *gin.Context) {
	languages := []gin.H{
		{
			"id":        models.LanguageJavaScript,
			"name":      "JavaScript",
			"extension": ".js",
			"template":  "function solution(input) {\n    // Your code here\n    return \"\";\n}",
		},
		{
			"id":        models.LanguagePython,
			"name":      "Python",
			"extension": ".py",
			"template":  "def solution(input_data):\n    # Your code here\n    return \"\"",
		},
		{
			"id":        models.LanguageJava,
			"name":      "Java",
			"extension": ".java",
			"template":  "public String solution(String input) {\n    // Your code here\n    return \"\";\n}",
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"languages": languages,
	})
}
