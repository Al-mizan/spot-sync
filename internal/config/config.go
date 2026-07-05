package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	Dsn         string
	JwtSecret   string
	Environment string
}

func LoadEnv() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		// Log but don't fatal, maybe env is set in OS
		log.Println("No .env file found, relying on environment variables")
	}

	return &Config{
		Port:        os.Getenv("PORT"),
		Dsn:         os.Getenv("DSN"),
		JwtSecret:   os.Getenv("JWT_SECRET"),
		Environment: os.Getenv("ENVIRONMENT"),
	}, nil
}
