package config

import (
	"context"
	"os"

	"github.com/Burtcam/encounter-builder-backend/logger"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Config struct {
	GH_TOKEN      string
	REPO_URL      string
	DB_CONNECTION string
	DBPool        *pgxpool.Pool
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		logger.Log.Error("Warning no .env file found, assuming system environment vars are already set")
	}

	cfg := Config{
		GH_TOKEN:      os.Getenv("GH_TOKEN"),
		REPO_URL:      os.Getenv("REPO_URL"),
		DB_CONNECTION: os.Getenv("DB_CONNECTION_STRING"),
	}
	logger.Log.Info("Configuration succesfully Loaded")
	ctx := context.Background()
	logger.Log.Debug("Creating DB Connection Pool")
	// Create connection pool
	pool, err := pgxpool.New(ctx, cfg.DB_CONNECTION)
	if err != nil {
		logger.Log.Error("failed to create DB pool: %v", "err", err)
		return nil
	}
	cfg.DBPool = pool

	logger.Log.Info("Configuration succesfully Loaded")

	return &cfg
}
