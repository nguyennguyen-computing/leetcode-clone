package handlers

import (
	"net/http"
	"strconv"

	"leetcode-clone-backend/pkg/auth"
	"leetcode-clone-backend/pkg/repository"
	"leetcode-clone-backend/pkg/services"

	"github.com/gin-gonic/gin"
)

// SubmissionHandlers handles submission-related HTTP requests
type SubmissionHandlers struct {
	submissionService services.SubmissionServiceInterface
}

// NewSubmissionHandlers creates a new submission handlers instance
func NewSubmissionHandlers(submissionService services.SubmissionServiceInterface) *SubmissionHandlers {
	return &SubmissionHandlers{
		submissionService: submissionService,
	}
}

// SubmitCodeRequest represents the request payload for code submission
type SubmitCodeRequest struct {
	ProblemID int    `json:"problem_id" binding:"required"`
	Language  string `json:"language" binding:"required"`
	Code      string `json:"code" binding:"required"`
}

// CreateSubmission handles POST /api/v1/submissions
func (sh *SubmissionHandlers) CreateSubmission(c *gin.Context) {
	// Get user from context (set by auth middleware)
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := userInterface.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	var req SubmitCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Create submission request
	submissionReq := &services.SubmissionRequest{
		UserID:    user.UserID,
		ProblemID: req.ProblemID,
		Language:  req.Language,
		Code:      req.Code,
	}

	// Process the submission
	result, err := sh.submissionService.ProcessSubmission(submissionReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, result)
}

// GetSubmission handles GET /api/v1/submissions/:id
func (sh *SubmissionHandlers) GetSubmission(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid submission ID"})
		return
	}

	submission, err := sh.submissionService.GetSubmissionByID(id)
	if err != nil {
		if repository.IsNotFound(err) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Submission not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve submission"})
		return
	}

	// Check if user has permission to view this submission
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := userInterface.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	// Users can only view their own submissions (unless they're admin)
	if submission.UserID != user.UserID {
		// TODO: Check if user is admin - for now, deny access
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
		return
	}

	c.JSON(http.StatusOK, submission)
}

// GetUserSubmissions handles GET /api/v1/submissions/user/:userId or GET /api/v1/submissions/me
func (sh *SubmissionHandlers) GetUserSubmissions(c *gin.Context) {
	// Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := userInterface.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	// Determine which user's submissions to retrieve
	var targetUserID int
	userIDParam := c.Param("userId")

	if userIDParam == "me" || userIDParam == "" {
		// Get current user's submissions
		targetUserID = user.UserID
	} else {
		// Get specific user's submissions
		var err error
		targetUserID, err = strconv.Atoi(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Users can only view their own submissions (unless they're admin)
		if targetUserID != user.UserID {
			// TODO: Check if user is admin - for now, deny access
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	// Parse problem filter if provided
	problemIDStr := c.Query("problem_id")
	if problemIDStr != "" {
		problemID, err := strconv.Atoi(problemIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
			return
		}

		// Get submissions for specific user and problem
		result, err := sh.submissionService.GetUserProblemSubmissions(targetUserID, problemID, page, pageSize)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve submissions"})
			return
		}

		c.JSON(http.StatusOK, result)
		return
	}

	// Get all submissions for the user
	result, err := sh.submissionService.GetUserSubmissions(targetUserID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve submissions"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetProblemSubmissions handles GET /api/v1/problems/:problemId/submissions
func (sh *SubmissionHandlers) GetProblemSubmissions(c *gin.Context) {
	problemIDStr := c.Param("problemId")
	problemID, err := strconv.Atoi(problemIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid problem ID"})
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	result, err := sh.submissionService.GetProblemSubmissions(problemID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve submissions"})
		return
	}

	c.JSON(http.StatusOK, result)
}

// GetUserSubmissionStats handles GET /api/v1/submissions/stats/me or GET /api/v1/submissions/stats/:userId
func (sh *SubmissionHandlers) GetUserSubmissionStats(c *gin.Context) {
	// Get user from context
	userInterface, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user, ok := userInterface.(*auth.Claims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user context"})
		return
	}

	// Determine which user's stats to retrieve
	var targetUserID int
	userIDParam := c.Param("userId")

	if userIDParam == "me" || userIDParam == "" {
		// Get current user's stats
		targetUserID = user.UserID
	} else {
		// Get specific user's stats
		var err error
		targetUserID, err = strconv.Atoi(userIDParam)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		// Users can only view their own stats (unless they're admin)
		if targetUserID != user.UserID {
			// TODO: Check if user is admin - for now, deny access
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	stats, err := sh.submissionService.GetUserSubmissionStats(targetUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve submission stats"})
		return
	}

	c.JSON(http.StatusOK, stats)
}
