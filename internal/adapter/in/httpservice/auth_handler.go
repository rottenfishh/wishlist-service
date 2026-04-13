package httpservice

import (
	"cdek/internal/adapter/in/dto"
	"cdek/internal/service/auth"
	"log/slog"
	"net/http"

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
// @Router /auth/register [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest

	if err := c.BindJSON(&req); err != nil {
		slog.Error("failed to bind register request", "error", err)
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	savedUser, err := h.service.Register(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		slog.Error("failed to register user", "error", err)
		c.JSON(http.StatusInternalServerError, "internal server error")
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
// @Router /auth/login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.BindJSON(&req); err != nil {
		slog.Error("failed to bind loginrequest", "error", err)
		c.JSON(http.StatusBadRequest, "invalid request")
		return
	}

	token, err := h.service.Login(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		slog.Error("failed to log in user", "error", err)
		c.JSON(http.StatusForbidden, "forbidden")
		return
	}

	tokenResp := dto.Token{Token: token}
	c.IndentedJSON(http.StatusOK, tokenResp)
}

func (h *UserHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/register", h.Register)
	router.POST("/login", h.Login)
}
