package utils

import "os"

// DBConfig includes database variables.
type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Name     string
}

// RabbitmqConfig includes rabbitmq variables.
type RabbitmqConfig struct {
	RabbitmqURL string
}

// Config includes config variables.
type Config struct {
	DBConfig DBConfig
	TokenTTL string
	Rabbitmq RabbitmqConfig
}

// NewConfig returns a new Config struct
func NewConfig() *Config {
	return &Config{
		DBConfig: DBConfig{
			User:     getEnv("DB_USER", ""),
			Password: getEnv("DB_PASSWORD", ""),
			Host:     getEnv("DB_HOST", ""),
			Port:     getEnv("DB_PORT", ""),
			Name:     getEnv("DB_NAME", ""),
		},
		TokenTTL: getEnv("TOKEN_TTL", "12h"),
		Rabbitmq: RabbitmqConfig{
			RabbitmqURL: getEnv("RABBITMQ_URL", ""),
		},
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key string, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}
