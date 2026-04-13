package app

import (
	"cdek/internal/database"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib" // register pgx driver
)

type App struct {
}

func NewApp() *App {
	return &App{}
}

func (a *App) Run(ctx context.Context) error {
	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	_, err = ConnectDB(ctx, cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	//defer db.Close()
	return nil
}

// TODO: migrations from docker compose
func ConnectDB(ctx context.Context, dbConfig *DatabaseConfig) (*sql.DB, error) {
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

	return db, nil
}
