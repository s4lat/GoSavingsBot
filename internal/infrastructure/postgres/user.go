package postgres

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/s4lat/gosavingsbot/internal/domain"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

// NewUserRepo initializes a new instance of UserRepo with the provided pgxpool.Pool.
func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		pool: pool,
	}
}

// Create inserts a new user into the database.
func (r *UserRepo) Create(ctx context.Context, user domain.User) (domain.User, error) {
	query := `INSERT INTO users (id, username) VALUES ($1, $2) RETURNING id, username;`
	if err := r.pool.QueryRow(ctx, query, user.Id, user.Username).Scan(&user.Id, &user.Username); err != nil {
		return domain.User{}, wrapError(err)
	}
	return user, nil
}

// GetById retrieves a user by their ID.
func (r *UserRepo) GetById(ctx context.Context, id int64) (domain.User, error) {
	query := `SELECT id, username FROM users WHERE id = $1`
	row := r.pool.QueryRow(ctx, query, id)

	var user domain.User
	err := row.Scan(&user.Id, &user.Username)
	if err != nil {
		return domain.User{}, wrapError(err)
	}

	return user, nil
}

// GetByUsername retrieves a user by their username.
func (r *UserRepo) GetByUsername(ctx context.Context, username string) (domain.User, error) {
	query := `SELECT id, username FROM users WHERE username = $1`
	row := r.pool.QueryRow(ctx, query, username)

	var user domain.User
	err := row.Scan(&user.Id, &user.Username)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

// UpdateUsername updates the username of a user in the database and updates the user object.
func (r *UserRepo) UpdateUsername(ctx context.Context, user domain.User, newUsername string) (domain.User, error) {
	query := `UPDATE users SET username = $1 WHERE id = $2`
	commandTag, err := r.pool.Exec(ctx, query, newUsername, user.Id)
	if err != nil {
		return domain.User{}, wrapError(err)
	}

	// Check if the update affected exactly one row
	if commandTag.RowsAffected() != 1 {
		return domain.User{}, wrapError(errors.New("no rows were updated"))
	}

	// Update the user object
	user.Username = newUsername
	return user, nil
}

func (r *UserRepo) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	query := `UPDATE users
				SET username = $2
			  WHERE id = $1
				RETURNING id, username;`

	err := r.pool.QueryRow(ctx, query, user.Id, user.Username).Scan(&user.Id, &user.Username)
	if err != nil {
		return domain.User{}, wrapError(err)
	}

	return user, nil
}
