package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseUrl string
	ServerUrl   string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("No .env file found")
	}
	return &Config{
		DatabaseUrl: getEnvValue("DATABASE_URL", ""),
		ServerUrl:   getEnvValue("SERVER_URL", ""),
	}
}

func getEnvValue(key string, defaultVal string) string {
	if val, found := os.LookupEnv(key); found {
		return val
	}

	return defaultVal
}
