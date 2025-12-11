package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"backend/services/user-service/internal/domain/entity"
)

type StudentRepo struct {
	db *pgxpool.Pool
}

func NewStudentRepo(db *pgxpool.Pool) *StudentRepo {
	return &StudentRepo{db: db}
}

func (r *StudentRepo) Create(ctx context.Context, s *entity.Student) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `INSERT INTO students(user_id, first_name, last_name, middle_name, photo_path, facility, course, description, created_dt, updated_dt)
		VALUES($1,$2,$3,$4,$5,$6,$7,$8, now(), now()) RETURNING id`,
		s.UserID, s.FirstName, s.LastName, s.MiddleName, s.PhotoPath, s.Facility, s.Course, s.Description).Scan(&id)
	return id, err
}

func (r *StudentRepo) GetByUserID(ctx context.Context, userID int64) (*entity.Student, error) {
	s := &entity.Student{}
	row := r.db.QueryRow(ctx, `SELECT id,user_id,first_name,last_name,middle_name,photo_path,facility,course,description,created_dt,updated_dt FROM students WHERE user_id=$1`, userID)
	if err := row.Scan(&s.ID, &s.UserID, &s.FirstName, &s.LastName, &s.MiddleName, &s.PhotoPath, &s.Facility, &s.Course, &s.Description, &s.CreatedAt, &s.UpdatedAt); err != nil {
		return nil, err
	}
	return s, nil
}

func (r *StudentRepo) UpdateByUserID(ctx context.Context, s *entity.Student) error {
	ct, err := r.db.Exec(ctx, `UPDATE students SET first_name=$1, last_name=$2, middle_name=$3, photo_path=$4, facility=$5, course=$6, description=$7, updated_dt=now() WHERE user_id=$8`,
		s.FirstName, s.LastName, s.MiddleName, s.PhotoPath, s.Facility, s.Course, s.Description, s.UserID)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}

func (r *StudentRepo) DeleteByUserID(ctx context.Context, userID int64) error {
	_, err := r.db.Exec(ctx, `DELETE FROM students WHERE user_id=$1`, userID)
	return err
}
