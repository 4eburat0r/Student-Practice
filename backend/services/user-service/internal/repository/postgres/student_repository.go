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

func (r *StudentRepo) Create(ctx context.Context, st *entity.Student) (int64, error) {
	var id int64
	err := r.db.QueryRow(ctx, `
        INSERT INTO students (user_id, first_name, last_name, middle_name, birth_date, photo_path, facility, speciality, course, description)
        VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10) RETURNING id
    `,
		st.UserID, st.FirstName, st.LastName, st.MiddleName, nullIfEmptyDate(st.BirthDate),
		st.PhotoPath, st.Facility, st.Speciality, st.Course, st.Description,
	).Scan(&id)
	return id, err
}

func (r *StudentRepo) GetByUserID(ctx context.Context, userID int64) (*entity.Student, error) {
	st := &entity.Student{}
	var birth pgNullDate
	err := r.db.QueryRow(ctx, `
        SELECT id, user_id, first_name, last_name, middle_name, birth_date, photo_path, facility, speciality, course, description
        FROM students WHERE user_id=$1
    `, userID).Scan(&st.ID, &st.UserID, &st.FirstName, &st.LastName, &st.MiddleName, &birth, &st.PhotoPath, &st.Facility, &st.Speciality, &st.Course, &st.Description)
	if err != nil {
		return nil, err
	}
	if birth.Valid {
		st.BirthDate = birth.Value
	}
	return st, nil
}

func (r *StudentRepo) UpdateByUserID(ctx context.Context, st *entity.Student) error {
	ct, err := r.db.Exec(ctx, `
        UPDATE students SET first_name=$1, last_name=$2, middle_name=$3, birth_date=$4, photo_path=$5, facility=$6, speciality=$7, course=$8, description=$9
        WHERE user_id=$10
    `, st.FirstName, st.LastName, st.MiddleName, nullIfEmptyDate(st.BirthDate), st.PhotoPath, st.Facility, st.Speciality, st.Course, st.Description, st.UserID)
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
