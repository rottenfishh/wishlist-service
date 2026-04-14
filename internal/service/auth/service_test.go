package auth

import (
	"context"
	"errors"
	"testing"

	"wishlist-service/internal/model"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type authRepoMock struct {
	saveUserFn       func(ctx context.Context, user *model.User) (*model.User, error)
	getUserByIDFn    func(ctx context.Context, id string) (*model.User, error)
	getUserByEmailFn func(ctx context.Context, email string) (*model.User, error)
}

func (m *authRepoMock) SaveUser(ctx context.Context, user *model.User) (*model.User, error) {
	return m.saveUserFn(ctx, user)
}

func (m *authRepoMock) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return m.getUserByIDFn(ctx, id)
}

func (m *authRepoMock) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	return m.getUserByEmailFn(ctx, email)
}

func TestRegister(t *testing.T) {
	t.Parallel()

	var savedInput *model.User
	repo := &authRepoMock{
		saveUserFn: func(ctx context.Context, user *model.User) (*model.User, error) {
			savedInput = user
			return user, nil
		},
		getUserByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return nil, nil
		},
		getUserByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
			return nil, nil
		},
	}

	svc := NewUserService(repo, "secret", 10000)
	user, err := svc.Register(context.Background(), "test@example.com", "plain-password")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	if user == nil {
		t.Fatal("Register() returned nil user")
	}
	if savedInput == nil {
		t.Fatal("SaveUser() was not called")
	}
	if savedInput.PasswordHash == "" {
		t.Fatal("password hash is empty")
	}
	if savedInput.PasswordHash == "plain-password" {
		t.Fatal("password was not hashed")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(savedInput.PasswordHash), []byte("plain-password")); err != nil {
		t.Fatalf("saved hash does not match password: %v", err)
	}
}

func Test_RegisterReturnsRepositoryError(t *testing.T) {
	t.Parallel()

	expectedErr := model.ErrUserAlreadyExists
	repo := &authRepoMock{
		saveUserFn: func(ctx context.Context, user *model.User) (*model.User, error) {
			return nil, expectedErr
		},
		getUserByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return nil, nil
		},
		getUserByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
			return nil, nil
		},
	}

	svc := NewUserService(repo, "secret", 10000)
	_, err := svc.Register(context.Background(), "test@example.com", "plain-password")
	if !errors.Is(err, expectedErr) {
		t.Fatalf("Register() error = %v, want %v", err, expectedErr)
	}
}

func TestLoginSuccess(t *testing.T) {
	t.Parallel()

	hash, err := bcrypt.GenerateFromPassword([]byte("plain-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	userID := uuid.New()
	repo := &authRepoMock{
		saveUserFn: func(ctx context.Context, user *model.User) (*model.User, error) {
			return nil, nil
		},
		getUserByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return nil, nil
		},
		getUserByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
			return &model.User{
				ID:           userID,
				Email:        email,
				PasswordHash: string(hash),
			}, nil
		},
	}

	svc := NewUserService(repo, "secret", 10000)
	token, err := svc.Login(context.Background(), "test@example.com", "plain-password")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}

	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		t.Fatalf("jwt.Parse() error = %v", err)
	}
	if !parsedToken.Valid {
		t.Fatal("generated token is invalid")
	}
}

func Test_LoginReturnsUnauthorized(t *testing.T) {
	t.Parallel()

	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("GenerateFromPassword() error = %v", err)
	}

	repo := &authRepoMock{
		saveUserFn: func(ctx context.Context, user *model.User) (*model.User, error) {
			return nil, nil
		},
		getUserByIDFn: func(ctx context.Context, id string) (*model.User, error) {
			return nil, nil
		},
		getUserByEmailFn: func(ctx context.Context, email string) (*model.User, error) {
			return &model.User{
				ID:           uuid.New(),
				Email:        email,
				PasswordHash: string(hash),
			}, nil
		},
	}

	svc := NewUserService(repo, "secret", 10000)
	_, err = svc.Login(context.Background(), "test@example.com", "wrong-password")
	if !errors.Is(err, model.ErrUnauthorized) {
		t.Fatalf("Login() error = %v, want %v", err, model.ErrUnauthorized)
	}
}
