package http

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

type AuthConfig struct {
	JWTSecret string
}

type UserClaims struct {
	Email string
	jwt.RegisteredClaims
}

func AuthMiddleware(cfg AuthConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			slog.Error("No Authorization header")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			slog.Error("invalid authorization header: ", "header", authHeader)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		claims := &UserClaims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenUnverifiable
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			slog.Error("invalid token", slog.String("token", tokenStr), slog.String("error", err.Error()))
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
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
