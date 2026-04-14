package gift

import (
	"context"
	"errors"
	"testing"

	"wishlist-service/internal/model"

	"github.com/google/uuid"
)

type giftRepoMock struct {
	saveFn          func(ctx context.Context, gift *model.Gift) (*model.Gift, error)
	updateFn        func(ctx context.Context, gift *model.Gift) (*model.Gift, error)
	bookFn          func(ctx context.Context, id int64, token uuid.UUID) (*model.Gift, error)
	getByIDFn       func(ctx context.Context, id int64) (*model.Gift, error)
	getByIDUserFn   func(ctx context.Context, id int64, userID uuid.UUID) (*model.Gift, error)
	getByWishlistFn func(ctx context.Context, id int64) ([]model.Gift, error)
	deleteFn        func(ctx context.Context, giftID int64) (*model.Gift, error)
}

func (m *giftRepoMock) Save(ctx context.Context, gift *model.Gift) (*model.Gift, error) {
	return m.saveFn(ctx, gift)
}

func (m *giftRepoMock) Update(ctx context.Context, gift *model.Gift) (*model.Gift, error) {
	return m.updateFn(ctx, gift)
}

func (m *giftRepoMock) Book(ctx context.Context, id int64, token uuid.UUID) (*model.Gift, error) {
	return m.bookFn(ctx, id, token)
}

func (m *giftRepoMock) GetByID(ctx context.Context, id int64) (*model.Gift, error) {
	return m.getByIDFn(ctx, id)
}

func (m *giftRepoMock) GetByIDAndUserID(ctx context.Context, id int64, userID uuid.UUID) (*model.Gift, error) {
	return m.getByIDUserFn(ctx, id, userID)
}

func (m *giftRepoMock) GetByWishlistID(ctx context.Context, id int64) ([]model.Gift, error) {
	return m.getByWishlistFn(ctx, id)
}

func (m *giftRepoMock) Delete(ctx context.Context, giftID int64) (*model.Gift, error) {
	return m.deleteFn(ctx, giftID)
}

type giftWishlistReaderMock struct {
	getByIDFn    func(ctx context.Context, id int64) (*model.Wishlist, error)
	getByTokenFn func(ctx context.Context, token uuid.UUID) (*model.Wishlist, error)
}

func (m *giftWishlistReaderMock) GetByID(ctx context.Context, id int64) (*model.Wishlist, error) {
	return m.getByIDFn(ctx, id)
}

func (m *giftWishlistReaderMock) GetByToken(ctx context.Context, token uuid.UUID) (*model.Wishlist, error) {
	return m.getByTokenFn(ctx, token)
}

func TestSave_ForeignWishlistForUser(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	ownerID := uuid.New()

	repo := &giftRepoMock{
		saveFn: func(ctx context.Context, gift *model.Gift) (*model.Gift, error) {
			t.Fatal("Save() should not reach repository save")
			return nil, nil
		},
		updateFn:        func(ctx context.Context, gift *model.Gift) (*model.Gift, error) { return nil, nil },
		bookFn:          func(ctx context.Context, id int64, token uuid.UUID) (*model.Gift, error) { return nil, nil },
		getByIDFn:       func(ctx context.Context, id int64) (*model.Gift, error) { return nil, nil },
		getByIDUserFn:   func(ctx context.Context, id int64, userID uuid.UUID) (*model.Gift, error) { return nil, nil },
		getByWishlistFn: func(ctx context.Context, id int64) ([]model.Gift, error) { return nil, nil },
		deleteFn:        func(ctx context.Context, giftID int64) (*model.Gift, error) { return nil, nil },
	}
	wishlistReader := &giftWishlistReaderMock{
		getByIDFn: func(ctx context.Context, id int64) (*model.Wishlist, error) {
			return &model.Wishlist{ID: id, UserID: ownerID}, nil
		},
		getByTokenFn: func(ctx context.Context, token uuid.UUID) (*model.Wishlist, error) { return nil, nil },
	}

	svc := NewService(repo, wishlistReader)
	_, err := svc.Save(context.Background(), userID, 5, "Gift", "desc", "link", 3)
	if !errors.Is(err, model.ErrForbidden) {
		t.Fatalf("Save() error = %v, want %v", err, model.ErrForbidden)
	}
}

