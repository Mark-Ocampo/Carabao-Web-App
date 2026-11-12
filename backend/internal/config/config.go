package config

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatabaseURL  string
	CookieSecure bool
}

func Load() (Config, error) {
	// .env is optional: in production the variables come from the environment
	_ = godotenv.Load()

	cfg := Config{
		Port:         os.Getenv("PORT"),
		DatabaseURL:  os.Getenv("DATABASE_URL"),
		CookieSecure: os.Getenv("COOKIE_SECURE") == "true",
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is not set")
	}

	return cfg, nil
}
