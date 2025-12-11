// internal/repository/postgres/vacancy_repository.go
package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	domainErrors "vacancy-service/internal/domain/errors"
	"vacancy-service/internal/domain/entity"
	"vacancy-service/internal/repository"
)

type vacancyRepository struct {
	pool *pgxpool.Pool
}

func NewVacancyRepository(pool *pgxpool.Pool) repository.VacancyRepository {
	return &vacancyRepository{pool: pool}
}

func (r *vacancyRepository) Create(ctx context.Context, v *entity.Vacancy) error {
	query := `
		INSERT INTO vacancies (employer_id, title, salary, description, hours, format, status, created_dt, updated_dt)
		VALUES ($1,$2,$3,$4,$5,$6,$7,NOW(),NOW())
		RETURNING id, created_dt, updated_dt`
	err := r.pool.QueryRow(ctx, query,
		v.EmployerID, v.Title, v.Salary, v.Description, v.Hours, v.Format, v.Status,
	).Scan(&v.ID, &v.CreatedAt, &v.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert vacancy: %w", err)
	}
	return nil
}

func (r *vacancyRepository) GetByID(ctx context.Context, id int64) (*entity.Vacancy, error) {
	query := `
		SELECT id, employer_id, title, salary, description, hours, format, status, created_dt, updated_dt
		FROM vacancies
		WHERE id = $1`
	var v entity.Vacancy
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&v.ID, &v.EmployerID, &v.Title, &v.Salary, &v.Description,
		&v.Hours, &v.Format, &v.Status, &v.CreatedAt, &v.UpdatedAt,
	)
	if err != nil {
		if domainErrors.IsNoRows(err) {
			return nil, domainErrors.ErrNotFound
		}
		return nil, fmt.Errorf("get vacancy by id: %w", err)
	}
	return &v, nil
}

func (r *vacancyRepository) GetByEmployer(ctx context.Context, employerID int64, limit, offset int) ([]entity.Vacancy, int, error) {
	countQuery := `SELECT COUNT(*) FROM vacancies WHERE employer_id = $1`
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, employerID).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count vacancies: %w", err)
	}

	query := `
		SELECT id, employer_id, title, salary, description, hours, format, status, created_dt, updated_dt
		FROM vacancies
		WHERE employer_id = $1
		ORDER BY created_dt DESC
		LIMIT $2 OFFSET $3`
	rows, err := r.pool.Query(ctx, query, employerID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("list vacancies: %w", err)
	}
	defer rows.Close()

	var res []entity.Vacancy
	for rows.Next() {
		var v entity.Vacancy
		if err := rows.Scan(
			&v.ID, &v.EmployerID, &v.Title, &v.Salary, &v.Description,
			&v.Hours, &v.Format, &v.Status, &v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan vacancy: %w", err)
		}
		res = append(res, v)
	}
	return res, total, nil
}

func (r *vacancyRepository) GetFeed(ctx context.Context, format *string, limit, offset int) ([]entity.Vacancy, int, error) {
	where := `WHERE status = 'active'`
	args := []any{}
	if format != nil {
		where += ` AND format = $1`
		args = append(args, *format)
	}

	countQuery := `SELECT COUNT(*) FROM vacancies ` + where
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count feed: %w", err)
	}

	// добавим limit и offset в аргументы
	args = append(args, limit, offset)

	// формируем плейсхолдеры для limit/offset
	limitPos := len(args) - 1
	offsetPos := len(args)

	query := fmt.Sprintf(`
		SELECT id, employer_id, title, salary, description, hours, format, status, created_dt, updated_dt
		FROM vacancies
		%s
		ORDER BY created_dt DESC
		LIMIT $%d OFFSET $%d`, where, limitPos, offsetPos)

	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("feed vacancies: %w", err)
	}
	defer rows.Close()

	var res []entity.Vacancy
	for rows.Next() {
		var v entity.Vacancy
		if err := rows.Scan(
			&v.ID, &v.EmployerID, &v.Title, &v.Salary, &v.Description,
			&v.Hours, &v.Format, &v.Status, &v.CreatedAt, &v.UpdatedAt,
		); err != nil {
			return nil, 0, fmt.Errorf("scan vacancy: %w", err)
		}
		res = append(res, v)
	}
	return res, total, nil
}

func (r *vacancyRepository) Update(ctx context.Context, v *entity.Vacancy) error {
	query := `
		UPDATE vacancies
		SET title = $1, salary = $2, description = $3, hours = $4, format = $5, updated_dt = NOW()
		WHERE id = $6 AND employer_id = $7
		RETURNING updated_dt`
	err := r.pool.QueryRow(ctx, query,
		v.Title, v.Salary, v.Description, v.Hours, v.Format, v.ID, v.EmployerID,
	).Scan(&v.UpdatedAt)
	if err != nil {
		if domainErrors.IsNoRows(err) {
			return domainErrors.ErrNotFound
		}
		return fmt.Errorf("update vacancy: %w", err)
	}
	return nil
}

func (r *vacancyRepository) Delete(ctx context.Context, id, employerID int64) error {
	cmdTag, err := r.pool.Exec(ctx,
		`DELETE FROM vacancies WHERE id = $1 AND employer_id = $2`, id, employerID)
	if err != nil {
		return fmt.Errorf("delete vacancy: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domainErrors.ErrNotFound
	}
	return nil
}

func (r *vacancyRepository) UpdateStatus(ctx context.Context, id, employerID int64, status entity.VacancyStatus) error {
	cmdTag, err := r.pool.Exec(ctx,
		`UPDATE vacancies SET status = $1, updated_dt = NOW() WHERE id = $2 AND employer_id = $3`,
		status, id, employerID)
	if err != nil {
		return fmt.Errorf("update vacancy status: %w", err)
	}
	if cmdTag.RowsAffected() == 0 {
		return domainErrors.ErrNotFound
	}
	return nil
}
