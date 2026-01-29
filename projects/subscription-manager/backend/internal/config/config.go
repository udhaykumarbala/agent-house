package config

import (
	"os"
	"strconv"
)

type Config struct {
	DatabaseURL   string
	SessionSecret string
	CORSOrigins   []string
	Port          string
	Environment   string
}

func Load() *Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}

	corsOrigins := []string{"http://localhost:3000"}
	if env == "production" {
		corsOrigins = []string{os.Getenv("FRONTEND_URL")}
	}

	return &Config{
		DatabaseURL:   os.Getenv("DATABASE_URL"),
		SessionSecret: os.Getenv("SESSION_SECRET"),
		CORSOrigins:   corsOrigins,
		Port:          port,
		Environment:   env,
	}
}

func GetEnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	i, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return i
}
