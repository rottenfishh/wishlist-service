package httpservice

import (
	"net/http"
	"strconv"
	"wishlist-service/internal/adapter/in/dto"
	"wishlist-service/internal/model"
	"wishlist-service/internal/service/gift"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GiftHandler struct {
	service gift.Service
}

func NewGiftHandler(service gift.Service) *GiftHandler {
	return &GiftHandler{
		service: service,
	}
}

// Create godoc
// @Summary Create gift
// @Description Creates a new gift in the specified wishlist owned by the authenticated user
// @Tags gifts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param wishlistId path int true "Wishlist ID"
// @Param request body dto.CreateGiftRequest true "Gift data"
// @Success 200 {object} dto.GiftResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/wishlists/{wishlistId}/gifts [post]
func (h *GiftHandler) Create(c *gin.Context) {
	wishlistIDParam := c.Param("wishlistId")
	wishlistID, err := strconv.ParseInt(wishlistIDParam, 10, 64)
	if err != nil {
		writeError(c, model.ErrInvalidRequest)
		return
	}

	userID, err := extractUserID(c)
	if err != nil {
		return
	}

	var req dto.CreateGiftRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		writeError(c, model.ErrInvalidRequest)
		return
	}

	newGift, err := h.service.Save(c.Request.Context(), userID, wishlistID, req.Name,
		req.Description, req.Link, req.Priority)
	if err != nil {
		writeError(c, err)
		return
	}

	resp := dto.ToGiftResponse(newGift)
	c.JSON(http.StatusOK, resp)
}

// Update godoc
// @Summary Update gift
// @Description Updates a gift in the specified wishlist owned by the authenticated user
// @Tags gifts
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param wishlistId path int true "Wishlist ID"
// @Param id path int true "Gift ID"
// @Param request body dto.UpdateGiftRequest true "Updated gift data"
// @Success 200 {object} dto.GiftResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/wishlists/{wishlistId}/gifts/{id} [put]
func (h *GiftHandler) Update(c *gin.Context) {
	wishlistIDParam := c.Param("wishlistId")
	wishlistID, err := strconv.Atoi(wishlistIDParam)
	if err != nil {
		writeError(c, model.ErrInvalidRequest)
		return
	}

	giftIDParam := c.Param("id")
	giftID, err := strconv.Atoi(giftIDParam)
	if err != nil {
		writeError(c, model.ErrInvalidRequest)
		return
	}
	userID, err := extractUserID(c)
	if err != nil {
		writeError(c, err)
		return
	}

	var req dto.UpdateGiftRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		writeError(c, model.ErrInvalidRequest)
		return
	}

	newGift, err := h.service.Update(c.Request.Context(), userID, int64(wishlistID),
		int64(giftID), req.Name, req.Description, req.Link, req.Priority)
	if err != nil {
		writeError(c, err)
		return
	}

	resp := dto.ToGiftResponse(newGift)
	c.JSON(http.StatusOK, resp)
}

// Delete godoc
// @Summary Delete gift
// @Description Deletes a gift from the specified wishlist owned by the authenticated user
// @Tags gifts
// @Produce json
// @Security BearerAuth
// @Param wishlistId path int true "Wishlist ID"
// @Param id path int true "Gift ID"
// @Success 200 {object} dto.GiftResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/wishlists/{wishlistId}/gifts/{id} [delete]
func (h *GiftHandler) Delete(c *gin.Context) {
	wishlistIDParam := c.Param("wishlistId")
	wishlistID, err := strconv.Atoi(wishlistIDParam)
	if err != nil {
		writeError(c, model.ErrInvalidRequest)
		return
	}

	giftIDParam := c.Param("id")
	giftID, err := strconv.Atoi(giftIDParam)
	if err != nil {
		writeError(c, model.ErrInvalidRequest)
		return
	}

	userID, err := extractUserID(c)
	if err != nil {
		writeError(c, err)
		return
	}

	deletedGift, err := h.service.Delete(c.Request.Context(), userID, int64(wishlistID), int64(giftID))
	if err != nil {
		writeError(c, err)
		return
	}

	resp := dto.ToGiftResponse(deletedGift)
	c.JSON(http.StatusOK, resp)
}

// Book godoc
// @Summary Book gift
// @Description Books a gift by public wishlist token
// @Tags gifts
// @Produce json
// @Param token path string true "Wishlist public token (UUID)"
// @Param id path int true "Gift ID"
// @Success 200 {object} dto.GiftResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/public/wishlists/{token}/gifts/{id} [post]
func (h *GiftHandler) Book(c *gin.Context) {
	tokenParam := c.Param("token")
	token, err := uuid.Parse(tokenParam)
	if err != nil {
		writeError(c, model.ErrInvalidRequest)
		return
	}

	giftIDParam := c.Param("id")
	giftID, err := strconv.Atoi(giftIDParam)
	if err != nil {
		writeError(c, model.ErrInvalidRequest)
		return
	}

	bookedGift, err := h.service.Book(c.Request.Context(), int64(giftID), token)
	if err != nil {
		writeError(c, err)
		return
	}

	resp := dto.ToGiftResponse(bookedGift)
	c.JSON(http.StatusOK, resp)
}

func (h *GiftHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	router.POST("/:wishlistId/gifts", authMiddleware, h.Create)
	router.PUT("/:wishlistId/gifts/:id", authMiddleware, h.Update)
	router.DELETE("/:wishlistId/gifts/:id", authMiddleware, h.Delete)
}
