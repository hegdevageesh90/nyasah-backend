package api

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"nyasah/config"
	"nyasah/api/handlers"
	"nyasah/api/middleware"
)

type Server struct {
	router *gin.Engine
	db     *gorm.DB
	config *config.Config
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	server := &Server{
		router: gin.Default(),
		db:     db,
		config: cfg,
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Create handlers
	authHandler := handlers.NewAuthHandler(s.db, s.config)
	reviewHandler := handlers.NewReviewHandler(s.db)
	socialProofHandler := handlers.NewSocialProofHandler(s.db)

	// Public routes
	s.router.POST("/api/auth/register", authHandler.Register)
	s.router.POST("/api/auth/login", authHandler.Login)

	// Protected routes
	protected := s.router.Group("/api")
	protected.Use(middleware.AuthMiddleware(s.config.JWTSecret))
	{
		// Reviews
		protected.POST("/reviews", reviewHandler.Create)
		protected.GET("/reviews", reviewHandler.List)
		protected.GET("/reviews/:id", reviewHandler.Get)

		// Social Proof
		protected.POST("/social-proof", socialProofHandler.Create)
		protected.GET("/social-proof", socialProofHandler.List)
		protected.GET("/social-proof/analytics", socialProofHandler.GetAnalytics)
	}
}

func (s *Server) Start() error {
	return s.router.Run(":" + s.config.Port)
}