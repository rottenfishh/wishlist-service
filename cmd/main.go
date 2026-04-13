package main

import (
	"cdek/internal/app"
	"context"
	"log/slog"
	"os/signal"
	"syscall"
)

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
