package config

import (
	"errors"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL      string
	ServerPort   string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

func NewConfig() (*Config, error) {
	_, err := os.Stat(".env")
	if err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("Warning: .env file exists but failed to load: %v", err)
		}
	}

	investURL := os.Getenv("T_INVEST_URL")
	if investURL == "" {
		log.Println("T_INVEST_URL variable is not found in environment")
		return nil, errors.New("T_INVEST_URL is required")
	}

	serverPort := os.Getenv("PORT")
	if serverPort == "" {
		log.Println("PORT variable is not found in environment\nSetting default 8080")
		serverPort = "8080"
	}
	readTimeoutStr := os.Getenv("READ_TIMEOUT")
	if readTimeoutStr == "" {
		log.Println("READ_TIMEOUT variable is not found in environment\nSetting default 10s")
		readTimeoutStr = "10"
	}
	writeTimeoutStr := os.Getenv("WRITE_TIMEOUT")
	if writeTimeoutStr == "" {
		log.Println("WRITE_TIMEOUT variable is not found in environment\nSetting default 10s")
		writeTimeoutStr = "10"
	}
	idleTimeoutStr := os.Getenv("IDLE_TIMEOUT")
	if idleTimeoutStr == "" {
		log.Println("IDLE_TIMEOUT variable is not found in environment\nSetting default 30s")
		idleTimeoutStr = "30"
	}

	readTimeout, err := strconv.Atoi(readTimeoutStr)
	if err != nil {
		log.Println("Wrong format of READ_TIMEOUT\nSetting default 10s")
		readTimeout = 10
	}
	writeTimeout, err := strconv.Atoi(writeTimeoutStr)
	if err != nil {
		log.Println("Wrong format of WRITE_TIMEOUT\nSetting default 10s")
		writeTimeout = 10
	}
	idleTimeout, err := strconv.Atoi(idleTimeoutStr)
	if err != nil {
		log.Println("Wrong format of IDLE_TIMEOUT\nSetting default 30s")
		idleTimeout = 30
	}

	return &Config{
		BaseURL:      investURL,
		ServerPort:   serverPort,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
	}, nil
}
