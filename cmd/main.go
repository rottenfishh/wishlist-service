package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"
	"wishlist-service/internal/app"
)

// @title Wishlist service
// @version 1.0
// @description Wishlist service for creating wishlists, adding items, sharing public link for wishlists and booking gifts
// @host localhost:8080
// @BasePath /
// @schemes http https

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.
func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	config, err := app.LoadConfig()
	if err != nil {
		slog.Error("Error loading config", "err", err)
		return
	}

	wishlistApp, err := app.NewApp(ctx, config)
	if err != nil {
		slog.Error("Error creating wishlist app", "err", err)
		return
	}

	err = wishlistApp.Run(ctx)
	if err != nil {
		slog.Error("Error running app", "error", err)
		return
	}
}
