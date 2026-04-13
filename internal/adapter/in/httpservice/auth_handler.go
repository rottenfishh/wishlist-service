package httpservice

import (
	"log/slog"
	"net/http"
	"wishlist-service/internal/adapter/in/dto"
	"wishlist-service/internal/model"
	"wishlist-service/internal/service/auth"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service auth.Service
}

func NewUserHandler(service auth.Service) *UserHandler {
	return &UserHandler{service: service}
}

// Register creates a new user with email and password
// @Summary Register a new user
// @Description Creates a new user account with email, password and optional role. Password is hashed before storing.
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration details"
// @Success 201 {object} model.User "User created successfully"
// @Failure 400 {object} dto.ErrorResponse "Invalid request body"
// @Failure 409 {object} dto.ErrorResponse "User with this email already exists"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.BindJSON(&req); err != nil {
		slog.Error("failed to bind register request", "error", err)
		writeError(c, model.ErrInvalidRequest)
		return
	}

	savedUser, err := h.service.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		slog.Error("failed to register user", "error", err)
		writeError(c, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, savedUser)
}

// Login authenticates user and returns JWT token
// @Summary Login user
// @Description Authenticates user with email and password, returns JWT token for subsequent requests
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.Token "JWT token"
// @Failure 400 {object} dto.ErrorResponse "Invalid request body"
// @Failure 401 {object} dto.ErrorResponse "Invalid email or password"
// @Failure 500 {object} dto.ErrorResponse "Internal server error"
// @Router /api/auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		slog.Error("failed to bind login request", "error", err)
		writeError(c, model.ErrInvalidRequest)
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		slog.Error("failed to log in user", "error", err)
		writeError(c, err)
		return
	}

	tokenResp := dto.Token{Token: token}
	c.IndentedJSON(http.StatusOK, tokenResp)
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)
}
