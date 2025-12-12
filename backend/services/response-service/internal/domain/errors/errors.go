package errors

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrNotFound  = errors.New("not found")
	ErrForbidden = errors.New("forbidden")
)

func IsNoRows(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
