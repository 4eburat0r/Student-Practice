package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"

	"backend/services/user-service/internal/domain/entity"
)

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{db: db}
}

func (r *UserRepo) Create(ctx context.Context, email, passwordHash, role string, idOverride *int64) (int64, error) {
	var id int64
	if idOverride != nil {
		id = *idOverride
		_, err := r.db.Exec(ctx,
			`INSERT INTO users (id, email, password_hash, role, email_verified, created_at, updated_at) VALUES ($1,$2,$3,$4,true,now(),now())`,
			id, email, passwordHash, role)
		if err != nil {
			return 0, err
		}
		// обновляем последовательность, чтобы дальше авто-инкремент не конфликтовал
		_, _ = r.db.Exec(ctx, `SELECT setval(pg_get_serial_sequence('users','id'), GREATEST((SELECT MAX(id) FROM users), $1))`, id)
		return id, nil
	}
	err := r.db.QueryRow(ctx,
		`INSERT INTO users (email, password_hash, role, email_verified, created_at, updated_at) VALUES ($1,$2,$3,true,now(),now()) RETURNING id`,
		email, passwordHash, role).Scan(&id)
	return id, err
}

func (r *UserRepo) GetByEmail(ctx context.Context, email string) (*entity.User, string, error) {
	u := &entity.User{}
	var ph string
	row := r.db.QueryRow(ctx, `SELECT id, email, role, email_verified, created_at, updated_at, password_hash FROM users WHERE email=$1`, email)
	if err := row.Scan(&u.ID, &u.Email, &u.Role, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt, &ph); err != nil {
		return nil, "", err
	}
	return u, ph, nil
}

func (r *UserRepo) GetByID(ctx context.Context, id int64) (*entity.User, error) {
	u := &entity.User{}
	row := r.db.QueryRow(ctx, `SELECT id, email, role, email_verified, created_at, updated_at FROM users WHERE id=$1`, id)
	if err := row.Scan(&u.ID, &u.Email, &u.Role, &u.EmailVerified, &u.CreatedAt, &u.UpdatedAt); err != nil {
		return nil, err
	}
	return u, nil
}

func (r *UserRepo) Delete(ctx context.Context, id int64) error {
	ct, err := r.db.Exec(ctx, `DELETE FROM users WHERE id=$1`, id)
	if err != nil {
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("not found")
	}
	return nil
}
