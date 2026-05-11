package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	AuthServiceURL string
	PostServiceURL string
}

func Load() *Config {

	_ = godotenv.Load()

	return &Config{
		Port:           getEnv("PORT", "8000"),
		AuthServiceURL: mustGetEnv("AUTH_SERVICE_URL"),
		PostServiceURL: mustGetEnv("POST_SERVICE_URL"),
	}
}

func getEnv(key, fallback string) string {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func mustGetEnv(key string) string {

	value := os.Getenv(key)

	if value == "" {
		log.Fatalf("missing env variable: %s", key)
	}

	return value
}