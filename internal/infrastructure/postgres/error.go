package postgres

import (
	"errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/s4lat/gosavingsbot/internal/domain"
)

func wrapError(err error) error {
	pgError, ok := err.(*pgconn.PgError)

	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return errors.Join(err, domain.ErrNotFound)
	case ok && pgError.Code == "23505":
		// Handle unique constraint violation using PostgreSQL error code 23505
		return errors.Join(err, domain.ErrAlreadyExists)
	case err == nil:
		return nil
	default:
		return errors.Join(err, domain.ErrUnknown)
	}
}
