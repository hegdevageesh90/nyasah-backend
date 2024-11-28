package api

import (
	"nyasah/api/handlers"
	"nyasah/api/middleware"
	"nyasah/config"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
	// Admin routes for tenant management
	admin := s.router.Group("/api/admin")
	admin.Use(middleware.AuthMiddleware(s.config.JWTSecret))
	{
		tenantHandler := handlers.NewTenantHandler(s.db)
		admin.POST("/tenants", tenantHandler.Create)
		admin.GET("/tenants/:id", tenantHandler.Get)
		admin.PUT("/tenants/:id", tenantHandler.Update)
	}

	// Tenant-specific routes
	api := s.router.Group("/api")
	api.Use(middleware.TenantMiddleware(s.db))
	{
		// Auth routes
		authHandler := handlers.NewAuthHandler(s.db, s.config)
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)

		// Protected routes
		protected := api.Group("")
		protected.Use(middleware.AuthMiddleware(s.config.JWTSecret))
		{
			reviewHandler := handlers.NewReviewHandler(s.db)
			protected.POST("/reviews", reviewHandler.Create)
			protected.GET("/reviews", reviewHandler.List)
			protected.GET("/reviews/:id", reviewHandler.Get)

			socialProofHandler := handlers.NewSocialProofHandler(s.db)
			protected.POST("/social-proof", socialProofHandler.Create)
			protected.GET("/social-proof", socialProofHandler.List)
			protected.GET("/social-proof/analytics", socialProofHandler.GetAnalytics)
		}
	}
}

func (s *Server) Start() error {
	return s.router.Run(":" + s.config.Port)
}
