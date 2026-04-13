package wishlist

import (
	"cdek/internal/model"
	"context"
	"time"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
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
	description string, date time.Time) (*model.Wishlist, error) {

	wishlist, err := s.repo.GetByID(ctx, wishlistID)
	if err != nil {
		return nil, err
	}

	if wishlist.UserID != userID {
		return nil, model.ErrForbidden
	}

	wishlist.Title = name
	wishlist.Description = description
	wishlist.Date = date
	return s.repo.Update(ctx, wishlist)
}

func (s *service) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Wishlist, error) {
	return s.repo.GetByUserID(ctx, userID)
}

func (s *service) GetByID(ctx context.Context, id int64) (*model.Wishlist, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *service) GetByToken(ctx context.Context, token uuid.UUID) (*model.Wishlist, error) {
	return s.repo.GetByToken(ctx, token)
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