func Test_ForeignWishlistNotfound(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	wishlistID := int64(10)

	repo := &giftRepoMock{
		saveFn: func(ctx context.Context, gift *model.Gift) (*model.Gift, error) { return nil, nil },
		updateFn: func(ctx context.Context, gift *model.Gift) (*model.Gift, error) {
			t.Fatal("Update() should not reach repository update")
			return nil, nil
		},
		bookFn: func(ctx context.Context, id int64, token uuid.UUID) (*model.Gift, error) { return nil, nil },
		getByIDFn: func(ctx context.Context, id int64) (*model.Gift, error) {
			return &model.Gift{ID: id, WishlistID: wishlistID + 1}, nil
		},
		getByIDUserFn:   func(ctx context.Context, id int64, userID uuid.UUID) (*model.Gift, error) { return nil, nil },
		getByWishlistFn: func(ctx context.Context, id int64) ([]model.Gift, error) { return nil, nil },
		deleteFn:        func(ctx context.Context, giftID int64) (*model.Gift, error) { return nil, nil },
	}
	wishlistReader := &giftWishlistReaderMock{
		getByIDFn: func(ctx context.Context, id int64) (*model.Wishlist, error) {
			return &model.Wishlist{ID: id, UserID: userID}, nil
		},
		getByTokenFn: func(ctx context.Context, token uuid.UUID) (*model.Wishlist, error) { return nil, nil },
	}

	svc := NewService(repo, wishlistReader)
	name := "updated"
	_, err := svc.Update(context.Background(), userID, wishlistID, 22, &name, nil, nil, nil)
	if !errors.Is(err, model.ErrNotFound) {
		t.Fatalf("Update() error = %v, want %v", err, model.ErrNotFound)
	}
}

func Test_DeleteForeignWishlistForUser(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	ownerID := uuid.New()

	repo := &giftRepoMock{
		saveFn:   func(ctx context.Context, gift *model.Gift) (*model.Gift, error) { return nil, nil },
		updateFn: func(ctx context.Context, gift *model.Gift) (*model.Gift, error) { return nil, nil },
		bookFn:   func(ctx context.Context, id int64, token uuid.UUID) (*model.Gift, error) { return nil, nil },
		getByIDFn: func(ctx context.Context, id int64) (*model.Gift, error) {
			return &model.Gift{ID: id, WishlistID: 5}, nil
		},
		getByIDUserFn:   func(ctx context.Context, id int64, userID uuid.UUID) (*model.Gift, error) { return nil, nil },
		getByWishlistFn: func(ctx context.Context, id int64) ([]model.Gift, error) { return nil, nil },
		deleteFn: func(ctx context.Context, giftID int64) (*model.Gift, error) {
			t.Fatal("Delete() should not reach repository delete")
			return nil, nil
		},
	}
	wishlistReader := &giftWishlistReaderMock{
		getByIDFn: func(ctx context.Context, id int64) (*model.Wishlist, error) {
			return &model.Wishlist{ID: id, UserID: ownerID}, nil
		},
		getByTokenFn: func(ctx context.Context, token uuid.UUID) (*model.Wishlist, error) { return nil, nil },
	}

	svc := NewService(repo, wishlistReader)
	_, err := svc.Delete(context.Background(), userID, 5, 1)
	if !errors.Is(err, model.ErrForbidden) {
		t.Fatalf("Delete() error = %v, want %v", err, model.ErrForbidden)
	}
}

func Test_ServiceReturnsAlreadyBooked(t *testing.T) {
	t.Parallel()

	token := uuid.New()

	repo := &giftRepoMock{
		saveFn:   func(ctx context.Context, gift *model.Gift) (*model.Gift, error) { return nil, nil },
		updateFn: func(ctx context.Context, gift *model.Gift) (*model.Gift, error) { return nil, nil },
		bookFn: func(ctx context.Context, id int64, token uuid.UUID) (*model.Gift, error) {
			return nil, model.ErrNotUpdated
		},
		getByIDFn: func(ctx context.Context, id int64) (*model.Gift, error) {
			return &model.Gift{ID: id, WishlistID: 5, Booked: true}, nil
		},
		getByIDUserFn:   func(ctx context.Context, id int64, userID uuid.UUID) (*model.Gift, error) { return nil, nil },
		getByWishlistFn: func(ctx context.Context, id int64) ([]model.Gift, error) { return nil, nil },
		deleteFn:        func(ctx context.Context, giftID int64) (*model.Gift, error) { return nil, nil },
	}
	wishlistReader := &giftWishlistReaderMock{
		getByIDFn: func(ctx context.Context, id int64) (*model.Wishlist, error) { return nil, nil },
		getByTokenFn: func(ctx context.Context, token uuid.UUID) (*model.Wishlist, error) {
			return &model.Wishlist{ID: 5, Token: token}, nil
		},
	}

	svc := NewService(repo, wishlistReader)
	_, err := svc.Book(context.Background(), 1, token)
	if !errors.Is(err, model.ErrAlreadyBooked) {
		t.Fatalf("Book() error = %v, want %v", err, model.ErrAlreadyBooked)
	}
}
