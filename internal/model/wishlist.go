package model

import (
	"time"

	"github.com/google/uuid"
)

type Wishlist struct {
	ID          int64
	UserID      uuid.UUID
	Token       uuid.UUID
	Title       string
	Description string
	Date        time.Time
}
