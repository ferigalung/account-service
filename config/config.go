package config

import (
	"os"

	"github.com/gofiber/fiber/v2/log"
	"github.com/joho/godotenv"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type BasicAuthConfig struct {
	Username string
	Password string
}

type AppConfig struct {
	BasicAuth *BasicAuthConfig
	DB        *DBConfig
}

// LoadConfig load env variable from .env file
func LoadConfig() *AppConfig {
	err := godotenv.Load()
	if err != nil {
		log.Errorf("Failed to load environment variable: error:%s", err.Error())
	}

	return &AppConfig{
		DB: &DBConfig{
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     GetEnv("DB_PORT", "5432"),
			User:     GetEnv("DB_USER", "admin"),
			Password: GetEnv("DB_PASS", "admin"),
			DBName:   GetEnv("DB_NAME", "mydb"),
		},
		BasicAuth: &BasicAuthConfig{
			Username: GetEnv("BASIC_USER", "basic"),
			Password: GetEnv("BASIC_PASS", "basic"),
		},
	}
}

// GetEnv get spesific env variable, set to default if null
func GetEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
