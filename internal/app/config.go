package app

import (
	"fmt"
	"log/slog"
	"os"
	"wishlist-service/internal/adapter/in/httpservice"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseConfig
	httpservice.AuthConfig
}

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

func LoadConfig() (*Config, error) {
	err := LoadEnv()
	if err != nil {
		return nil, fmt.Errorf("loading environment: %w", err)
	}

	cfg := DatabaseConfig{}
	cfg.DatabaseURL = os.Getenv("DATABASE_URL")
	cfg.Name = os.Getenv("DATABASE_NAME")
	cfg.Password = os.Getenv("DATABASE_PASSWORD")
	cfg.Host = os.Getenv("DATABASE_HOST")
	cfg.Port = os.Getenv("DATABASE_PORT")

	authCfg := httpservice.AuthConfig{
		JWTSecret: os.Getenv("JWT_SECRET"),
	}

	return &Config{DatabaseConfig: cfg, AuthConfig: authCfg}, nil
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
