package config

import (
	"errors"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL string
	Token   string
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

	token := os.Getenv("TOKEN")
	if token == "" {
		return nil, errors.New("TOKEN is required")
	}

	return &Config{
		BaseURL: baseURL,
		Token:   token,
	}, nil
}
