package postgres

import (
    "auth-service/internal/domain/entity"
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type EmployerRepository struct {
    db *pgxpool.Pool
}

func NewEmployerRepository(db *pgxpool.Pool) *EmployerRepository {
    return &EmployerRepository{db: db}
}

func (r *EmployerRepository) Create(ctx context.Context, userID int) (*entity.Employer, error) {
    query := `
        INSERT INTO employers (user_id, company_name, created_dt, updated_dt)
        VALUES ($1, $2, $3, $4)
        RETURNING id, user_id, company_name, created_dt, updated_dt
    `

    var employer entity.Employer
    now := time.Now()
    defaultCompanyName := "My Company"

    err := r.db.QueryRow(ctx, query, userID, defaultCompanyName, now, now).Scan(
        &employer.ID,
        &employer.UserID,
        &employer.CompanyName,
        &employer.CreatedDt,
        &employer.UpdatedDt,
    )

    if err != nil {
        return nil, fmt.Errorf("failed to create employer profile: %w", err)
    }

    return &employer, nil
}

func (r *EmployerRepository) GetByID(ctx context.Context, id int) (*entity.Employer, error) {
    query := `
        SELECT id, user_id, company_name, description, photo_path,
               website_url, requisites, created_dt, updated_dt
        FROM employers
        WHERE id = $1
    `

    var employer entity.Employer
    err := r.db.QueryRow(ctx, query, id).Scan(
        &employer.ID,
        &employer.UserID,
        &employer.CompanyName,
        &employer.Description,
        &employer.PhotoPath,
        &employer.WebsiteURL,
        &employer.Requisites,
        &employer.CreatedDt,
        &employer.UpdatedDt,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get employer by id: %w", err)
    }

    return &employer, nil
}

func (r *EmployerRepository) GetByUserID(ctx context.Context, userID int) (*entity.Employer, error) {
    query := `
        SELECT id, user_id, company_name, description, photo_path,
               website_url, requisites, created_dt, updated_dt
        FROM employers
        WHERE user_id = $1
    `

    var employer entity.Employer
    err := r.db.QueryRow(ctx, query, userID).Scan(
        &employer.ID,
        &employer.UserID,
        &employer.CompanyName,
        &employer.Description,
        &employer.PhotoPath,
        &employer.WebsiteURL,
        &employer.Requisites,
        &employer.CreatedDt,
        &employer.UpdatedDt,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get employer by user id: %w", err)
    }

    return &employer, nil
}

func (r *EmployerRepository) Update(ctx context.Context, employer *entity.Employer) error {
    query := `
        UPDATE employers
        SET company_name = $1,
            description = $2,
            photo_path = $3,
            website_url = $4,
            requisites = $5,
            updated_dt = $6
        WHERE id = $7
    `

    result, err := r.db.Exec(
        ctx,
        query,
        employer.CompanyName,
        employer.Description,
        employer.PhotoPath,
        employer.WebsiteURL,
        employer.Requisites,
        time.Now(),
        employer.ID,
    )

    if err != nil {
        return fmt.Errorf("failed to update employer profile: %w", err)
    }

    if result.RowsAffected() == 0 {
        return fmt.Errorf("employer with id %d not found", employer.ID)
    }

    return nil
}

func (r *EmployerRepository) Delete(ctx context.Context, id int) error {
    query := `DELETE FROM employers WHERE id = $1`

    result, err := r.db.Exec(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete employer profile: %w", err)
    }

    if result.RowsAffected() == 0 {
        return fmt.Errorf("employer with id %d not found", id)
    }

    return nil
}
