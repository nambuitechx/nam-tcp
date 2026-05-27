package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/nambuitechx/nam-tcp/internal/models"
	"github.com/nambuitechx/nam-tcp/internal/repositories"
)

type UserPATService struct {
	Repository *repositories.UserPATRepository
}

func NewUserPATService(repo *repositories.UserPATRepository) *UserPATService {
	return &UserPATService{Repository: repo}
}

func (s *UserPATService) GetUserPATs(limit int, offset int) ([]models.UserPATModel, error) {
	return s.Repository.GetUserPATs(limit, offset)
}

func (s *UserPATService) CreateUserPAT(payload *models.CreateUserPATPayload) (*models.UserPATModel, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return nil, err
	}

	plaintext := "nam_tcp_" + hex.EncodeToString(buf)
	sum := sha256.Sum256([]byte(plaintext))
	hash := hex.EncodeToString(sum[:])
	now := int(time.Now().Unix())

	newUserPAT := &models.UserPATModel{
		ID: uuid.NewString(),
		UserID: payload.UserID,
		HashToken: hash,
		CreatedAt: now,
		ExpiresAt: now + 24 * payload.TTLInHour * 3600,
		RevokedAt: 0,
		TargetHost: payload.TargetHost,
		TargetPort: payload.TargetPort,
	}

	err := s.Repository.CreateUserPAT(newUserPAT)
	if err != nil {
		return nil, err
	}

	return newUserPAT, nil
}
