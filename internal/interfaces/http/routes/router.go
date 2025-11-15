package routes

import (
	"meobeo-talk-api/internal/application"
	"meobeo-talk-api/internal/config"
	"meobeo-talk-api/internal/infrastructure/persistence/postgres"
	"meobeo-talk-api/internal/interfaces/http/handler"
	"meobeo-talk-api/internal/interfaces/http/middleware"
	"meobeo-talk-api/pkg/logger"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
)

func SetupRouter(cfg *config.Config, db *sqlx.DB, log *logger.Logger) *gin.Engine {
	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// Global middleware
	router.Use(middleware.Recovery(log))
	router.Use(middleware.Logger(log))
	router.Use(middleware.CORS())

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)

	// Initialize services
	userService := application.NewUserService(userRepo)

	// Initialize handlers
	userHandler := handler.NewUserHandler(userService, log)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		users := v1.Group("/users")
		{
			users.GET("", userHandler.GetUsers)
			users.GET("/:id", userHandler.GetUserByID)
			users.POST("", userHandler.CreateUser)
		}
	}

	return router
}
