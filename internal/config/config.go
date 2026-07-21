package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL     string
	Token       string
	TargetsPath string
	DBPath      string
}

func NewConfig() (*Config, error) {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("env file found but failed to load: %v", err)
		}
	}

	baseURL := os.Getenv("T_INVEST_URL")
	if baseURL == "" {
		return nil, errors.New("T_INVEST_URL is required")
	}

	token := os.Getenv("TOKEN_TORG")
	if token == "" {
		return nil, errors.New("TOKEN is required")
	}

	targetsPath := os.Getenv("TARGETS_PATH")
	if targetsPath == "" {
		return nil, errors.New("TARGETS_PATH is required")
	}

	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		return nil, errors.New("DB_PATH is required")
	}

	return &Config{
		BaseURL:     baseURL,
		Token:       token,
		TargetsPath: targetsPath,
		DBPath:      dbPath,
	}, nil
}
