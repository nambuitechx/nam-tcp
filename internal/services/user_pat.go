package services

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
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

func (s *UserPATService) CreateUserPAT(payload *models.CreateUserPATPayload) (string, *models.UserPATModel, error) {
	plaintext, hash, err := GenerateNewUserPAT()
	if err != nil {
		return "", nil, err
	}

	now := int(time.Now().Unix())

	newUserPAT := &models.UserPATModel{
		ID: uuid.NewString(),
		UserID: payload.UserID,
		TargetID: payload.TargetID,
		HashToken: hash,
		CreatedAt: now,
		ExpiresAt: now + 24 * payload.TTLInHour * 3600,
		RevokedAt: 0,
	}

	err = s.Repository.CreateUserPAT(newUserPAT)
	if err != nil {
		return "", nil, err
	}

	return plaintext, newUserPAT, nil
}

func (s *UserPATService) RevokeUserPAT(id string) error {
	user_pat, err := s.Repository.GetUserPATByID(id)
	if err != nil {
		return err
	}

	if user_pat.RevokedAt != 0 {
		return errors.New("user pat already revoked")
	}

	err = s.Repository.RevokeUserPAT(id)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserPATService) RevokeExpiredUserPATs() error {
	return s.Repository.RevokeExpiredUserPATs()
}

func GenerateNewUserPAT() (string, string, error) {
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return "", "", err
	}

	plaintext := "nam_tcp_" + hex.EncodeToString(buf)
	sum := sha256.Sum256([]byte(plaintext))
	hash := hex.EncodeToString(sum[:])

	return plaintext, hash, nil
}
