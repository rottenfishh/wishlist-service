package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	DatabaseURL string `config:"database_url"`
	Name        string `config:"name"`
	Password    string `config:"password"`
	Host        string `config:"host"`
	Port        string `config:"port"`
}

func (c *DatabaseConfig) DSN() string {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.Name, c.Password, c.Host, c.Port, c.DatabaseURL)
	return dsn
}

func LoadConfig() (*DatabaseConfig, error) {
	err := LoadEnv()
	if err != nil {
		return nil, fmt.Errorf("loading environment: %w", err)
	}

	cfg := &DatabaseConfig{}
	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	cfg.Name = os.Getenv("DATABASE_NAME")
	cfg.Password = os.Getenv("DATABASE_PASSWORD")
	cfg.Host = os.Getenv("DATABASE_HOST")
	cfg.Port = os.Getenv("DATABASE_PORT")

	return cfg, nil
}

func LoadEnv() error {
	err := godotenv.Load(".env")
	if err != nil {
		return fmt.Errorf("loading .env: %w", err)
	}
	return nil
}

func InitLogging() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	slog.SetLogLoggerLevel(slog.LevelDebug)
}
