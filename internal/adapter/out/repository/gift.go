package repository

import (
	"cdek/internal/model"
	"context"
	"errors"
	"fmt"

	"database/sql"

	"github.com/google/uuid"
)

type GiftRepository struct {
	db *sql.DB
}

func NewGiftRepository(db *sql.DB) *GiftRepository {
	return &GiftRepository{db: db}
}

func (r *GiftRepository) Save(ctx context.Context, gift *model.Gift) (*model.Gift, error) {
	query := `INSERT INTO gifts(wishlist_id, name, description, link, priority) 
            VALUES ($1, $2, $3, $4, $5)
            RETURNING id, wishlist_id, name, description, link, priority, booked`

	row := r.db.QueryRowContext(ctx, query, gift.WishlistID, gift.Name, gift.Description, gift.Link, gift.Priority)

	var savedGift model.Gift
	err := row.Scan(&savedGift.ID, &savedGift.WishlistID, &savedGift.Name,
		&savedGift.Description, &savedGift.Link, &savedGift.Priority, &savedGift.Booked)
	if err != nil {
		return nil, fmt.Errorf("save gift: %w", err)
	}

	return &savedGift, nil
}

func (r *GiftRepository) Update(ctx context.Context, gift *model.Gift) (*model.Gift, error) {
	query := `UPDATE gifts SET name = $1, description = $2, link = $3, priority = $4
            WHERE id = $5
            RETURNING id, wishlist_id, name, description, link, priority, booked`

	row := r.db.QueryRowContext(ctx, query, gift.Name, gift.Description, gift.Link, gift.Priority,
		gift.ID)

	var updatedGift model.Gift
	err := row.Scan(&updatedGift.ID, &updatedGift.WishlistID, &updatedGift.Name,
		&updatedGift.Description, &updatedGift.Link, &updatedGift.Priority, &updatedGift.Booked)
	if err != nil {
		return nil, fmt.Errorf("update gift: %w", err)
	}

	return &updatedGift, nil
}

func (r *GiftRepository) Book(ctx context.Context, ID int64, token uuid.UUID) (*model.Gift, error) {
	query := `UPDATE gifts
              SET booked = true
              WHERE id = $1
                AND wishlist_id = (
                    SELECT id FROM wishlists WHERE token = $2
                )
                AND booked = false
              RETURNING id, wishlist_id, name, description, link, priority  `

	row := r.db.QueryRowContext(ctx, query, ID, token)

	var gift model.Gift
	err := row.Scan(&gift.ID, &gift.WishlistID, &gift.Name,
		&gift.Description, &gift.Link, &gift.Priority, &gift.Booked)

	if err != nil {
		return nil, fmt.Errorf("book gift: %w", err)
	}

	return &gift, nil

}

func (r *GiftRepository) GetByID(ctx context.Context, id int64) (*model.Gift, error) {
	query := `SELECT id, wishlist_id, name, description, link, priority, booked
            FROM gifts WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	var gift model.Gift
	err := row.Scan(&gift.ID, &gift.WishlistID, &gift.Name,
		&gift.Description, &gift.Link, &gift.Priority, &gift.Booked)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get by id gift: %w", err)
	}

	return &gift, nil
}

func (r *GiftRepository) GetByWishlistID(ctx context.Context, id int64) ([]model.Gift, error) {
	query := `SELECT id, wishlist_id, name, description, link, priority, booked
            FROM gifts WHERE wishlist_id = $1`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var gifts []model.Gift
	for rows.Next() {
		var gift model.Gift
		err = rows.Scan(&gift.ID, &gift.WishlistID, &gift.Name,
			&gift.Description, &gift.Link, &gift.Priority, &gift.Booked)
		if err != nil {
			return nil, fmt.Errorf("get by wishlist id gift: %w", err)
		}
		gifts = append(gifts, gift)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("get by wishlist id gift: %w", err)
	}

	return gifts, nil
}

func (r *GiftRepository) Delete(ctx context.Context, giftID int64) (*model.Gift, error) {
	query := `DELETE FROM gifts WHERE id = $1 
              RETURNING id, wishlist_id, name, description, link, priority, booked`
	row := r.db.QueryRowContext(ctx, query, giftID)

	var gift model.Gift
	err := row.Scan(&gift.ID, &gift.WishlistID, &gift.Name,
		&gift.Description, &gift.Link, &gift.Priority, &gift.Booked)
	if err != nil {
		return nil, fmt.Errorf("delete gift: %w", err)
	}

	return &gift, nil
}
