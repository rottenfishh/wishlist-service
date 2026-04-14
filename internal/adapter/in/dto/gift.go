package dto

import "wishlist-service/internal/model"

type CreateGiftRequest struct {
	Name        string `json:"name" binding:"required,min=1,max=255"`
	Description string `json:"description"`
	Link        string `json:"link" binding:"omitempty,url"`
	Priority    int    `json:"priority" binding:"required,gte=1,lte=5"`
}

type UpdateGiftRequest struct {
	Name        *string `json:"name" binding:"omitempty,min=1,max=255"`
	Description *string `json:"description"`
	Link        *string `json:"link" binding:"omitempty,url"`
	Priority    *int    `json:"priority" binding:"omitempty,gte=1,lte=5"`
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
