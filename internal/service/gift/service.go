package gift

import (
	"cdek/internal/model"
	"context"

	"github.com/google/uuid"
)

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) Save(ctx context.Context, wishlistID int64, name, description, link string, priority int) (*model.Gift, error) {
	gift := &model.Gift{
		WishlistID:  wishlistID,
		Name:        name,
		Description: description,
		Link:        link,
		Priority:    priority,
		Booked:      false,
	}

	return s.repo.Save(ctx, gift)
}

func (s *service) Update(ctx context.Context, ID int64, name, description, link string, priority int) (*model.Gift, error) {
	gift, err := s.repo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	gift.Name = name
	gift.Description = description
	gift.Link = link
	gift.Priority = priority

	return s.repo.Update(ctx, gift)
}

func (s *service) Book(ctx context.Context, ID int64, token uuid.UUID) (*model.Gift, error) {
	return s.repo.Book(ctx, ID, token)
}

func (s *service) Delete(ctx context.Context, ID int64) (*model.Gift, error) {
	return s.repo.Delete(ctx, ID)
}
