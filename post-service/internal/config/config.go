package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBURL     string
	JWTSecret string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		Port:      getEnv("PORT", "8081"),
		DBURL:     mustGetEnv("DB_URL"),
		JWTSecret: mustGetEnv("JWT_SECRET"),
	}

	return cfg
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

func mustGetEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("%s is required", key)
	}
	return v
}