package config

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"sync"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// GlobalConfig holds the application's configuration values loaded from environment variables.
type GlobalConfig struct {
	PublicHost        string
	Port              string
	DBUser            string
	DBPassword        string
	DBHost            string
	DBPort            string
	DBName            string
	From              string
	To                string
	Password          string
	SMTPHost          string
	SMTPPort          string
	MatrixUser        string
	MatrixPassword    string
	MatrixAccessToken string
}

var (
	config     GlobalConfig
	once       sync.Once
)

// LoadConfig initializes configuration values from environment variables.
func LoadConfig() GlobalConfig {
	once.Do(func() {
		config = GlobalConfig{
			PublicHost:        getEnv("PUBLIC_HOST", "http://localhost"),
			Port:              getEnv("PORT", "3000"),
			DBUser:            getEnv("DB_USER", "localdev"),
			DBPassword:        getEnv("DB_PASSWORD", "localdev"),
			DBHost:            getEnv("DB_HOST", "192.168.8.35"),
			DBPort:            getEnv("DB_PORT", "5432"),
			DBName:            getEnv("DB_NAME", "localdev"),
			From:              getEnv("FROM", ""),
			To:                getEnv("TO", ""),
			Password:          getEnv("PASSWORD", ""),
			SMTPHost:          getEnv("SMTP_HOST", ""),
			SMTPPort:          getEnv("SMTP_PORT", ""),
			MatrixUser:        getEnv("MATRIX_USER", ""),
			MatrixPassword:    getEnv("MATRIX_PASSWORD", ""),
			MatrixAccessToken: getEnv("MATRIX_ACCESS_TOKEN", ""),
		}
	})
	return config
}

// getEnv retrieves an environment variable or returns the provided fallback value.
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// ConnectDB establishes and returns a database connection using database/sql with pgx compatibility.
func ConnectDB(ctx context.Context) (*sql.DB, error) {
	cfg := LoadConfig()
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

