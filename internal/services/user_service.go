package services

import (
	"errors"

	"github.com/realworld-api/internal/dtos"
	"github.com/realworld-api/internal/models"
	"github.com/realworld-api/internal/repositories"
	"github.com/realworld-api/internal/utils"
	"gorm.io/gorm"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{
		userRepo: repositories.NewUserRepository(db),
	}
}

func (s *UserService) Register(req *dtos.RegisterRequest) (*models.User, string, error) {
	if s.userRepo.ExistsByEmail(req.User.Email) {
		return nil, "", errors.New("email already exists")
	}

	if s.userRepo.ExistsByUsername(req.User.Username) {
		return nil, "", errors.New("username already exists")
	}

	hashedPassword, err := utils.HashPassword(req.User.Password)
	if err != nil {
		return nil, "", errors.New("failed to hash password")
	}

	user := &models.User{
		Username:     req.User.Username,
		Email:        req.User.Email,
		PasswordHash: hashedPassword,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, "", errors.New("failed to create user")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}

func (s *UserService) Login(req *dtos.LoginRequest) (*models.User, string, error) {
	user, err := s.userRepo.FindByEmail(req.User.Email)
	if err != nil {
		return nil, "", errors.New("invalid email or password")
	}

	if !utils.CheckPassword(req.User.Password, user.PasswordHash) {
		return nil, "", errors.New("invalid email or password")
	}

	token, err := utils.GenerateToken(user.ID, user.Username, user.Email)
	if err != nil {
		return nil, "", errors.New("failed to generate token")
	}

	return user, token, nil
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) Update(user *models.User, req *dtos.UpdateUserRequest) error {
	if req.User.Email != nil && *req.User.Email != user.Email {
		if s.userRepo.ExistsByEmail(*req.User.Email) {
			return errors.New("email already exists")
		}
		user.Email = *req.User.Email
	}

	if req.User.Username != nil && *req.User.Username != user.Username {
		if s.userRepo.ExistsByUsername(*req.User.Username) {
			return errors.New("username already exists")
		}
		user.Username = *req.User.Username
	}

	if req.User.Password != nil {
		hashedPassword, err := utils.HashPassword(*req.User.Password)
		if err != nil {
			return errors.New("failed to hash password")
		}
		user.PasswordHash = hashedPassword
	}

	if req.User.Bio != nil {
		user.Bio = *req.User.Bio
	}

	if req.User.Image != nil {
		user.Image = *req.User.Image
	}

	return s.userRepo.Update(user)
}

func (s *UserService) GenerateToken(user *models.User) (string, error) {
	return utils.GenerateToken(user.ID, user.Username, user.Email)
}
