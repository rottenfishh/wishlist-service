package dto

import (
	"time"
	"wishlist-service/internal/model"

	"github.com/google/uuid"
)

type CreateWishlistRequest struct {
	Title       string    `json:"title" binding:"required,min=1,max=255"`
	Description string    `json:"description"`
	Date        time.Time `json:"date" binding:"required"`
}

type UpdateWishlistRequest struct {
	Title       *string    `json:"title" binding:"omitempty,min=1,max=255"`
	Description *string    `json:"description"`
	Date        *time.Time `json:"date"`
}

type WishlistResponse struct {
	ID          int64     `json:"id"`
	Token       uuid.UUID `json:"token"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Date        time.Time `json:"date"`
}

type ListOfWishlistResponse struct {
	List []WishlistResponse `json:"list"`
}

type WishlistDetailsResponse struct {
	WishlistResponse
	Items []GiftResponse `json:"items"`
}

func ToWishlistResponse(wishlist model.Wishlist) WishlistResponse {
	return WishlistResponse{
		ID:          wishlist.ID,
		Token:       wishlist.Token,
		Title:       wishlist.Title,
		Description: wishlist.Description,
		Date:        wishlist.Date,
	}
}

func ToListOfWishlistResponse(wishlist []model.Wishlist) ListOfWishlistResponse {
	list := make([]WishlistResponse, 0)
	for _, wishlistItem := range wishlist {
		wishlistResp := ToWishlistResponse(wishlistItem)
		list = append(list, wishlistResp)
	}

	return ListOfWishlistResponse{List: list}
}

func ToWishListDetailsResponse(wishlist model.WishlistDetails) WishlistDetailsResponse {
	gifts := make([]GiftResponse, 0)
	for _, gift := range wishlist.Gifts {
		giftResp := ToGiftResponse(&gift)
		gifts = append(gifts, *giftResp)
	}

	return WishlistDetailsResponse{WishlistResponse: ToWishlistResponse(wishlist.Wishlist), Items: gifts}
}
