package postgres

import (
    "auth-service/internal/domain/entity"
    "context"
    "fmt"
    "time"

    "github.com/jackc/pgx/v5"
    "github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
    db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
    return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, user *entity.User) error {
    query := `
        INSERT INTO users (email, password_hash, role, created_dt, updated_dt)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, created_dt, updated_dt
    `

    now := time.Now()
    err := r.db.QueryRow(
        ctx,
        query,
        user.Email,
        user.PasswordHash,
        user.Role,
        now,
        now,
    ).Scan(&user.ID, &user.CreatedDt, &user.UpdatedDt)

    if err != nil {
        return fmt.Errorf("failed to create user: %w", err)
    }

    return nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*entity.User, error) {
    query := `
        SELECT id, email, password_hash, role, created_dt, updated_dt
        FROM users
        WHERE email = $1
    `

    var user entity.User
    err := r.db.QueryRow(ctx, query, email).Scan(
        &user.ID,
        &user.Email,
        &user.PasswordHash,
        &user.Role,
        &user.CreatedDt,
        &user.UpdatedDt,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get user by email: %w", err)
    }

    return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id int) (*entity.User, error) {
    query := `
        SELECT id, email, password_hash, role, created_dt, updated_dt
        FROM users
        WHERE id = $1
    `

    var user entity.User
    err := r.db.QueryRow(ctx, query, id).Scan(
        &user.ID,
        &user.Email,
        &user.PasswordHash,
        &user.Role,
        &user.CreatedDt,
        &user.UpdatedDt,
    )

    if err != nil {
        if err == pgx.ErrNoRows {
            return nil, nil
        }
        return nil, fmt.Errorf("failed to get user by id: %w", err)
    }

    return &user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *entity.User) error {
    query := `
        UPDATE users
        SET email = $1,
            password_hash = $2,
            role = $3,
            updated_dt = $4
        WHERE id = $5
    `

    result, err := r.db.Exec(
        ctx,
        query,
        user.Email,
        user.PasswordHash,
        user.Role,
        time.Now(),
        user.ID,
    )

    if err != nil {
        return fmt.Errorf("failed to update user: %w", err)
    }

    if result.RowsAffected() == 0 {
        return fmt.Errorf("user with id %d not found", user.ID)
    }

    return nil
}

func (r *UserRepository) Delete(ctx context.Context, id int) error {
    query := `DELETE FROM users WHERE id = $1`

    result, err := r.db.Exec(ctx, query, id)
    if err != nil {
        return fmt.Errorf("failed to delete user: %w", err)
    }

    if result.RowsAffected() == 0 {
        return fmt.Errorf("user with id %d not found", id)
    }

    return nil
}

func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
    query := `SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)`

    var exists bool
    err := r.db.QueryRow(ctx, query, email).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("failed to check email existence: %w", err)
    }

    return exists, nil
}
