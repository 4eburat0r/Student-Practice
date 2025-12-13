package postgres

import (
	"context"
	"fmt"

	"resume-service/internal/domain/entity"
	"resume-service/internal/domain/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ResumeRepository struct {
	db *pgxpool.Pool
}

func NewResumeRepository(db *pgxpool.Pool) *ResumeRepository {
	return &ResumeRepository{db: db}
}

func (r *ResumeRepository) Create(ctx context.Context, resume *entity.Resume) error {
	query := `
		INSERT INTO resumes (student_id, title, summary, status, created_dt, updated_dt)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_dt, updated_dt
	`

	err := r.db.QueryRow(
		ctx,
		query,
		resume.StudentID,
		resume.Title,
		resume.Summary,
		resume.Status,
	).Scan(&resume.ID, &resume.CreatedAt, &resume.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create resume: %w", err)
	}

	return nil
}

func (r *ResumeRepository) GetByID(ctx context.Context, id int) (*entity.Resume, error) {
	query := `
		SELECT id, student_id, title, summary, status, created_dt, updated_dt
		FROM resumes
		WHERE id = $1
	`

	var resume entity.Resume
	err := r.db.QueryRow(ctx, query, id).Scan(
		&resume.ID,
		&resume.StudentID,
		&resume.Title,
		&resume.Summary,
		&resume.Status,
		&resume.CreatedAt,
		&resume.UpdatedAt,
	)

	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, errors.ErrResumeNotFound
		}
		return nil, fmt.Errorf("failed to get resume by id: %w", err)
	}

	return &resume, nil
}

func (r *ResumeRepository) GetByStudentID(ctx context.Context, studentID int) ([]*entity.Resume, error) {
	query := `
		SELECT id, student_id, title, summary, status, created_dt, updated_dt
		FROM resumes
		WHERE student_id = $1
		ORDER BY created_dt DESC
	`

	rows, err := r.db.Query(ctx, query, studentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get resumes by student id: %w", err)
	}
	defer rows.Close()

	var resumes []*entity.Resume
	for rows.Next() {
		var resume entity.Resume
		if err := rows.Scan(
			&resume.ID,
			&resume.StudentID,
			&resume.Title,
			&resume.Summary,
			&resume.Status,
			&resume.CreatedAt,
			&resume.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan resume: %w", err)
		}
		resumes = append(resumes, &resume)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return resumes, nil
}

func (r *ResumeRepository) GetFeed(ctx context.Context, limit, offset int) ([]*entity.Resume, error) {
	query := `
		SELECT id, student_id, title, summary, status, created_dt, updated_dt
		FROM resumes
		WHERE status = $1
		ORDER BY created_dt DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(ctx, query, entity.ResumeStatusActive, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get resume feed: %w", err)
	}
	defer rows.Close()

	var resumes []*entity.Resume
	for rows.Next() {
		var resume entity.Resume
		if err := rows.Scan(
			&resume.ID,
			&resume.StudentID,
			&resume.Title,
			&resume.Summary,
			&resume.Status,
			&resume.CreatedAt,
			&resume.UpdatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan resume: %w", err)
		}
		resumes = append(resumes, &resume)
	}

	return resumes, rows.Err()
}

func (r *ResumeRepository) Update(ctx context.Context, resume *entity.Resume) error {
	query := `
		UPDATE resumes
		SET title = $1,
		    summary = $2,
		    updated_dt = NOW()
		WHERE id = $3
		RETURNING updated_dt
	`

	err := r.db.QueryRow(
		ctx,
		query,
		resume.Title,
		resume.Summary,
		resume.ID,
	).Scan(&resume.UpdatedAt)

	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.ErrResumeNotFound
		}
		return fmt.Errorf("failed to update resume: %w", err)
	}

	return nil
}

func (r *ResumeRepository) Delete(ctx context.Context, id int) error {
	query := `DELETE FROM resumes WHERE id = $1`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete resume: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.ErrResumeNotFound
	}

	return nil
}

func (r *ResumeRepository) Publish(ctx context.Context, id int) error {
	query := `
		UPDATE resumes
		SET status = $1,
		    updated_dt = NOW()
		WHERE id = $2
	`

	result, err := r.db.Exec(ctx, query, entity.ResumeStatusActive, id)
	if err != nil {
		return fmt.Errorf("failed to publish resume: %w", err)
	}

	if result.RowsAffected() == 0 {
		return errors.ErrResumeNotFound
	}

	return nil
}
