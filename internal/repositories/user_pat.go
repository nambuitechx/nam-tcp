package repositories

import (
	"database/sql"

	"github.com/nambuitechx/nam-tcp/internal/models"
)

type UserPATRepository struct {
	DB *sql.DB 
}

func NewUserPATRepository(db *sql.DB) *UserPATRepository {
	return &UserPATRepository{DB: db}
}

func (r *UserPATRepository) GetUserPATs(limit int, offset int) ([]models.UserPATModel, error) {
	rows, err := r.DB.Query(
		`SELECT id, user_id, hash_token, created_at, expires_at, revoked_at, target_host, target_port FROM user_pats LIMIT ? OFFSET ?`,
		limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var user_pats []models.UserPATModel
	for rows.Next() {
		var user_pat models.UserPATModel
		rows.Scan(&user_pat.ID, &user_pat.UserID, &user_pat.HashToken, &user_pat.CreatedAt, &user_pat.ExpiresAt, &user_pat.RevokedAt, &user_pat.TargetHost, &user_pat.TargetPort)
		user_pats = append(user_pats, user_pat)
	}
	return user_pats, rows.Err()
}

func (r *UserPATRepository) CreateUserPAT(user_pat *models.UserPATModel) error {
	_, err := r.DB.Exec(
		`INSERT INTO user_pats(id, user_id, hash_token, created_at, expires_at, revoked_at, target_host, target_port) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		user_pat.ID, user_pat.UserID, user_pat.HashToken, user_pat.CreatedAt, user_pat.ExpiresAt, user_pat.RevokedAt, user_pat.TargetHost, user_pat.TargetPort,
	)
	return err
}
