package utils

import (
	"fmt"
	"os"
)

// ConfigService allows you to interact with the config file.
type ConfigService interface {
	GetEnv(string) (string, error)
}

// Configure includes a config file.
type Configure struct {
	path string
}

// GetEnv gets variables from .env file.
func (conf Configure) GetEnv(key string) (string, error) {
	value, ok := os.LookupEnv(key)
	if !ok {
		return value, fmt.Errorf("error in fetching value from .env")
	}
	return value, nil
}

// NewConfig is config constructor.
func NewConfig(filepath string) ConfigService {
	return &Configure{
		path: filepath,
	}
}
