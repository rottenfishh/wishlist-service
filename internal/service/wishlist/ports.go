package wishlist

import (
	"context"
	"time"
	"wishlist-service/internal/model"

	"github.com/google/uuid"
)

type Service interface {
	Create(ctx context.Context, userID uuid.UUID, name, description string, date time.Time) (*model.Wishlist, error)
	Update(ctx context.Context, userID uuid.UUID, wishlistID int64, name, description *string, date *time.Time) (*model.Wishlist, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Wishlist, error)
	GetByID(ctx context.Context, id int64) (*model.WishlistDetails, error)
	GetByToken(ctx context.Context, token uuid.UUID) (*model.WishlistDetails, error)
	Delete(ctx context.Context, userID uuid.UUID, wishlistID int64) (*model.Wishlist, error)
}

type Repository interface {
	Save(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error)
	Update(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error)
	Delete(ctx context.Context, ID int64) (*model.Wishlist, error)
	GetByID(ctx context.Context, ID int64) (*model.Wishlist, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Wishlist, error)
	GetByToken(ctx context.Context, token uuid.UUID) (*model.Wishlist, error)
}

type GiftReader interface {
	GetByWishlistID(ctx context.Context, id int64) ([]model.Gift, error)
}
