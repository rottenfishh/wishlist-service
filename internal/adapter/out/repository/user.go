package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"wishlist-service/internal/model"

	"github.com/jackc/pgx/v5/pgconn"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) SaveUser(ctx context.Context, user *model.User) (*model.User, error) {
	query := `INSERT INTO users(id, email, password_hash) VALUES ($1, $2, $3)
                RETURNING id, email, created_at`

	row := r.db.QueryRowContext(ctx, query, user.ID, user.Email, user.PasswordHash)

	var savedUser model.User
	err := row.Scan(&savedUser.ID, &savedUser.Email, &savedUser.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, model.ErrUserAlreadyExists
		}
		return nil, fmt.Errorf("save user: %w", err)
	}

	return &savedUser, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT id, email, created_at, password_hash FROM users WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var user model.User
	err := row.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, created_at, password_hash FROM users WHERE email = $1`

	row := r.db.QueryRowContext(ctx, query, email)

	var user model.User
	err := row.Scan(&user.ID, &user.Email, &user.CreatedAt, &user.PasswordHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}
