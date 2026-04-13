package main

import (
	"cdek/internal/app"
	"context"
	"log/slog"
	"os"
)

func main() {
	ctx := context.Background()

	wishlistApp := app.NewApp()

	err := wishlistApp.Run(ctx)
	if err != nil {
		slog.Error("Error running app", "error", err)
		ctx.Done()
		os.Exit(1)
	}
}
