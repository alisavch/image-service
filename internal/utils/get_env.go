package utils

import (
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

// GetEnvWithKey gets env value.
func GetEnvWithKey(key string) string {
	return os.Getenv(key)
}

// LoadEnv initially load env.
func LoadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}
