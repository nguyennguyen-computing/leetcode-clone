package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"leetcode-clone-backend/pkg/models"
	"leetcode-clone-backend/pkg/repository"
	"leetcode-clone-backend/pkg/services"
)

// ProblemHandlers handles HTTP requests for problems
type ProblemHandlers struct {
	problemService *services.ProblemService
}

// NewProblemHandlers creates a new problem handlers instance
func NewProblemHandlers(problemService *services.ProblemService) *ProblemHandlers {
	return &ProblemHandlers{
		problemService: problemService,
	}
}

// CreateProblem handles POST /problems
func (h *ProblemHandlers) CreateProblem(c *gin.Context) {
	var problem models.Problem
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	created, err := h.problemService.CreateProblem(&problem)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "slug_exists") {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Problem with this slug already exists",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create problem",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetProblem handles GET /problems/:id
func (h *ProblemHandlers) GetProblem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid problem ID",
			"details": "Problem ID must be a valid integer",
		})
		return
	}

	problem, err := h.problemService.GetProblem(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Problem not found",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get problem",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, problem)
}

// GetProblemBySlug handles GET /problems/slug/:slug
func (h *ProblemHandlers) GetProblemBySlug(c *gin.Context) {
	slug := c.Param("slug")
	if strings.TrimSpace(slug) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid slug",
			"details": "Slug cannot be empty",
		})
		return
	}

	problem, err := h.problemService.GetProblemBySlug(slug)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Problem not found",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get problem",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, problem)
}

// UpdateProblem handles PUT /problems/:id
func (h *ProblemHandlers) UpdateProblem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid problem ID",
			"details": "Problem ID must be a valid integer",
		})
		return
	}

	var problem models.Problem
	if err := c.ShouldBindJSON(&problem); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Set the ID from the URL parameter
	problem.ID = id

	updated, err := h.problemService.UpdateProblem(&problem)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Problem not found",
				"details": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "slug_exists") {
			c.JSON(http.StatusConflict, gin.H{
				"error":   "Problem with this slug already exists",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update problem",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteProblem handles DELETE /problems/:id
func (h *ProblemHandlers) DeleteProblem(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid problem ID",
			"details": "Problem ID must be a valid integer",
		})
		return
	}

	err = h.problemService.DeleteProblem(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Problem not found",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete problem",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListProblems handles GET /problems
func (h *ProblemHandlers) ListProblems(c *gin.Context) {
	filters := repository.ProblemFilters{}

	// Parse difficulty filters
	if difficultyStr := c.Query("difficulty"); difficultyStr != "" {
		filters.Difficulty = strings.Split(difficultyStr, ",")
	}

	// Parse tag filters
	if tagsStr := c.Query("tags"); tagsStr != "" {
		filters.Tags = strings.Split(tagsStr, ",")
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	// Parse offset
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	// Parse sort by
	if sortBy := c.Query("sort_by"); sortBy != "" {
		filters.SortBy = sortBy
	}

	// Parse sort order
	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		filters.SortOrder = sortOrder
	}

	problems, err := h.problemService.ListProblems(filters)
	if err != nil {
		if strings.Contains(err.Error(), "invalid filters") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid filters",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to list problems",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"problems": problems,
		"count":    len(problems),
	})
}

// SearchProblems handles GET /problems/search
func (h *ProblemHandlers) SearchProblems(c *gin.Context) {
	query := c.Query("q")
	if strings.TrimSpace(query) == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Search query is required",
			"details": "Query parameter 'q' cannot be empty",
		})
		return
	}

	filters := repository.ProblemFilters{}

	// Parse difficulty filters
	if difficultyStr := c.Query("difficulty"); difficultyStr != "" {
		filters.Difficulty = strings.Split(difficultyStr, ",")
	}

	// Parse tag filters
	if tagsStr := c.Query("tags"); tagsStr != "" {
		filters.Tags = strings.Split(tagsStr, ",")
	}

	// Parse limit
	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			filters.Limit = limit
		}
	}

	// Parse offset
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			filters.Offset = offset
		}
	}

	// Parse sort by
	if sortBy := c.Query("sort_by"); sortBy != "" {
		filters.SortBy = sortBy
	}

	// Parse sort order
	if sortOrder := c.Query("sort_order"); sortOrder != "" {
		filters.SortOrder = sortOrder
	}

	problems, err := h.problemService.SearchProblems(query, filters)
	if err != nil {
		if strings.Contains(err.Error(), "invalid filters") || strings.Contains(err.Error(), "cannot be empty") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Invalid search parameters",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to search problems",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"problems": problems,
		"count":    len(problems),
		"query":    query,
	})
}

// CreateTestCase handles POST /problems/:id/testcases
func (h *ProblemHandlers) CreateTestCase(c *gin.Context) {
	idStr := c.Param("id")
	problemID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid problem ID",
			"details": "Problem ID must be a valid integer",
		})
		return
	}

	var testCase models.TestCase
	if err := c.ShouldBindJSON(&testCase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Set the problem ID from the URL parameter
	testCase.ProblemID = problemID

	created, err := h.problemService.CreateTestCase(&testCase)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Problem not found",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create test case",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, created)
}

// GetTestCases handles GET /problems/:id/testcases
func (h *ProblemHandlers) GetTestCases(c *gin.Context) {
	idStr := c.Param("id")
	problemID, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid problem ID",
			"details": "Problem ID must be a valid integer",
		})
		return
	}

	// Check if only public test cases are requested
	publicOnly := c.Query("public") == "true"

	var testCases []*models.TestCase
	if publicOnly {
		testCases, err = h.problemService.GetPublicTestCases(problemID)
	} else {
		testCases, err = h.problemService.GetTestCases(problemID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get test cases",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"test_cases": testCases,
		"count":      len(testCases),
	})
}

// UpdateTestCase handles PUT /testcases/:id
func (h *ProblemHandlers) UpdateTestCase(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid test case ID",
			"details": "Test case ID must be a valid integer",
		})
		return
	}

	var testCase models.TestCase
	if err := c.ShouldBindJSON(&testCase); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	// Set the ID from the URL parameter
	testCase.ID = id

	updated, err := h.problemService.UpdateTestCase(&testCase)
	if err != nil {
		if strings.Contains(err.Error(), "validation failed") {
			c.JSON(http.StatusBadRequest, gin.H{
				"error":   "Validation failed",
				"details": err.Error(),
			})
			return
		}
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Test case not found",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to update test case",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, updated)
}

// DeleteTestCase handles DELETE /testcases/:id
func (h *ProblemHandlers) DeleteTestCase(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid test case ID",
			"details": "Test case ID must be a valid integer",
		})
		return
	}

	err = h.problemService.DeleteTestCase(id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			c.JSON(http.StatusNotFound, gin.H{
				"error":   "Test case not found",
				"details": err.Error(),
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to delete test case",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusNoContent, nil)
}