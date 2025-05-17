package config

import (
	"github.com/Burtcam/encounter-builder-backend/logger"
	"os"
	"github.com/joho/godotenv"
)
type Config struct {
	GH_TOKEN string
	REPO_URL string
}

func Load() *Config {
	err:= godotenv.Load()
	if err != nil {
		logger.Log.Error("Warning no .env file found, assuming system environment vars are already set")
	}
	cfg := Config{
		GH_TOKEN : os.Getenv("GH_TOKEN"),
		REPO_URL : os.Getenv("REPO_URL"),
	}
	return &cfg
}