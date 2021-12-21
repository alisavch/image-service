package utils

import "os"

// DBConfig includes database variables.
type DBConfig struct {
	DatabaseURL string
	User        string
	Password    string
	Host        string
	Port        string
	DBName      string
}

// RabbitmqConfig includes rabbitmq variables.
type RabbitmqConfig struct {
	RabbitmqURL string
}

// BucketConfig includes bucket variables.
type BucketConfig struct {
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	BucketName         string
}

// Authentication includes variables for generating token.
type Authentication struct {
	TokenTTL   string
	SigningKey string
}

// Config includes config variables.
type Config struct {
	DBConfig DBConfig
	Auth     Authentication
	Rabbitmq RabbitmqConfig
	Bucket   BucketConfig
	Storage  string
}

// NewConfig returns a new Config struct
func NewConfig() *Config {
	return &Config{
		DBConfig: DBConfig{
			DatabaseURL: getEnv("DATABASE_URL", ""),
		},
		Auth: Authentication{
			TokenTTL:   getEnv("TOKEN_TTL", "12h"),
			SigningKey: getEnv("SIGNING_KEY", ""),
		},
		Rabbitmq: RabbitmqConfig{
			RabbitmqURL: getEnv("RABBITMQ_URL", ""),
		},
		Bucket: BucketConfig{
			AWSRegion:          getEnv("AWS_REGION", ""),
			AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
			AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
			BucketName:         getEnv("BUCKET_NAME", ""),
		},
		Storage: getEnv("REMOTE_STORAGE", "local"),
	}
}

// Simple helper function to read an environment or return a default value
func getEnv(key, defaultVal string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultVal
}
