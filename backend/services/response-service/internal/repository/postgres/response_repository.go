package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	domainErrors "response-service/internal/domain/errors"
	"response-service/internal/domain/entity"
	"response-service/internal/repository"
)

type responseRepository struct {
	pool *pgxpool.Pool
}

func NewResponseRepository(pool *pgxpool.Pool) repository.ResponseRepository {
	return &responseRepository{pool: pool}
}

func (r *responseRepository) Create(ctx context.Context, resp *entity.Response) error {
	query := `
		INSERT INTO responses (vacancy_id, student_id, status, message, created_dt)
		VALUES ($1,$2,$3,$4,NOW())
		RETURNING id, created_dt`
	err := r.pool.QueryRow(ctx, query,
		resp.VacancyID, resp.StudentID, resp.Status, resp.Message,
	).Scan(&resp.ID, &resp.CreatedAt)
	if err != nil {
		return fmt.Errorf("insert response: %w", err)
	}
	return nil
}

func (r *responseRepository) GetByID(ctx context.Context, id int64, studentID int64) (*entity.Response, error) {
	query := `
		SELECT id, vacancy_id, student_id, status, message, created_dt
		FROM responses
		WHERE id = $1 AND student_id = $2`
	var resp entity.Response
	err := r.pool.QueryRow(ctx, query, id, studentID).Scan(
		&resp.ID, &resp.VacancyID, &resp.StudentID,
		&resp.Status, &resp.Message, &resp.CreatedAt,
	)
	if err != nil {
		if domainErrors.IsNoRows(err) {
			return nil, domainErrors.ErrNotFound
		}
		return nil, fmt.Errorf("get response by id: %w", err)
	}
	return &resp, nil
}

func (r *responseRepository) ListByStudent(ctx context.Context, studentID int64, limit, offset int) ([]entity.Response, int, error) {
	countQuery := `SELECT COUNT(*) FROM responses WHERE student_id = $1`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, studentID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count responses: %w", err)
	}

	query := `
		SELECT id, vacancy_id, student_id, status, message, created_dt
		FROM responses
		WHERE student_id = $1
		ORDER BY created_dt DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, studentID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list responses: %w", err)
	}
	defer rows.Close()

	var res []entity.Response
	for rows.Next() {
		var resp entity.Response
		if err := rows.Scan(
			&resp.ID, &resp.VacancyID, &resp.StudentID,
			&resp.Status, &resp.Message, &resp.CreatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan response: %w", err)
		}
		res = append(res, resp)
	}
	return res, total, nil
}

func (r *responseRepository) Delete(ctx context.Context, id int64, studentID int64) error {
	cmdTag, err := r.pool.Exec(ctx,
		`DELETE FROM responses WHERE id = $1 AND student_id = $2`, id, studentID)
	if err != nil {
		return fmt.Errorf("delete response: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domainErrors.ErrNotFound
	}
	return nil
}
