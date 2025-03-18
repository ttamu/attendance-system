package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	AllowedOrigin string
}

func LoadEnv() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")

	return &Config{
		AllowedOrigin: allowedOrigin,
	}
}
