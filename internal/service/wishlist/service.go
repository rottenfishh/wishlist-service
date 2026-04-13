package wishlist

import (
	"cdek/internal/model"
	"cdek/internal/service/gift"
	"context"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo      Repository
	giftsRepo gift.Repository
}

func NewService(repo Repository, giftsRepo gift.Repository) Service {
	return &service{repo: repo, giftsRepo: giftsRepo}
}

func (s *service) Create(ctx context.Context, userID uuid.UUID, name, description string,
	date time.Time) (*model.Wishlist, error) {

	wishlist := &model.Wishlist{
		UserID:      userID,
		Token:       uuid.New(),
		Title:       name,
		Description: description,
		Date:        date,
	}

	return s.repo.Save(ctx, wishlist)
}

func (s *service) Update(ctx context.Context, userID uuid.UUID, wishlistID int64, name,
	description *string, date *time.Time) (*model.Wishlist, error) {

	wishlist, err := s.repo.GetByID(ctx, wishlistID)
	if err != nil {
		return nil, err
	}

	if wishlist.UserID != userID {
		return nil, model.ErrForbidden
	}

	if name != nil {
		wishlist.Title = *name
	}
	if description != nil {
		wishlist.Description = *description
	}
	if date != nil {
		wishlist.Date = *date
	}
	return s.repo.Update(ctx, wishlist)
}

func (s *service) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Wishlist, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *service) GetByID(ctx context.Context, id int64) (*model.WishlistDetails, error) {
	wishlist, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	items, err := s.giftsRepo.GetByWishlistID(ctx, wishlist.ID)
	if err != nil {
		return nil, err
	}

	wishlistWithItems := &model.WishlistDetails{
		Wishlist: *wishlist,
		Gifts:    items,
	}
	return wishlistWithItems, nil
}

func (s *service) GetByToken(ctx context.Context, token uuid.UUID) (*model.WishlistDetails, error) {
	wishlist, err := s.repo.GetByToken(ctx, token)
	if err != nil {
		return nil, err
	}

	items, err := s.giftsRepo.GetByWishlistID(ctx, wishlist.ID)
	if err != nil {
		return nil, err
	}

	wishlistWithItems := &model.WishlistDetails{
		Wishlist: *wishlist,
		Gifts:    items,
	}
	return wishlistWithItems, nil
}

func (s *service) Delete(ctx context.Context, userID uuid.UUID, wishlistID int64) (*model.Wishlist, error) {
	wishlist, err := s.repo.GetByID(ctx, wishlistID)
	if err != nil {
		return nil, err
	}

	if wishlist.UserID != userID {
		return nil, model.ErrForbidden
	}

	return s.repo.Delete(ctx, wishlistID)
}
