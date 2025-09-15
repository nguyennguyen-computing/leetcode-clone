package main

import (
	"database/sql"
	"log"
	"os"

	"leetcode-clone-backend/pkg/auth"
	"leetcode-clone-backend/pkg/database"
	"leetcode-clone-backend/pkg/execution"
	"leetcode-clone-backend/pkg/handlers"
	"leetcode-clone-backend/pkg/middleware"
	"leetcode-clone-backend/pkg/repository"
	"leetcode-clone-backend/pkg/services"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	router            *gin.Engine
	db                *sql.DB
	repo              *repository.Repository
	authService       *auth.AuthService
	problemService    *services.ProblemService
	submissionService *services.SubmissionService
	executionService  *execution.ExecutionService
	authHandler       *handlers.AuthHandlers
	problemHandler    *handlers.ProblemHandlers
	submissionHandler *handlers.SubmissionHandlers
	executionHandler  *handlers.ExecutionHandlers
}

func main() {
	// Load database configuration
	dbConfig := database.LoadConfigFromEnv()

	// Connect to database
	db, err := database.Connect(dbConfig)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Run database migrations
	log.Println("Running database migrations...")
	if err := database.RunMigrations(db); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repository
	repo := repository.NewRepository(db)

	// Initialize authentication service
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key-change-in-production"
	}
	authService := auth.NewAuthService(jwtSecret)

	// Initialize services
	problemService := services.NewProblemService(repo.Problem, repo.TestCase)
	executionService := execution.NewExecutionService()
	submissionService := services.NewSubmissionService(repo.Submission, repo.TestCase, repo.UserProgress, executionService)

	// Initialize handlers
	authHandler := handlers.NewAuthHandlers(authService, repo.User)
	problemHandler := handlers.NewProblemHandlers(problemService)
	submissionHandler := handlers.NewSubmissionHandlers(submissionService)
	executionHandler := handlers.NewExecutionHandlers(executionService, repo.TestCase)

	server := &Server{
		router:            gin.Default(),
		db:                db,
		repo:              repo,
		authService:       authService,
		problemService:    problemService,
		submissionService: submissionService,
		executionService:  executionService,
		authHandler:       authHandler,
		problemHandler:    problemHandler,
		submissionHandler: submissionHandler,
		executionHandler:  executionHandler,
	}

	// Setup CORS
	server.router.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:4200",
			"http://127.0.0.1:4200",
		},
		AllowMethods: []string{
			"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH",
		},
		AllowHeaders: []string{
			"Origin", "Content-Type", "Authorization", "Accept",
			"X-Requested-With", "sec-ch-ua", "sec-ch-ua-mobile",
			"sec-ch-ua-platform", "User-Agent", "Referer",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // 12 hours
	}))

	// Setup routes
	server.setupRoutes()

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(server.router.Run(":" + port))
}

func (s *Server) setupRoutes() {
	api := s.router.Group("/api/v1")

	// Health check
	api.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes with rate limiting
	auth := api.Group("/auth")
	auth.Use(middleware.RateLimitMiddleware(10)) // 10 requests per minute for auth endpoints
	auth.POST("/register", s.authHandler.Register)
	auth.POST("/login", s.authHandler.Login)
	auth.POST("/password-reset", s.authHandler.RequestPasswordReset)
	auth.POST("/password-reset/confirm", s.authHandler.ResetPassword)

	// Public problem routes (read-only)
	api.GET("/problems", s.problemHandler.ListProblems)
	api.GET("/problems/search", s.problemHandler.SearchProblems)
	api.GET("/problems/:id", s.problemHandler.GetProblem)
	api.GET("/problems/slug/:slug", s.problemHandler.GetProblemBySlug)
	api.GET("/problems/:id/testcases", s.problemHandler.GetTestCases)

	// Protected routes
	protected := api.Group("/")
	protected.Use(handlers.AuthMiddleware(s.authService))

	// Admin-only problem management routes
	admin := protected.Group("/admin")
	admin.Use(s.adminMiddleware())
	admin.POST("/problems", s.problemHandler.CreateProblem)
	admin.PUT("/problems/:id", s.problemHandler.UpdateProblem)
	admin.DELETE("/problems/:id", s.problemHandler.DeleteProblem)
	admin.POST("/problems/:id/testcases", s.problemHandler.CreateTestCase)
	admin.PUT("/testcases/:id", s.problemHandler.UpdateTestCase)
	admin.DELETE("/testcases/:id", s.problemHandler.DeleteTestCase)

	// Code execution routes
	protected.POST("/execute/run", s.executionHandler.RunCode)
	protected.POST("/execute/submit", s.executionHandler.SubmitCode)
	protected.POST("/execute/validate", s.executionHandler.ValidateCode)
	protected.GET("/execute/languages", s.executionHandler.GetSupportedLanguages)

	// Submission routes
	protected.POST("/submissions", s.submissionHandler.CreateSubmission)
	protected.GET("/submissions/:id", s.submissionHandler.GetSubmission)
	protected.GET("/submissions/me", s.submissionHandler.GetUserSubmissions)
	protected.GET("/submissions/user/:userId", s.submissionHandler.GetUserSubmissions)
	protected.GET("/submissions/stats/me", s.submissionHandler.GetUserSubmissionStats)
	protected.GET("/submissions/stats/:userId", s.submissionHandler.GetUserSubmissionStats)

	// Problem-specific submission routes
	protected.GET("/problems/:problemId/submissions", s.submissionHandler.GetProblemSubmissions)
}

// adminMiddleware checks if the user is an admin
func (s *Server) adminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get user from context (set by auth middleware)
		userInterface, exists := c.Get("user")
		if !exists {
			c.JSON(401, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		user, ok := userInterface.(*auth.Claims)
		if !ok {
			c.JSON(401, gin.H{"error": "Invalid user context"})
			c.Abort()
			return
		}

		// Get full user details to check admin status
		fullUser, err := s.repo.User.GetByID(user.UserID)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to verify admin status"})
			c.Abort()
			return
		}

		if !fullUser.IsAdmin {
			c.JSON(403, gin.H{"error": "Admin access required"})
			c.Abort()
			return
		}

		c.Next()
	}
}
