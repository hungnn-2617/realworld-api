package bootstrap

import (
	"realworld-api/internal/config"
	"realworld-api/internal/handlers"
	"realworld-api/internal/services"

	"gorm.io/gorm"
)

type AppContainer struct {
	DB              *gorm.DB
	UserService     *services.UserService
	UserHandler     *handlers.UserHandler
	ProfileService  *services.ProfileService
	ProfileHandler  *handlers.ProfileHandler
	ArticleService  *services.ArticleService
	CommentService  *services.CommentService
	TagService      *services.TagService
	ArticleHandler  *handlers.ArticleHandler
	CommentHandler  *handlers.CommentHandler
	TagHandler      *handlers.TagHandler
	FavoriteHandler *handlers.FavoriteHandler
}

func NewAppContainer() *AppContainer {
	db := config.GetDB()

	// Initialize services
	userService := services.NewUserService(db)
	profileService := services.NewProfileService(db)
	articleService := services.NewArticleService(db)
	commentService := services.NewCommentService(db)
	tagService := services.NewTagService(db)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	profileHandler := handlers.NewProfileHandler(profileService)
	articleHandler := handlers.NewArticleHandler(articleService)
	commentHandler := handlers.NewCommentHandler(commentService)
	tagHandler := handlers.NewTagHandler(tagService)
	favoriteHandler := handlers.NewFavoriteHandler(articleService)

	return &AppContainer{
		DB:              db,
		UserService:     userService,
		UserHandler:     userHandler,
		ProfileService:  profileService,
		ProfileHandler:  profileHandler,
		ArticleService:  articleService,
		CommentService:  commentService,
		TagService:      tagService,
		ArticleHandler:  articleHandler,
		CommentHandler:  commentHandler,
		TagHandler:      tagHandler,
		FavoriteHandler: favoriteHandler,
	}
}
