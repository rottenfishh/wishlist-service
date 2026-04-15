package gift

import (
	"context"
	"errors"
	"wishlist-service/internal/model"

	"github.com/google/uuid"
)

type service struct {
	repo     Repository
	wishlist WishlistReader
}

func (s *service) GetByID(ctx context.Context, userID uuid.UUID, wishlistID, ID int64) (*model.Gift, error) {
	whList, err := s.wishlist.GetByID(ctx, wishlistID)
	if err != nil {
		return nil, err
	}

	if whList.UserID != userID {
		return nil, model.ErrForbidden
	}

	gift, err := s.repo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	if gift.WishlistID != wishlistID {
		return nil, model.ErrNotFound
	}

	return gift, nil
}

func NewService(repo Repository, wishlist WishlistReader) Service {
	return &service{repo: repo, wishlist: wishlist}
}

func (s *service) Save(ctx context.Context, userID uuid.UUID, wishlistID int64, name, description, link string, priority int) (*model.Gift, error) {
	whList, err := s.wishlist.GetByID(ctx, wishlistID)
	if err != nil {
		return nil, err
	}

	if whList.UserID != userID {
		return nil, model.ErrForbidden
	}

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

func (s *service) Update(ctx context.Context, userID uuid.UUID, wishlistID, ID int64, name, description,
	link *string, priority *int) (*model.Gift, error) {

	wishlist, err := s.wishlist.GetByID(ctx, wishlistID)
	if err != nil {
		return nil, err
	}
	if wishlist.UserID != userID {
		return nil, model.ErrForbidden
	}

	gift, err := s.repo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	if gift.WishlistID != wishlistID {
		return nil, model.ErrNotFound
	}

	if name != nil {
		gift.Name = *name
	}
	if description != nil {
		gift.Description = *description
	}
	if link != nil {
		gift.Link = *link
	}
	if priority != nil {
		gift.Priority = *priority
	}

	return s.repo.Update(ctx, gift)
}

func (s *service) Book(ctx context.Context, ID int64, token uuid.UUID) (*model.Gift, error) {
	gift, err := s.repo.Book(ctx, ID, token)
	if err == nil {
		return gift, nil
	}

	return nil, s.resolveBookingError(ctx, ID, token, err)
}

func (s *service) Delete(ctx context.Context, userID uuid.UUID, wishlistID, ID int64) (*model.Gift, error) {
	wishlist, err := s.wishlist.GetByID(ctx, wishlistID)
	if err != nil {
		return nil, err
	}
	if wishlist.UserID != userID {
		return nil, model.ErrForbidden
	}

	gift, err := s.repo.GetByID(ctx, ID)
	if err != nil {
		return nil, err
	}
	if gift.WishlistID != wishlistID {
		return nil, model.ErrNotFound
	}

	return s.repo.Delete(ctx, ID)
}

func (s *service) resolveBookingError(ctx context.Context, id int64, token uuid.UUID, err error) error {
	if !errors.Is(err, model.ErrNotUpdated) {
		return err
	}

	if _, err = s.wishlist.GetByToken(ctx, token); err != nil {
		return err
	}

	if _, err = s.repo.GetByID(ctx, id); err != nil {
		return err
	}

	return model.ErrAlreadyBooked
}
