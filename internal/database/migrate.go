package database

import (
	"database/sql"
	"log/slog"

	"github.com/pressly/goose/v3"
)

// MigrateDB applies SQL migrations located in migrationsDir using goose.
func MigrateDB(db *sql.DB, migrationsDir string) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		slog.Error("error setting postgres dialect for goose migrations", "error", err)
		return err
	}
	return goose.Up(db, migrationsDir)
}
