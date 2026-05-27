package repositories

import (
	"database/sql"

	"github.com/nambuitechx/nam-tcp/internal/models"
)

type TargetRepository struct {
	DB *sql.DB 
}

func NewTargetRepository(db *sql.DB) *TargetRepository {
	return &TargetRepository{DB: db}
}

func (r *TargetRepository) GetTargets(limit int, offset int) ([]models.TargetModel, error) {
	var rows *sql.Rows
	var err error
	var targets []models.TargetModel = []models.TargetModel{}

	if limit == -1 {
		rows, err = r.DB.Query("SELECT id, name, host, port, created_at, updated_at FROM targets")
	} else {
		rows, err = r.DB.Query("SELECT id, name, host, port, created_at, updated_at FROM targets LIMIT ? OFFSET ?", limit, offset)
	}

	if err != nil {
		return targets, err
	}
	defer rows.Close()

	for rows.Next() {
		var target models.TargetModel
		err = rows.Scan(&target.ID, &target.Name, &target.Host, &target.Port, &target.CreatedAt, &target.UpdatedAt)
		if err != nil {
			return targets, err
		}
		targets = append(targets, target)
	}
	return targets, rows.Err()
}

func (r *TargetRepository) CreateTarget(target *models.TargetModel) error {
	_, err := r.DB.Exec(
		`INSERT INTO targets(id, name, host, port, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)`,
		target.ID, target.Name, target.Host, target.Port, target.CreatedAt, target.UpdatedAt,
	)
	return err
}
