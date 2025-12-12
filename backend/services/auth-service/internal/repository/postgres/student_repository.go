package postgres

import (
    "auth-service/internal/domain/entity"
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type StudentRepository struct {
    db *pgxpool.Pool
}

func NewStudentRepository(db *pgxpool.Pool) *StudentRepository {
    return &StudentRepository{db: db}
}

func (r *StudentRepository) Create(ctx context.Context, userID int) (*entity.Student, error) {
    query := `
        INSERT INTO students (user_id, birth_date, created_dt, updated_dt)
        VALUES ($1, $2, $3, $4)
        RETURNING id, user_id, birth_date, created_dt, updated_dt
    `

    var student entity.Student
    now := time.Now()
    birthDate := time.Now().AddDate(-18, 0, 0)

    err := r.db.QueryRow(ctx, query, userID, birthDate, now, now).Scan(
        &student.ID,
        &student.UserID,
        &student.BirthDate,
        &student.CreatedDt,
        &student.UpdatedDt,
    )

    if err != nil {
        return nil, fmt.Errorf("failed to create student profile: %w", err)
    }

    return &student, nil
}

func (r *StudentRepository) GetByID(ctx context.Context, id int) (*entity.Student, error) {
    query := `
        SELECT id, user_id, first_name, last_name, middle_name, birth_date,
               photo_path, facility, speciality, course, description,
               created_dt, updated_dt
        FROM students
        WHERE id = $1
    `

    var student entity.Student
    err := r.db.QueryRow(ctx, query, id).Scan(
        &student.ID,
        &student.UserID,
        &student.FirstName,
        &student.LastName,
        &student.MiddleName,
        &student.BirthDate,
        &student.PhotoPath,
        &student.Facility,
        &student.Speciality,
        &student.Course,
        &student.Description,
        &student.CreatedDt,
        &student.UpdatedDt,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get student by id: %w", err)
    }

    return &student, nil
}

func (r *StudentRepository) GetByUserID(ctx context.Context, userID int) (*entity.Student, error) {
    query := `
        SELECT id, user_id, first_name, last_name, middle_name, birth_date,
               photo_path, facility, speciality, course, description,
               created_dt, updated_dt
        FROM students
        WHERE user_id = $1
    `

    var student entity.Student
    err := r.db.QueryRow(ctx, query, userID).Scan(
        &student.ID,
        &student.UserID,
        &student.FirstName,
        &student.LastName,
        &student.MiddleName,
        &student.BirthDate,
        &student.PhotoPath,
        &student.Facility,
        &student.Speciality,
        &student.Course,
        &student.Description,
        &student.CreatedDt,
        &student.UpdatedDt,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get student by user id: %w", err)
    }

    return &student, nil
}

func (r *StudentRepository) Update(ctx context.Context, student *entity.Student) error {
    query := `
        UPDATE students
        SET first_name = $1,
            last_name = $2,
            middle_name = $3,
            birth_date = $4,
            photo_path = $5,
            facility = $6,
            speciality = $7,
            course = $8,
            description = $9,
            updated_dt = $10
        WHERE id = $11
    `

    result, err := r.db.Exec(
        ctx,
        query,
        student.FirstName,
        student.LastName,
        student.MiddleName,
        student.BirthDate,
        student.PhotoPath,
        student.Facility,
        student.Speciality,
        student.Course,
        student.Description,
        time.Now(),
        student.ID,
    )

    if err != nil {
        return fmt.Errorf("failed to update student profile: %w", err)
    }

    if result.RowsAffected() == 0 {
        return fmt.Errorf("student with id %d not found", student.ID)
    }

    return nil
}

func (r *StudentRepository) Delete(ctx context.Context, id int) error {
    query := `DELETE FROM students WHERE id = $1`

    result, err := r.db.Exec(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete student profile: %w", err)
    }

    if result.RowsAffected() == 0 {
        return fmt.Errorf("student with id %d not found", id)
    }

    return nil
}
