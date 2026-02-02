package bootstrap

import (
	"github.com/realworld-api/internal/config"
	"github.com/realworld-api/internal/handlers"
	"github.com/realworld-api/internal/services"
	"gorm.io/gorm"
)

type AppContainer struct {
	DB          *gorm.DB
	UserService *services.UserService
	UserHandler *handlers.UserHandler
}

func NewAppContainer() *AppContainer {
	db := config.GetDB()

	// Initialize services
	userService := services.NewUserService(db)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)

	return &AppContainer{
		DB:          db,
		UserService: userService,
		UserHandler: userHandler,
	}
}
