package services

import (
	"time"

	"github.com/google/uuid"
	"github.com/nambuitechx/nam-tcp/internal/models"
	"github.com/nambuitechx/nam-tcp/internal/repositories"
)

type TargetService struct {
	Repository *repositories.TargetRepository
}

func NewTargetService(repo *repositories.TargetRepository) *TargetService {
	return &TargetService{Repository: repo}
}

func (s *TargetService) GetTargets(limit int, offset int) ([]models.TargetModel, error) {
	return s.Repository.GetTargets(limit, offset)
}

func (s *TargetService) CreateTarget(payload *models.CreateTargetPayload) (*models.TargetModel, error) {
	newTarget := &models.TargetModel{
		ID: uuid.NewString(),
		Name: payload.Name,
		Host: payload.Host,
		Port: payload.Port,
		CreatedAt: int(time.Now().Unix()),
		UpdatedAt: int(time.Now().Unix()),
	}

	err := s.Repository.CreateTarget(newTarget)
	if err != nil {
		return nil, err
	}

	return newTarget, nil
}
