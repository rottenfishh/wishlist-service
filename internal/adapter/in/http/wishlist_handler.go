package http

import (
	"cdek/internal/adapter/in/dto"
	"cdek/internal/model"
	"cdek/internal/service/wishlist"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WishlistHandler struct {
	service wishlist.Service
}

func NewWishlistHandler(service wishlist.Service) *WishlistHandler {
	return &WishlistHandler{service: service}
}

func (h *WishlistHandler) Create(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		writeError(c, err)
		return
	}

	var req dto.CreateWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, model.ErrInvalidRequest)
		return
	}

	newWishlist, err := h.service.Create(c.Request.Context(), userID, req.Title, req.Description, req.Date)
	if err != nil {
		writeError(c, err)
		return
	}

	resp := dto.ToWishlistResponse(*newWishlist)
	c.JSON(http.StatusCreated, resp)
}

func (h *WishlistHandler) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		writeError(c, model.ErrInvalidRequest)
		return
	}

	userID, err := extractUserID(c)
	if err != nil {
		writeError(c, err)
		return
	}

	var req dto.UpdateWishlistRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeError(c, model.ErrInvalidRequest)
		return
	}

	newWishlist, err := h.service.Update(c.Request.Context(), userID, int64(id),
		req.Title, req.Description, req.Date)
	if err != nil {
		writeError(c, err)
		return
	}

	resp := dto.ToWishlistResponse(*newWishlist)
	c.JSON(http.StatusOK, resp)
}

func (h *WishlistHandler) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		writeError(c, model.ErrInvalidRequest)
		return
	}

	userID, err := extractUserID(c)
	if err != nil {
		writeError(c, err)
		return
	}

	deletedWishlist, err := h.service.Delete(c.Request.Context(), userID, int64(id))
	if err != nil {
		writeError(c, err)
		return
	}

	resp := dto.ToWishlistResponse(*deletedWishlist)
	c.JSON(http.StatusOK, resp)
}

func (h *WishlistHandler) List(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		writeError(c, err)
		return
	}

	list, err := h.service.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		writeError(c, err)
		return
	}

	resp := dto.ToListOfWishlistResponse(list)
	c.JSON(http.StatusOK, resp)
}

func (h *WishlistHandler) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		writeError(c, model.ErrInvalidRequest)
		return
	}

	userID, err := extractUserID(c)
	if err != nil {
		writeError(c, err)
		return
	}

	list, err := h.service.GetByID(c.Request.Context(), int64(id))
	if err != nil {
		writeError(c, err)
		return
	}

	if list.UserID != userID {
		writeError(c, model.ErrForbidden)
		return
	}

	resp := dto.ToWishListDetailsResponse(*list)
	c.JSON(http.StatusOK, resp)
}

func (h *WishlistHandler) GetByToken(c *gin.Context) {
	tokenParam := c.Param("token")
	token, err := uuid.Parse(tokenParam)
	if err != nil {
		writeError(c, model.ErrInvalidRequest)
		return
	}

	list, err := h.service.GetByToken(c.Request.Context(), token)
	if err != nil {
		writeError(c, err)
		return
	}

	resp := dto.ToWishListDetailsResponse(*list)
	c.JSON(http.StatusOK, resp)
}

func (h *WishlistHandler) RegisterRoutes(router *gin.RouterGroup, authMiddleware gin.HandlerFunc) {
	router.POST("", authMiddleware, h.Create)
	router.GET("", authMiddleware, h.List)
	router.GET("/:id", authMiddleware, h.GetByID)
	router.PUT("/:id", authMiddleware, h.Update)
	router.DELETE("/:id", authMiddleware, h.Delete)
}
