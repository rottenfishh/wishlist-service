package app

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"
	"wishlist-service/internal/adapter/in/httpservice"
	"wishlist-service/internal/adapter/out/repository"
	"wishlist-service/internal/database"
	"wishlist-service/internal/service/auth"
	"wishlist-service/internal/service/gift"
	"wishlist-service/internal/service/wishlist"

	_ "github.com/jackc/pgx/v5/stdlib" // register pgx driver
)

type App struct {
	server *httpservice.Server
	db     *sql.DB
}

func NewApp(ctx context.Context, cfg *Config) (*App, error) {
	InitLogging()

	db, err := ConnectDB(ctx, cfg.DatabaseConfig)
	if err != nil {
		return nil, err
	}

	userRepo := repository.NewUserRepository(db)
	wishlistRepo := repository.NewWishlistRepository(db)
	giftRepo := repository.NewGiftRepository(db)

	userService := auth.NewUserService(userRepo, cfg.JWTSecret)
	wishlistService := wishlist.NewService(wishlistRepo, giftRepo)
	giftService := gift.NewService(giftRepo, wishlistRepo)

	userHandler := httpservice.NewUserHandler(userService)
	wishlistHandler := httpservice.NewWishlistHandler(wishlistService)
	giftHandler := httpservice.NewGiftHandler(giftService)

	server := httpservice.NewServer(":8080", cfg.AuthConfig, httpservice.Handlers{
		User:     userHandler,
		Gift:     giftHandler,
		Wishlist: wishlistHandler,
	})

	return &App{server: server, db: db}, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		err := a.server.Run()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			slog.Error("starting http server", "error", err)
			errCh <- err
			return
		}
	}()

	select {
	case <-ctx.Done():
		err := a.Shutdown()
		if err != nil {
			return err
		}
	case serverError := <-errCh:
		err := a.Shutdown()
		if err != nil {
			slog.Error("shutting down app", "error", err)
		}
		return serverError
	}

	return nil
}

func (a *App) Shutdown() error {
	slog.Info("shutting down app")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var err error
	if err = a.server.Shutdown(shutdownCtx); err != nil {
		slog.Error("shutting down http server error", "err", err)
	}

	if err = a.db.Close(); err != nil {
		slog.Error("closing db error", "err", err)
	}

	return err
}

// TODO: migrations from docker compose
func ConnectDB(ctx context.Context, dbConfig DatabaseConfig) (*sql.DB, error) {
	dsn := dbConfig.DSN()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to database: %w", err)
	}

	if err = db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping database: %w", err)
	}

	err = database.MigrateDB(db, "database")
	if err != nil {
		return nil, err
	}
	slog.Info("migrations applied successfully")
	return db, nil
}
