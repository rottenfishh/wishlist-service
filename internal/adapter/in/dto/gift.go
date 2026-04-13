package dto

import "cdek/internal/model"

type CreateGiftRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Priority    int    `json:"priority"`
}

type UpdateGiftRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Link        *string `json:"link"`
	Priority    *int    `json:"priority"`
}

type GiftResponse struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Link        string `json:"link"`
	Priority    int    `json:"priority"`
	Booked      bool   `json:"booked"`
}

func ToGiftResponse(gift *model.Gift) *GiftResponse {
	return &GiftResponse{
		ID:          int64(gift.ID),
		Name:        gift.Name,
		Description: gift.Description,
		Link:        gift.Link,
		Priority:    gift.Priority,
		Booked:      gift.Booked,
	}
}
