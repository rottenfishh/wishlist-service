package gift

import (
	"context"
	"wishlist-service/internal/model"

	"github.com/google/uuid"
)

type Service interface {
	Save(ctx context.Context, userID uuid.UUID, wishlistID int64, name, description, link string, priority int) (*model.Gift, error)
	GetByID(ctx context.Context, userID uuid.UUID, wishlistID, ID int64) (*model.Gift, error)
	Update(ctx context.Context, userID uuid.UUID, wishlistID, ID int64, name, description, link *string, priority *int) (*model.Gift, error)
	Book(ctx context.Context, ID int64, token uuid.UUID) (*model.Gift, error)
	Delete(ctx context.Context, userID uuid.UUID, wishlistID, ID int64) (*model.Gift, error)
}

type Repository interface {
	Save(ctx context.Context, gift *model.Gift) (*model.Gift, error)
	Update(ctx context.Context, gift *model.Gift) (*model.Gift, error)
	Book(ctx context.Context, ID int64, token uuid.UUID) (*model.Gift, error)
	GetByID(ctx context.Context, id int64) (*model.Gift, error)
	GetByWishlistID(ctx context.Context, id int64) ([]model.Gift, error)
	Delete(ctx context.Context, giftID int64) (*model.Gift, error)
}

type WishlistReader interface {
	GetByID(ctx context.Context, id int64) (*model.Wishlist, error)
	GetByToken(ctx context.Context, token uuid.UUID) (*model.Wishlist, error)
}
