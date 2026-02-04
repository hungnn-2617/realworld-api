package bootstrap

import (
	"realworld-api/internal/config"
	"realworld-api/internal/handlers"
	"realworld-api/internal/services"

	"gorm.io/gorm"
)

type AppContainer struct {
	DB             *gorm.DB
	UserService    *services.UserService
	UserHandler    *handlers.UserHandler
	ArticleService *services.ArticleService
	CommentService *services.CommentService
	ArticleHandler *handlers.ArticleHandler
	CommentHandler *handlers.CommentHandler
}

func NewAppContainer() *AppContainer {
	db := config.GetDB()

	// Initialize services
	userService := services.NewUserService(db)
	articleService := services.NewArticleService(db)
	commentService := services.NewCommentService(db)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	articleHandler := handlers.NewArticleHandler(articleService)
	commentHandler := handlers.NewCommentHandler(commentService)

	return &AppContainer{
		DB:             db,
		UserService:    userService,
		UserHandler:    userHandler,
		ArticleService: articleService,
		CommentService: commentService,
		ArticleHandler: articleHandler,
		CommentHandler: commentHandler,
	}
}
