package model

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	CreatedAt    time.Time `json:"createdAt,omitempty"`
	PasswordHash string    `json:"-"`
}

func NewUser(id uuid.UUID, email string, passwordHash string) *User {
	return &User{id, email, time.Now(), passwordHash}
}
