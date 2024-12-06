package api

import (
	"nyasah-backend/api/handlers"
	"nyasah-backend/api/middleware"
	"nyasah-backend/config"
	"nyasah-backend/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Server struct {
	router    *gin.Engine
	db        *gorm.DB
	config    *config.Config
	aiService *services.Service
}

func NewServer(cfg *config.Config, db *gorm.DB) *Server {
	server := &Server{
		router:    gin.Default(),
		db:        db,
		config:    cfg,
		aiService: services.NewAIService(db),
	}
	server.setupRoutes()
	return server
}

func (s *Server) setupRoutes() {
	// Create handlers
	authHandler := handlers.NewAuthHandler(s.db, s.config)
	reviewHandler := handlers.NewReviewHandler(s.db)
	socialProofHandler := handlers.NewSocialProofHandler(s.db)
	aiQueryHandler := handlers.NewAIQueryHandler(s.db, s.aiService)
	insightsHandler := handlers.NewInsightsHandler(s.db, s.aiService)
	tenantHandler := handlers.NewTenantHandler(s.db)

	// Public routes
	s.router.POST("/api/auth/register", authHandler.Register)
	s.router.POST("/api/auth/login", authHandler.Login)

	// Protected routes
	protected := s.router.Group("/api")
	protected.Use(middleware.AuthMiddleware(s.config.JWTSecret))
	{
		// Tenants
		protected.POST("/admin/tenats", tenantHandler.Create)
		protected.GET("/admin/tenants/:id", tenantHandler.Get)
		protected.PUT("/admin/tenants/:id", tenantHandler.Update)

		// Reviews
		protected.POST("/reviews", reviewHandler.Create)
		protected.GET("/reviews", reviewHandler.List)
		protected.GET("/reviews/:id", reviewHandler.Get)

		// Social Proof
		protected.POST("/social-proof", socialProofHandler.Create)
		protected.GET("/social-proof", socialProofHandler.List)
		protected.GET("/social-proof/analytics", socialProofHandler.GetAnalytics)

		// AI Features
		protected.POST("/ai/query", aiQueryHandler.Query)
		protected.GET("/ai/insights/product/:id", insightsHandler.GetProductInsights)
		protected.GET("/ai/insights/recommendations", insightsHandler.GetRecommendations)
		protected.GET("/ai/insights/trends", insightsHandler.GetTrendAnalysis)
	}
}

func (s *Server) Start() error {
	return s.router.Run(":" + s.config.Port)
}
