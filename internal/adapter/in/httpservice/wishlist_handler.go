package httpservice

import (
	"log/slog"
	"net/http"
	"strconv"
	"wishlist-service/internal/adapter/in/dto"
	"wishlist-service/internal/model"
	"wishlist-service/internal/service/wishlist"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type WishlistHandler struct {
	service wishlist.Service
}

func NewWishlistHandler(service wishlist.Service) *WishlistHandler {
	return &WishlistHandler{service: service}
}

// Create  godoc
// @Summary Create wishlist
// @Description Creates a new wishlist for the authenticated user
// @Tags wishlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateWishlistRequest true "Wishlist data"
// @Success 201 {object} dto.WishlistResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/wishlists [post]
func (h *WishlistHandler) Create(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		writeError(c, err)
		return
	}

	var req dto.CreateWishlistRequest
	if err = c.ShouldBindJSON(&req); err != nil {
		slog.Error("Binding wishlist request", "error", err)
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

// Update godoc
// @Summary Update wishlist
// @Description Updates a wishlist owned by the authenticated user
// @Tags wishlists
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Wishlist ID"
// @Param request body dto.UpdateWishlistRequest true "Updated wishlist data"
// @Success 200 {object} dto.WishlistResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/wishlists/details/{id} [put]
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

// Delete godoc
// @Summary Delete wishlist
// @Description Deletes a wishlist owned by the authenticated user
// @Tags wishlists
// @Produce json
// @Security BearerAuth
// @Param id path int true "Wishlist ID"
// @Success 200 {object} dto.WishlistResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/wishlists/details/{id} [delete]
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

// List godoc
// @Summary List wishlists
// @Description Returns all wishlists of the authenticated user
// @Tags wishlists
// @Produce json
// @Security BearerAuth
// @Success 200 {array} dto.WishlistResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/wishlists [get]
func (h *WishlistHandler) List(c *gin.Context) {
	userID, err := extractUserID(c)
	if err != nil {
		writeError(c, err)
		return
	}

	list, err := h.service.GetByUserID(c.Request.Context(), userID)
	if err != nil {
		slog.Error("getting wishlist by user", "error", err)
		writeError(c, err)
		return
	}

	resp := dto.ToListOfWishlistResponse(list)
	c.JSON(http.StatusOK, resp)
}

// GetByID godoc
// @Summary Get wishlist by ID
// @Description Get wishlist with items for the authenticated owner
// @Tags wishlists
// @Produce json
// @Security BearerAuth
// @Param id path int true "Wishlist ID"
// @Success 200 {object} dto.WishlistDetailsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 403 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/wishlists/details/{id} [get]
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

// GetByToken godoc
// @Summary Get wishlist with items by public token
// @Description Returns public wishlist details by token
// @Tags wishlists
// @Produce json
// @Param token path string true "Wishlist public token (UUID)"
// @Success 200 {object} dto.WishlistDetailsResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /api/public/wishlists/token/{token} [get]
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
	router.GET("/details/:id", authMiddleware, h.GetByID)
	router.PUT("/details/:id", authMiddleware, h.Update)
	router.DELETE("/details/:id", authMiddleware, h.Delete)
}
