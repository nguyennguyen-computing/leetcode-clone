package handlers

import (
	"database/sql"
	"net/http"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
	"leetcode-clone-backend/pkg/auth"
	"leetcode-clone-backend/pkg/models"
	"leetcode-clone-backend/pkg/repository"
)

// AuthHandlers handles authentication-related HTTP requests
type AuthHandlers struct {
	authService *auth.AuthService
	userRepo    repository.UserRepository
}

// NewAuthHandlers creates a new AuthHandlers instance
func NewAuthHandlers(authService *auth.AuthService, userRepo repository.UserRepository) *AuthHandlers {
	return &AuthHandlers{
		authService: authService,
		userRepo:    userRepo,
	}
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=3,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// AuthResponse represents the authentication response
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// PasswordResetRequest represents the password reset request payload
type PasswordResetRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// PasswordResetConfirmRequest represents the password reset confirmation payload
type PasswordResetConfirmRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// Register handles user registration
func (h *AuthHandlers) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Validate username format (alphanumeric and underscores only)
	if !isValidUsername(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid username",
			"message": "Username can only contain letters, numbers, and underscores",
		})
		return
	}

	// Check if user already exists
	existingUser, err := h.userRepo.GetByEmail(req.Email)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"message": "Failed to check existing user",
		})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "User already exists",
			"message": "A user with this email already exists",
		})
		return
	}

	// Check if username is taken
	existingUserByUsername, err := h.userRepo.GetByUsername(req.Username)
	if err != nil && err != sql.ErrNoRows {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"message": "Failed to check existing username",
		})
		return
	}
	if existingUserByUsername != nil {
		c.JSON(http.StatusConflict, gin.H{
			"error":   "Username taken",
			"message": "This username is already taken",
		})
		return
	}

	// Hash password
	hashedPassword, err := h.authService.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid password",
			"message": err.Error(),
		})
		return
	}

	// Create user
	user := &models.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
		IsAdmin:      false,
	}

	createdUser, err := h.userRepo.Create(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Registration failed",
			"message": "Failed to create user account",
		})
		return
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(createdUser.ID, createdUser.Username, createdUser.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Token generation failed",
			"message": "Failed to generate authentication token",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusCreated, AuthResponse{
		Token: token,
		User:  *createdUser,
	})
}

// Login handles user login
func (h *AuthHandlers) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid credentials",
				"message": "Email or password is incorrect",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Database error",
			"message": "Failed to authenticate user",
		})
		return
	}

	// Verify password
	if err := h.authService.VerifyPassword(user.PasswordHash, req.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error":   "Invalid credentials",
			"message": "Email or password is incorrect",
		})
		return
	}

	// Generate JWT token
	token, err := h.authService.GenerateToken(user.ID, user.Username, user.IsAdmin)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Token generation failed",
			"message": "Failed to generate authentication token",
		})
		return
	}

	// Return success response
	c.JSON(http.StatusOK, AuthResponse{
		Token: token,
		User:  *user,
	})
}

// RequestPasswordReset handles password reset requests
func (h *AuthHandlers) RequestPasswordReset(c *gin.Context) {
	var req PasswordResetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Get user by email
	user, err := h.userRepo.GetByEmail(req.Email)
	if err != nil {
		// Don't reveal if email exists or not for security
		c.JSON(http.StatusOK, gin.H{
			"message": "If the email exists, a password reset link has been sent",
		})
		return
	}

	// Generate password reset token
	resetToken, err := h.authService.GeneratePasswordResetToken(user.ID, user.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Token generation failed",
			"message": "Failed to generate password reset token",
		})
		return
	}

	// TODO: In a real application, you would send this token via email
	// For now, we'll return it in the response (NOT recommended for production)
	c.JSON(http.StatusOK, gin.H{
		"message":     "Password reset token generated",
		"reset_token": resetToken, // Remove this in production
	})
}

// ResetPassword handles password reset confirmation
func (h *AuthHandlers) ResetPassword(c *gin.Context) {
	var req PasswordResetConfirmRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request",
			"message": err.Error(),
		})
		return
	}

	// Validate reset token
	userID, email, err := h.authService.ValidatePasswordResetToken(req.Token)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid token",
			"message": "Password reset token is invalid or expired",
		})
		return
	}

	// Get user to verify email matches
	user, err := h.userRepo.GetByID(userID)
	if err != nil || user.Email != email {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid token",
			"message": "Password reset token is invalid",
		})
		return
	}

	// Hash new password
	hashedPassword, err := h.authService.HashPassword(req.NewPassword)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid password",
			"message": err.Error(),
		})
		return
	}

	// Update user password
	user.PasswordHash = hashedPassword
	_, err = h.userRepo.Update(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Update failed",
			"message": "Failed to update password",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Password reset successfully",
	})
}

// isValidUsername checks if username contains only alphanumeric characters and underscores
func isValidUsername(username string) bool {
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	return matched
}

// AuthMiddleware creates a JWT authentication middleware
func AuthMiddleware(authService *auth.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Missing authorization header",
				"message": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Check if header starts with "Bearer "
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid authorization header",
				"message": "Authorization header must be in format 'Bearer <token>'",
			})
			c.Abort()
			return
		}

		token := tokenParts[1]
		claims, err := authService.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error":   "Invalid token",
				"message": "Authentication token is invalid or expired",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("is_admin", claims.IsAdmin)

		c.Next()
	}
}