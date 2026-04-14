package gift

import (
	"context"
	"errors"
	"testing"

	"wishlist-service/internal/model"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type giftRepoMock struct {
	mock.Mock
}

func (m *giftRepoMock) Save(ctx context.Context, gift *model.Gift) (*model.Gift, error) {
	args := m.Called(ctx, gift)
	if result := args.Get(0); result != nil {
		return result.(*model.Gift), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *giftRepoMock) Update(ctx context.Context, gift *model.Gift) (*model.Gift, error) {
	args := m.Called(ctx, gift)
	if result := args.Get(0); result != nil {
		return result.(*model.Gift), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *giftRepoMock) Book(ctx context.Context, id int64, token uuid.UUID) (*model.Gift, error) {
	args := m.Called(ctx, id, token)
	if result := args.Get(0); result != nil {
		return result.(*model.Gift), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *giftRepoMock) GetByID(ctx context.Context, id int64) (*model.Gift, error) {
	args := m.Called(ctx, id)
	if result := args.Get(0); result != nil {
		return result.(*model.Gift), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *giftRepoMock) GetByIDAndUserID(ctx context.Context, id int64, userID uuid.UUID) (*model.Gift, error) {
	args := m.Called(ctx, id, userID)
	if result := args.Get(0); result != nil {
		return result.(*model.Gift), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *giftRepoMock) GetByWishlistID(ctx context.Context, id int64) ([]model.Gift, error) {
	args := m.Called(ctx, id)
	if result := args.Get(0); result != nil {
		return result.([]model.Gift), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *giftRepoMock) Delete(ctx context.Context, giftID int64) (*model.Gift, error) {
	args := m.Called(ctx, giftID)
	if result := args.Get(0); result != nil {
		return result.(*model.Gift), args.Error(1)
	}
	return nil, args.Error(1)
}

type giftWishlistReaderMock struct {
	mock.Mock
}

func (m *giftWishlistReaderMock) GetByID(ctx context.Context, id int64) (*model.Wishlist, error) {
	args := m.Called(ctx, id)
	if result := args.Get(0); result != nil {
		return result.(*model.Wishlist), args.Error(1)
	}
	return nil, args.Error(1)
}

func (m *giftWishlistReaderMock) GetByToken(ctx context.Context, token uuid.UUID) (*model.Wishlist, error) {
	args := m.Called(ctx, token)
	if result := args.Get(0); result != nil {
		return result.(*model.Wishlist), args.Error(1)
	}
	return nil, args.Error(1)
}

func TestSave_ForeignWishlistForUser(t *testing.T) {
	t.Parallel()

	userID := uuid.New()
	ownerID := uuid.New()

	repo := &giftRepoMock{}
	wishlistReader := &giftWishlistReaderMock{}
	wishlistReader.On("GetByID", mock.Anything, int64(5)).Return(&model.Wishlist{ID: 5, UserID: ownerID}, nil)

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

	repo := &giftRepoMock{}
	repo.On("GetByID", mock.Anything, int64(22)).Return(&model.Gift{ID: 22, WishlistID: wishlistID + 1}, nil)
	wishlistReader := &giftWishlistReaderMock{}
	wishlistReader.On("GetByID", mock.Anything, wishlistID).Return(&model.Wishlist{ID: wishlistID, UserID: userID}, nil)

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

	repo := &giftRepoMock{}
	repo.On("GetByID", mock.Anything, int64(1)).Return(&model.Gift{ID: 1, WishlistID: 5}, nil)
	wishlistReader := &giftWishlistReaderMock{}
	wishlistReader.On("GetByID", mock.Anything, int64(5)).Return(&model.Wishlist{ID: 5, UserID: ownerID}, nil)

	svc := NewService(repo, wishlistReader)
	_, err := svc.Delete(context.Background(), userID, 5, 1)
	if !errors.Is(err, model.ErrForbidden) {
		t.Fatalf("Delete() error = %v, want %v", err, model.ErrForbidden)
	}
}

func Test_ServiceReturnsAlreadyBooked(t *testing.T) {
	t.Parallel()

	token := uuid.New()

	repo := &giftRepoMock{}
	repo.On("Book", mock.Anything, int64(1), token).Return((*model.Gift)(nil), model.ErrNotUpdated)
	repo.On("GetByID", mock.Anything, int64(1)).Return(&model.Gift{ID: 1, WishlistID: 5, Booked: true}, nil)
	wishlistReader := &giftWishlistReaderMock{}
	wishlistReader.On("GetByToken", mock.Anything, token).Return(&model.Wishlist{ID: 5, Token: token}, nil)

	svc := NewService(repo, wishlistReader)
	_, err := svc.Book(context.Background(), 1, token)
	if !errors.Is(err, model.ErrAlreadyBooked) {
		t.Fatalf("Book() error = %v, want %v", err, model.ErrAlreadyBooked)
	}
}
