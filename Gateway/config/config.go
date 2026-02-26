package config

import (
	"log"
	"os"
)

type Config struct {
	UserServiceURL  string
	EventServiceURL string
	MatchServiceURL string
	JWTSecret       string
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func Load() *Config {
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == ""{
		log.Fatal("JWT_SECRET is not set")
	}

	return &Config{
		UserServiceURL:  getEnv("USER_SERVICE_URL", "http://user-service:8081"),
		MatchServiceURL: getEnv("MATCH_SERVICE_URL", "http://match-service:8082"),
		EventServiceURL: getEnv("EVENT_SERVICE_URL", "http://event-service:8083"),
		JWTSecret:       jwtSecret,
	}
}
