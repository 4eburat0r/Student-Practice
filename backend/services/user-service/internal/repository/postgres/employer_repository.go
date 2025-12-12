package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"backend/services/user-service/internal/domain/entity"
)

type EmployerRepo struct {
	db *pgxpool.Pool
}

func NewEmployerRepo(db *pgxpool.Pool) *EmployerRepo {
	return &EmployerRepo{db: db}
}

func (r *EmployerRepo) Create(ctx context.Context, e *entity.Employer) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `
        INSERT INTO employers (user_id, company_name, description, photo_path, website_url, requisites)
        VALUES ($1,$2,$3,$4,$5,$6) RETURNING id
    `, e.UserID, e.CompanyName, e.Description, e.PhotoPath, e.WebsiteURL, e.Requisites).Scan(&id)
	return id, err
}

func (r *EmployerRepo) GetByUserID(ctx context.Context, userID int64) (*entity.Employer, error) {
	e := &entity.Employer{}
	err := r.db.QueryRow(ctx, `
        SELECT id, user_id, company_name, description, photo_path, website_url, requisites
        FROM employers WHERE user_id=$1
    `, userID).Scan(&e.ID, &e.UserID, &e.CompanyName, &e.Description, &e.PhotoPath, &e.WebsiteURL, &e.Requisites)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (r *EmployerRepo) UpdateByUserID(ctx context.Context, e *entity.Employer) error {
	ct, err := r.db.Exec(ctx, `
        UPDATE employers SET company_name=$1, description=$2, photo_path=$3, website_url=$4, requisites=$5
        WHERE user_id=$6
    `, e.CompanyName, e.Description, e.PhotoPath, e.WebsiteURL, e.Requisites, e.UserID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *EmployerRepo) DeleteByUserID(ctx context.Context, userID int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM employers WHERE user_id=$1`, userID)
	return err
}
