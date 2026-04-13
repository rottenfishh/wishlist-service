package repository

import (
	"cdek/internal/model"
	"context"
	"database/sql"
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

	row := r.db.QueryRowContext(ctx, query, user.Id, user.Email, user.PasswordHash)

	var savedUser model.User
	err := row.Scan(&savedUser.Id, &savedUser.Email, &savedUser.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &savedUser, nil
}

func (r *UserRepository) GetUserById(ctx context.Context, id string) (*model.User, error) {
	query := `SELECT * FROM users WHERE id = $1`

	row := r.db.QueryRowContext(ctx, query, id)

	var user model.User
	err := row.Scan(&user.Id, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT * FROM users WHERE email = $1`

	row := r.db.QueryRowContext(ctx, query, email)

	var user model.User
	err := row.Scan(&user.Id, &user.Email, &user.CreatedAt)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
