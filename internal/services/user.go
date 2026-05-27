package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/nambuitechx/nam-tcp/internal/models"
	"github.com/nambuitechx/nam-tcp/internal/repositories"
)

type UserService struct {
	Repository *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{Repository: repo}
}

func (s *UserService) GetUsers(limit int, offset int) ([]models.UserModel, error) {
	return s.Repository.GetUsers(limit, offset)
}

func (s *UserService) CreateUser(payload *models.CreateUserPayload) (*models.UserModel, error) {
	newUser := &models.UserModel{
		ID: uuid.NewString(),
		Email: payload.Email,
		Password: payload.Password,
		CreatedAt: int(time.Now().Unix()),
		UpdatedAt: int(time.Now().Unix()),
	}

	err := s.Repository.CreateUser(newUser)
	if err != nil {
		return nil, err
	}

	return newUser, nil
}
