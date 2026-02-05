package services

import (
	"errors"

	"realworld-api/internal/dtos"
	"realworld-api/internal/repositories"

	"gorm.io/gorm"
)

type ProfileService struct {
	userRepo *repositories.UserRepository
}

func NewProfileService(db *gorm.DB) *ProfileService {
	return &ProfileService{
		userRepo: repositories.NewUserRepository(db),
	}
}

// GetProfile returns the profile of a user by username
func (s *ProfileService) GetProfile(username string, currentUserID uint) (*dtos.ProfileData, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	following := false
	if currentUserID > 0 {
		following = s.userRepo.IsFollowing(currentUserID, user.ID)
	}

	return &dtos.ProfileData{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: following,
	}, nil
}

// FollowUser follows a user
func (s *ProfileService) FollowUser(username string, currentUserID uint) (*dtos.ProfileData, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.ID == currentUserID {
		return nil, errors.New("cannot follow yourself")
	}

	if err := s.userRepo.Follow(currentUserID, user.ID); err != nil {
		return nil, errors.New("failed to follow user")
	}

	return &dtos.ProfileData{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: true,
	}, nil
}

// UnfollowUser unfollows a user
func (s *ProfileService) UnfollowUser(username string, currentUserID uint) (*dtos.ProfileData, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.ID == currentUserID {
		return nil, errors.New("cannot unfollow yourself")
	}

	if err := s.userRepo.Unfollow(currentUserID, user.ID); err != nil {
		return nil, errors.New("failed to unfollow user")
	}

	return &dtos.ProfileData{
		Username:  user.Username,
		Bio:       user.Bio,
		Image:     user.Image,
		Following: false,
	}, nil
}
