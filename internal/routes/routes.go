package routes

import (
	"realworld-api/internal/bootstrap"
	"realworld-api/internal/middleware"

	"github.com/gin-gonic/gin"
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

		// Article routes
		articles := api.Group("/articles")
		{
			// Feed endpoint (requires auth) - must be before /:slug
			articles.GET("/feed", middleware.AuthRequired(), appContainer.ArticleHandler.FeedArticles)

			// List articles (optional auth)
			articles.GET("", middleware.AuthOptional(), appContainer.ArticleHandler.ListArticles)

			// Get single article (no auth required)
			articles.GET("/:slug", appContainer.ArticleHandler.GetArticle)

			// Create article (requires auth)
			articles.POST("", middleware.AuthRequired(), appContainer.ArticleHandler.CreateArticle)

			// Update article (requires auth)
			articles.PUT("/:slug", middleware.AuthRequired(), appContainer.ArticleHandler.UpdateArticle)

			// Delete article (requires auth)
			articles.DELETE("/:slug", middleware.AuthRequired(), appContainer.ArticleHandler.DeleteArticle)

			// Favorite/unfavorite article (requires auth)
			articles.POST("/:slug/favorite", middleware.AuthRequired(), appContainer.ArticleHandler.FavoriteArticle)
			articles.DELETE("/:slug/favorite", middleware.AuthRequired(), appContainer.ArticleHandler.UnfavoriteArticle)

			// Comments (requires auth for POST/DELETE, optional for GET)
			articles.POST("/:slug/comments", middleware.AuthRequired(), appContainer.CommentHandler.AddComment)
			articles.GET("/:slug/comments", middleware.AuthOptional(), appContainer.CommentHandler.GetComments)
			articles.DELETE("/:slug/comments/:id", middleware.AuthRequired(), appContainer.CommentHandler.DeleteComment)
		}
	}
}
