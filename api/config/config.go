package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Config struct {
	AllowedOrigin     string
	LineChannelSecret string
	LineChannelToken  string
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
	lineChannelSecret := os.Getenv("LINE_CHANNEL_SECRET")
	lineChannelToken := os.Getenv("LINE_CHANNEL_ACCESS_TOKEN")

	return &Config{
		AllowedOrigin:     allowedOrigin,
		LineChannelSecret: lineChannelSecret,
		LineChannelToken:  lineChannelToken,
	}
}
