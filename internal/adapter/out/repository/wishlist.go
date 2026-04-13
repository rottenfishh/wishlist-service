package repository

import (
	"cdek/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type WishlistRepository struct {
	db *sql.DB
}

func NewWishlistRepository(db *sql.DB) *WishlistRepository {
	return &WishlistRepository{db: db}
}

func (r *WishlistRepository) Save(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error) {
	query := `INSERT INTO wishlists(user_id, token, title, description, date) 
              VALUES ($1, $2, $3, $4, $5)
              RETURNING id, user_id, token, title, description, date;`

	row := r.db.QueryRowContext(ctx, query, wishlist.UserID, wishlist.Token,
		wishlist.Title, wishlist.Description, wishlist.Date)

	var savedWishlist model.Wishlist

	err := row.Scan(&savedWishlist.ID, &savedWishlist.UserID, &savedWishlist.Token,
		&savedWishlist.Title, &savedWishlist.Description, &savedWishlist.Date)
	if err != nil {
		return nil, fmt.Errorf("save wishlist: %w", err)
	}

	return &savedWishlist, nil
}

func (r *WishlistRepository) Update(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error) {
	query := `UPDATE wishlists SET title = $1, description = $2, date = $3
              WHERE id = $4;`

	row := r.db.QueryRowContext(ctx, query, wishlist.Title, wishlist.Description, wishlist.Date, wishlist.ID)

	var updatedWishlist model.Wishlist

	err := row.Scan(&updatedWishlist.ID, &updatedWishlist.UserID, &updatedWishlist.Token,
		&updatedWishlist.Title, &updatedWishlist.Description, &updatedWishlist.Date)
	if err != nil {
		return nil, fmt.Errorf("update wishlist: %w", err)
	}

	return &updatedWishlist, nil
}

func (r *WishlistRepository) Delete(ctx context.Context, ID int64) (*model.Wishlist, error) {
	query := `DELETE FROM wishlists WHERE id = $1
              RETURNING id, user_id, token, title, description, date;`

	row := r.db.QueryRowContext(ctx, query, ID)

	var wishlist model.Wishlist

	err := row.Scan(&wishlist.ID, &wishlist.UserID, &wishlist.Token,
		&wishlist.Title, &wishlist.Description, &wishlist.Date)
	if err != nil {
		return nil, fmt.Errorf("delete wishlist: %w", err)
	}

	return &wishlist, nil

}

func (r *WishlistRepository) GetByID(ctx context.Context, ID int64) (*model.Wishlist, error) {
	query := `SELECT id, user_id, token, title, description, date
              FROM wishlists 
              WHERE id = $1;`

	row := r.db.QueryRowContext(ctx, query, ID)

	var wishlist model.Wishlist

	err := row.Scan(&wishlist.ID, &wishlist.UserID, &wishlist.Token,
		&wishlist.Title, &wishlist.Description, &wishlist.Date)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get by id wishlist: %w", err)
	}

	return &wishlist, nil
}

func (r *WishlistRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Wishlist, error) {
	query := `SELECT id, user_id, token, title, description, date
              FROM wishlists
              WHERE user_id = $1;`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("get by user id wishlist: %w", err)
	}
	defer rows.Close()

	var wishlists []model.Wishlist
	for rows.Next() {
		var wishlist model.Wishlist
		err = rows.Scan(&wishlist.ID, &wishlist.UserID, &wishlist.Token, &wishlist.Title,
			&wishlist.Description, &wishlist.Date)
		if err != nil {
			return nil, fmt.Errorf("get by user id wishlist: %w", err)
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get by user id wishlists: %w", err)
	}

	return wishlists, nil
}

func (r *WishlistRepository) GetByToken(ctx context.Context, token uuid.UUID) (*model.Wishlist, error) {
	query := `SELECT id, user_id, token, title, description, date
              FROM wishlists
              WHERE token = $1;`
	row := r.db.QueryRowContext(ctx, query, token)

	var wishlist model.Wishlist
	err := row.Scan(&wishlist.ID, &wishlist.UserID, &wishlist.Token,
		&wishlist.Title, &wishlist.Description, &wishlist.Date)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
	}

	return &wishlist, nil
}
