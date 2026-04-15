package wishlist

import (
	"context"
	"errors"
	"testing"
	"time"

	"wishlist-service/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type wishlistRepoMock struct {
	mock.Mock
}

func (m *wishlistRepoMock) Save(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error) {
	args := m.Called(ctx, wishlist)
	if result := args.Get(0); result != nil {
		return result.(*model.Wishlist), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *wishlistRepoMock) Update(ctx context.Context, wishlist *model.Wishlist) (*model.Wishlist, error) {
	args := m.Called(ctx, wishlist)
	if result := args.Get(0); result != nil {
		return result.(*model.Wishlist), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *wishlistRepoMock) Delete(ctx context.Context, id int64) (*model.Wishlist, error) {
	args := m.Called(ctx, id)
	if result := args.Get(0); result != nil {
		return result.(*model.Wishlist), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *wishlistRepoMock) GetByID(ctx context.Context, id int64) (*model.Wishlist, error) {
	args := m.Called(ctx, id)
	if result := args.Get(0); result != nil {
		return result.(*model.Wishlist), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *wishlistRepoMock) GetByUserID(ctx context.Context, userID uuid.UUID) ([]model.Wishlist, error) {
	args := m.Called(ctx, userID)
	if result := args.Get(0); result != nil {
		return result.([]model.Wishlist), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *wishlistRepoMock) GetByToken(ctx context.Context, token uuid.UUID) (*model.Wishlist, error) {
	args := m.Called(ctx, token)
	if result := args.Get(0); result != nil {
		return result.(*model.Wishlist), args.Error(1)
	}
	return nil, args.Error(1)
}

type wishlistGiftReaderMock struct {
	mock.Mock
}

func (m *wishlistGiftReaderMock) GetByWishlistID(ctx context.Context, id int64) ([]model.Gift, error) {
	args := m.Called(ctx, id)
	if result := args.Get(0); result != nil {
		return result.([]model.Gift), args.Error(1)
	}
	return nil, args.Error(1)
}

func Test_ServiceCreate(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	date := time.Now().UTC()
	var saved *model.Wishlist

	repo := &wishlistRepoMock{}
	repo.On("Save", mock.Anything, mock.AnythingOfType("*model.Wishlist")).Run(func(args mock.Arguments) {
		saved = args.Get(1).(*model.Wishlist)
	}).Return(&model.Wishlist{}, nil)
	gifts := &wishlistGiftReaderMock{}

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

	repo := &wishlistRepoMock{}
	repo.On("GetByID", mock.Anything, existing.ID).Return(existing, nil)
	gifts := &wishlistGiftReaderMock{}

	svc := NewService(repo, gifts)
	title := "new"
	_, err := svc.Update(context.Background(), otherUserID, existing.ID, &title, nil, nil)
	if !errors.Is(err, model.ErrForbidden) {
		t.Fatalf("Update() error = %v, want %v", err, model.ErrForbidden)
	}
}

func Test_ReturnWishlistDetails(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	wishlistID := int64(7)
	expectedWishlist := &model.Wishlist{ID: wishlistID, UserID: userID, Title: "Birthday"}
	expectedGifts := []model.Gift{{ID: 1, WishlistID: wishlistID, Name: "Book"}}

	repo := &wishlistRepoMock{}
	repo.On("GetByID", mock.Anything, wishlistID).Return(expectedWishlist, nil)
	gifts := &wishlistGiftReaderMock{}
	gifts.On("GetByWishlistID", mock.Anything, wishlistID).Return(expectedGifts, nil)

	svc := NewService(repo, gifts)
	result, err := svc.GetByID(context.Background(), userID, wishlistID)
	if err != nil {
		t.Fatalf("GetByID() error = %v", err)
	}

	if result.ID != wishlistID {
		t.Fatalf("result.Wishlist.ID = %d, want %d", result.ID, wishlistID)
	}
	if len(result.Gifts) != 1 || result.Gifts[0].Name != "Book" {
		t.Fatalf("result.Gifts = %#v, want one Book gift", result.Gifts)
	}
}

func Test_GiftReaderError(t *testing.T) {
	t.Parallel()

	expectedErr := errors.New("gift reader failed")
	token := uuid.New()

	repo := &wishlistRepoMock{}
	repo.On("GetByToken", mock.Anything, token).Return(&model.Wishlist{ID: 11, UserID: uuid.New()}, nil)
	gifts := &wishlistGiftReaderMock{}
	gifts.On("GetByWishlistID", mock.Anything, int64(11)).Return(([]model.Gift)(nil), expectedErr)

	svc := NewService(repo, gifts)
	_, err := svc.GetByToken(context.Background(), token)
	if !errors.Is(err, expectedErr) {
		t.Fatalf("GetByToken() error = %v, want %v", err, expectedErr)
	}
}
