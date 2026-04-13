package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"
	"wishlist-service/internal/model"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	repo   Repository
	config Config
}

func NewUserService(repo Repository, JWTSecret string) Service {
	return &service{repo: repo, config: Config{JWTSecret: JWTSecret, JWTExpiresInSec: 10000}}
}

type Config struct {
	JWTSecret       string
	JWTExpiresInSec int64
}

func (s *service) Register(ctx context.Context, email, password string) (*model.User, error) {

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := model.NewUser(uuid.New(), email, string(hashed))
	saveUser, err := s.repo.SaveUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return saveUser, nil
}

func (s *service) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		slog.Error("get user by email error", "error", err)
		if errors.Is(err, model.ErrNotFound) {
			return "", model.ErrUnauthorized
		}
		return "", fmt.Errorf("get user by email repo: %w", err)
	}
	if bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		slog.Error("wrong password")
		return "", model.ErrUnauthorized
	}

	token, err := s.generateToken(user)
	if err != nil {
		slog.Error("generate token", "error", err)
		return "", err
	}

	return token, nil

}

func (s *service) generateToken(user *model.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":   user.ID.String(),
		"email": user.Email,
		"exp":   time.Now().Add(time.Duration(s.config.JWTExpiresInSec) * time.Second).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(s.config.JWTSecret))
}
