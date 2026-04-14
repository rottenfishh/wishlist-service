package app

import (
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"wishlist-service/internal/adapter/in/httpservice"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseConfig
	httpservice.AuthConfig
	ServerPort string
}

type DatabaseConfig struct {
	MigrationsDir string
	DatabaseName  string
	Username      string
	Password      string
	Host          string
	Port          string
}

func (c *DatabaseConfig) DSN() string {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		c.Username, c.Password, c.Host, c.Port, c.DatabaseName)
	return dsn
}

func LoadConfig() (*Config, error) {
	cfg := DatabaseConfig{}
	cfg.DatabaseName = os.Getenv("DATABASE_NAME")
	cfg.Username = os.Getenv("DATABASE_USERNAME")
	cfg.Password = os.Getenv("DATABASE_PASSWORD")
	cfg.Host = os.Getenv("DATABASE_HOST")
	cfg.Port = os.Getenv("DATABASE_PORT")

	migrationsDir := os.Getenv("MIGRATIONS_DIR")
	if migrationsDir == "" {
		migrationsDir = "database"
	}
	cfg.MigrationsDir = migrationsDir

	JwtExpires, err := strconv.Atoi(os.Getenv("JWT_EXPIRES_IN_SEC"))
	if err != nil {
		slog.Error("incorrect JWT_EXPIRATION value", "error", err)
		JwtExpires = 10000
	}

	authCfg := httpservice.AuthConfig{
		JWTSecret:  os.Getenv("JWT_SECRET"),
		JwtExpires: int64(JwtExpires),
	}

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8080"
	}
	return &Config{DatabaseConfig: cfg, AuthConfig: authCfg,
		ServerPort: serverPort}, nil
}

// LoadEnv if run locally
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
