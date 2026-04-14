package wishlist

import (
	"context"
	"errors"
	"testing"
	"time"

	"wishlist-service/internal/model"

	"github.com/google/uuid"
)

type wishlistRepoMock struct {
	saveFn        func(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error)
	updateFn      func(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error)
	deleteFn      func(ctx context.Context, id int64) (*model.Wishlist, error)
	getByIDFn     func(ctx context.Context, id int64) (*model.Wishlist, error)
	getByUserIDFn func(ctx context.Context, userID uuid.UUID) ([]model.Wishlist, error)
	getByTokenFn  func(ctx context.Context, token uuid.UUID) (*model.Wishlist, error)
}

func (m *wishlistRepoMock) Save(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error) {
	return m.saveFn(ctx, wishlist)
}

func (m *wishlistRepoMock) Update(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error) {
	return m.updateFn(ctx, wishlist)
}

func (m *wishlistRepoMock) Delete(ctx context.Context, id int64) (*model.Wishlist, error) {
	return m.deleteFn(ctx, id)
}

func (m *wishlistRepoMock) GetByID(ctx context.Context, id int64) (*model.Wishlist, error) {
	return m.getByIDFn(ctx, id)
}

func (m *wishlistRepoMock) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Wishlist, error) {
	return m.getByUserIDFn(ctx, userID)
}

func (m *wishlistRepoMock) GetByToken(ctx context.Context, token uuid.UUID) (*model.Wishlist, error) {
	return m.getByTokenFn(ctx, token)
}

type wishlistGiftReaderMock struct {
	getByWishlistIDFn func(ctx context.Context, id int64) ([]model.Gift, error)
}

func (m *wishlistGiftReaderMock) GetByWishlistID(ctx context.Context, id int64) ([]model.Gift, error) {
	return m.getByWishlistIDFn(ctx, id)
}

func Test_ServiceCreate(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	date := time.Now().UTC()
	var saved *model.Wishlist

	repo := &wishlistRepoMock{
		saveFn: func(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error) {
			saved = wishlist
			return wishlist, nil
		},
		updateFn:      func(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error) { return nil, nil },
		deleteFn:      func(ctx context.Context, id int64) (*model.Wishlist, error) { return nil, nil },
		getByIDFn:     func(ctx context.Context, id int64) (*model.Wishlist, error) { return nil, nil },
		getByUserIDFn: func(ctx context.Context, userID uuid.UUID) ([]model.Wishlist, error) { return nil, nil },
		getByTokenFn:  func(ctx context.Context, token uuid.UUID) (*model.Wishlist, error) { return nil, nil },
	}
	gifts := &wishlistGiftReaderMock{
		getByWishlistIDFn: func(ctx context.Context, id int64) ([]model.Gift, error) { return nil, nil },
	}

	svc := NewService(repo, gifts)
	result, err := svc.Create(context.Background(), userID, "Birthday", "desc", date)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	if result == nil || saved == nil {
		t.Fatal("Create() did not save wishlist")
	}
	if saved.UserID != userID {
		t.Fatalf("saved.UserID = %v, want %v", saved.UserID, userID)
	}
	if saved.Token == uuid.Nil {
		t.Fatal("saved.Token is empty")
	}
}

func Test_UpdateForeignWishlist(t *testing.T) {
	t.Parallel()

	ownerID := uuid.New()
	otherUserID := uuid.New()
	existing := &model.Wishlist{ID: 10, UserID: ownerID, Title: "old"}

	repo := &wishlistRepoMock{
		saveFn: func(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error) { return nil, nil },
		updateFn: func(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error) {
			t.Fatal("Update() should not reach repository update")
			return nil, nil
		},
		deleteFn:      func(ctx context.Context, id int64) (*model.Wishlist, error) { return nil, nil },
		getByIDFn:     func(ctx context.Context, id int64) (*model.Wishlist, error) { return existing, nil },
		getByUserIDFn: func(ctx context.Context, userID uuid.UUID) ([]model.Wishlist, error) { return nil, nil },
		getByTokenFn:  func(ctx context.Context, token uuid.UUID) (*model.Wishlist, error) { return nil, nil },
	}
	gifts := &wishlistGiftReaderMock{
		getByWishlistIDFn: func(ctx context.Context, id int64) ([]model.Gift, error) { return nil, nil },
	}

	svc := NewService(repo, gifts)
	title := "new"
	_, err := svc.Update(context.Background(), otherUserID, existing.ID, &title, nil, nil)
	if !errors.Is(err, model.ErrForbidden) {
		t.Fatalf("Update() error = %v, want %v", err, model.ErrForbidden)
	}
}

func Test_ReturnWishlistDetails(t *testing.T) {
	t.Parallel()

	wishlistID := int64(7)
	expectedWishlist := &model.Wishlist{ID: wishlistID, UserID: uuid.New(), Title: "Birthday"}
	expectedGifts := []model.Gift{{ID: 1, WishlistID: wishlistID, Name: "Book"}}

	repo := &wishlistRepoMock{
		saveFn:        func(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error) { return nil, nil },
		updateFn:      func(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error) { return nil, nil },
		deleteFn:      func(ctx context.Context, id int64) (*model.Wishlist, error) { return nil, nil },
		getByIDFn:     func(ctx context.Context, id int64) (*model.Wishlist, error) { return expectedWishlist, nil },
		getByUserIDFn: func(ctx context.Context, userID uuid.UUID) ([]model.Wishlist, error) { return nil, nil },
		getByTokenFn:  func(ctx context.Context, token uuid.UUID) (*model.Wishlist, error) { return nil, nil },
	}
	gifts := &wishlistGiftReaderMock{
		getByWishlistIDFn: func(ctx context.Context, id int64) ([]model.Gift, error) { return expectedGifts, nil },
	}

	svc := NewService(repo, gifts)
	result, err := svc.GetByID(context.Background(), wishlistID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if result.Wishlist.ID != wishlistID {
		t.Fatalf("result.Wishlist.ID = %d, want %d", result.Wishlist.ID, wishlistID)
	}
	if len(result.Gifts) != 1 || result.Gifts[0].Name != "Book" {
		t.Fatalf("result.Gifts = %#v, want one Book gift", result.Gifts)
	}
}

func Test_GiftReaderError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("gift reader failed")
	token := uuid.New()

	repo := &wishlistRepoMock{
		saveFn:        func(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error) { return nil, nil },
		updateFn:      func(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error) { return nil, nil },
		deleteFn:      func(ctx context.Context, id int64) (*model.Wishlist, error) { return nil, nil },
		getByIDFn:     func(ctx context.Context, id int64) (*model.Wishlist, error) { return nil, nil },
		getByUserIDFn: func(ctx context.Context, userID uuid.UUID) ([]model.Wishlist, error) { return nil, nil },
		getByTokenFn: func(ctx context.Context, token uuid.UUID) (*model.Wishlist, error) {
			return &model.Wishlist{ID: 11, UserID: uuid.New()}, nil
		},
	}
	gifts := &wishlistGiftReaderMock{
		getByWishlistIDFn: func(ctx context.Context, id int64) ([]model.Gift, error) { return nil, expectedErr },
	}

	svc := NewService(repo, gifts)
	_, err := svc.GetByToken(context.Background(), token)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("GetByToken() error = %v, want %v", err, expectedErr)
	}
}
