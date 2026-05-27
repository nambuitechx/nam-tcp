package repositories

import (
	"database/sql"
	"time"

	"github.com/nambuitechx/nam-tcp/internal/models"
)

type UserPATRepository struct {
	DB *sql.DB 
}

func NewUserPATRepository(db *sql.DB) *UserPATRepository {
	return &UserPATRepository{DB: db}
}

func (r *UserPATRepository) GetUserPATs(limit int, offset int) ([]models.UserPATModel, error) {
	var rows *sql.Rows
	var err error
	var user_pats []models.UserPATModel = []models.UserPATModel{}

	if limit == -1 {
		rows, err = r.DB.Query(
			`SELECT id, user_id, target_id, hash_token, created_at, expires_at, revoked_at FROM user_pats`,
		)
	} else {
		rows, err = r.DB.Query(
			`SELECT id, user_id, target_id, hash_token, created_at, expires_at, revoked_at FROM user_pats LIMIT ? OFFSET ?`,
			limit, offset,
		)
	}

	if err != nil {
		return user_pats, err
	}
	defer rows.Close()

	for rows.Next() {
		var user_pat models.UserPATModel
		rows.Scan(&user_pat.ID, &user_pat.UserID, &user_pat.TargetID, &user_pat.HashToken, &user_pat.CreatedAt, &user_pat.ExpiresAt, &user_pat.RevokedAt)
		user_pats = append(user_pats, user_pat)
	}
	return user_pats, rows.Err()
}

func (r *UserPATRepository) GetUserPATByID(id string) (*models.UserPATModel, error) {
	row := r.DB.QueryRow(
		`SELECT id, user_id, target_id, hash_token, created_at, expires_at, revoked_at FROM user_pats WHERE id = ?`,
		id,
	)
	var user_pat models.UserPATModel
	err := row.Scan(&user_pat.ID, &user_pat.UserID, &user_pat.TargetID, &user_pat.HashToken, &user_pat.CreatedAt, &user_pat.ExpiresAt, &user_pat.RevokedAt)
	return &user_pat, err
}

func (r *UserPATRepository) CreateUserPAT(user_pat *models.UserPATModel) error {
	_, err := r.DB.Exec(
		`INSERT INTO user_pats(id, user_id, target_id, hash_token, created_at, expires_at, revoked_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		user_pat.ID, user_pat.UserID, user_pat.TargetID, user_pat.HashToken, user_pat.CreatedAt, user_pat.ExpiresAt, user_pat.RevokedAt,
	)
	return err
}

func (r *UserPATRepository) RevokeUserPAT(id string) error {
	_, err := r.DB.Exec(
		`UPDATE user_pats SET revoked_at = ? WHERE id = ?`,
		int(time.Now().Unix()), id,
	)
	return err
}

func (r *UserPATRepository) GetUserPATByHashToken(hashToken string) (*models.UserPATModel, error) {
	row := r.DB.QueryRow(
		`SELECT id, user_id, target_id, hash_token, created_at, expires_at, revoked_at FROM user_pats WHERE hash_token = ?`,
		hashToken,
	)
	var user_pat models.UserPATModel
	err := row.Scan(&user_pat.ID, &user_pat.UserID, &user_pat.TargetID, &user_pat.HashToken, &user_pat.CreatedAt, &user_pat.ExpiresAt, &user_pat.RevokedAt)
	return &user_pat, err
}

func (r *UserPATRepository) RevokeExpiredUserPATs() error {
	_, err := r.DB.Exec(
		`UPDATE user_pats SET revoked_at = ? WHERE expires_at < ? AND revoked_at = 0`,
		int(time.Now().Unix()), int(time.Now().Unix()),
	)
	return err
}
