package httpservice

import (
	"log/slog"
	"strings"
	"wishlist-service/internal/model"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	Email string
	jwt.RegisteredClaims
}

func AuthMiddleware(cfg AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			slog.Error("No Authorization header")
			writeError(c, model.ErrUnauthorized)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			slog.Error("invalid authorization header", "header", authHeader)
			writeError(c, model.ErrUnauthorized)
			c.Abort()
			return
		}

		tokenStr := parts[1]

		claims := &UserClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenUnverifiable
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			slog.Error("invalid token", "token", tokenStr, "error", err)
			writeError(c, model.ErrUnauthorized)
			c.Abort()
			return
		}

		c.Set("user", &UserClaims{
			Email: claims.Email,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: claims.Subject,
			},
		})
		slog.Debug("User authenticated", "subject:", claims.Subject, "email", claims.Email)
		c.Next()
	}
}
