package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/realworld-api/internal/config"
	"github.com/realworld-api/internal/middleware"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()

	if cfg.App.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"message": "RealWorld API is running",
		})
	})

	api := r.Group("/api")
	{
		api.GET("/", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Welcome to RealWorld API",
			})
		})

		// Authentication routes (no auth required)
		users := api.Group("/users")
		{
			users.POST("", placeholderHandler("Register user"))    // POST /api/users
			users.POST("/login", placeholderHandler("Login user")) // POST /api/users/login
		}

		user := api.Group("/user")
		user.Use(middleware.AuthRequired())
		{
			user.GET("", placeholderHandler("Get current user"))    // GET /api/user
			user.PUT("", placeholderHandler("Update current user")) // PUT /api/user
		}

		profiles := api.Group("/profiles")
		{
			profiles.GET("/:username", middleware.AuthOptional(), placeholderHandler("Get profile"))
			profiles.POST("/:username/follow", middleware.AuthRequired(), placeholderHandler("Follow user"))
			profiles.DELETE("/:username/follow", middleware.AuthRequired(), placeholderHandler("Unfollow user"))
		}

		// Article routes
		articles := api.Group("/articles")
		{
			// Public routes (no auth  or optional⚪)
			articles.GET("", middleware.AuthOptional(), placeholderHandler("List articles"))
			articles.GET("/:slug", placeholderHandler("Get article"))
			articles.GET("/:slug/comments", middleware.AuthOptional(), placeholderHandler("Get comments"))

			articles.GET("/feed", middleware.AuthRequired(), placeholderHandler("Get feed"))
			articles.POST("", middleware.AuthRequired(), placeholderHandler("Create article"))
			articles.PUT("/:slug", middleware.AuthRequired(), placeholderHandler("Update article"))
			articles.DELETE("/:slug", middleware.AuthRequired(), placeholderHandler("Delete article"))
			articles.POST("/:slug/comments", middleware.AuthRequired(), placeholderHandler("Add comment"))
			articles.DELETE("/:slug/comments/:id", middleware.AuthRequired(), placeholderHandler("Delete comment"))
			articles.POST("/:slug/favorite", middleware.AuthRequired(), placeholderHandler("Favorite article"))
			articles.DELETE("/:slug/favorite", middleware.AuthRequired(), placeholderHandler("Unfavorite article"))
		}

		api.GET("/tags", placeholderHandler("Get tags"))
	}

	log.Printf("Server starting on port %s...", cfg.App.Port)
	if err := r.Run(":" + cfg.App.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// placeholderHandler returns a placeholder handler for unimplemented endpoints
func placeholderHandler(description string) gin.HandlerFunc {
	return func(c *gin.Context) {
		response := gin.H{
			"message": description + " - Not implemented yet",
		}

		// If user is authenticated, include user info
		if userID, exists := c.Get(middleware.ContextUserID); exists {
			response["authenticated"] = true
			response["user_id"] = userID
		}

		c.JSON(200, response)
	}
}
