package consumer

import (
	"github.com/alisavch/image-service/internal/log"
	"github.com/sirupsen/logrus"
)

// Logger unites interfaces.
type Logger struct {
	DisplayLog
}

// NewLogger configures Logger.
func NewLogger() *Logger {
	return &Logger{
		DisplayLog: log.NewCustomLogger(logrus.New()),
	}
}
