package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/realworld-api/internal/bootstrap"
	"github.com/realworld-api/internal/middleware"
)

func SetupRoutes(router *gin.Engine, appContainer *bootstrap.AppContainer) {
	// Add CORS middleware
	router.Use(middleware.CORS())

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "RealWorld API is running",
		})
	})

	// API v1 routes
	api := router.Group("/api")
	{
		// User routes (authentication)
		users := api.Group("/users")
		{
			users.POST("", appContainer.UserHandler.Register)    // Register
			users.POST("/login", appContainer.UserHandler.Login) // Login
		}

		// Current user routes (requires auth middleware)
		user := api.Group("/user")
		user.Use(middleware.AuthRequired())
		{
			user.GET("", appContainer.UserHandler.GetCurrentUser) // Get current user
			user.PUT("", appContainer.UserHandler.UpdateUser)     // Update current user
		}
	}
}
