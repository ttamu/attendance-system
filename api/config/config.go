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
	if os.Getenv("GO_ENV") != "production" {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found")
		}
	} else {
		log.Println("Production mode detected. Skipping loading .env file.")
	}

	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")

	return &Config{
		AllowedOrigin: allowedOrigin,
	}
}
