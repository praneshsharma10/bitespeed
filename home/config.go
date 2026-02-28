package home

import "os"

type Config struct {
	Port        string
	DatabaseURL string
}

var AppConfig Config

func LoadConfig() {
	AppConfig = Config{
		Port:        getEnv("PORT", "8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
