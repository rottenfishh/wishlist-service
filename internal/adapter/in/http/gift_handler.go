package http

import (
	"cdek/internal/adapter/in/dto"
	"cdek/internal/model"
	"cdek/internal/service/gift"
	"net/http"
	"strconv"

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
	router.POST("", authMiddleware, h.Create)
	router.PUT("/:id", authMiddleware, h.Update)
	router.DELETE("/:id", authMiddleware, h.Delete)
}
