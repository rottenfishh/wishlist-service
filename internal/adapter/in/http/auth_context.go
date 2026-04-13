package http

import (
	"cdek/internal/model"
	"log/slog"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func extractUserId(c *gin.Context) (uuid.UUID, error) {
	userRaw, exists := c.Get("user")
	if !exists {
		return uuid.Nil, model.ErrUnauthorized
	}

	user, ok := userRaw.(*UserClaims)
	if !ok {
		slog.Error("invalid user claims")
		return uuid.Nil, model.ErrUnauthorized
	}

	parsedUserId, err := uuid.Parse(user.Subject)
	if err != nil {
		slog.Error("invalid user id", "user", user.Subject, "error", err)
		return uuid.Nil, model.ErrUnauthorized
	}

	return parsedUserId, nil
}
