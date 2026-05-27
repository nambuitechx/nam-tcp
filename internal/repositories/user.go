package repositories

import (
	"database/sql"

	"github.com/nambuitechx/nam-tcp/internal/models"
)

type UserRepository struct {
	DB *sql.DB 
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) GetUsers(limit int, offset int) ([]models.UserModel, error) {
	var rows *sql.Rows
	var err error
	var users []models.UserModel = []models.UserModel{}

	if limit == -1 {
		rows, err = r.DB.Query("SELECT id, email, password, created_at, updated_at FROM users")
	} else {
		rows, err = r.DB.Query("SELECT id, email, password, created_at, updated_at FROM users LIMIT ? OFFSET ?", limit, offset)
	}

	if err != nil {
		return users, err
	}
	defer rows.Close()

	for rows.Next() {
		var user models.UserModel
		rows.Scan(&user.ID, &user.Email, &user.Password, &user.CreatedAt, &user.UpdatedAt)
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *UserRepository) CreateUser(user *models.UserModel) error {
	_, err := r.DB.Exec(
		`INSERT INTO users(id, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, ?)`,
		user.ID, user.Email, user.Password, user.CreatedAt, user.UpdatedAt,
	)
	return err
}
