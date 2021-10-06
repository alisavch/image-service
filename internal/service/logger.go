package service

import (
	"github.com/alisavch/image-service/internal/log"
	"github.com/sirupsen/logrus"
)

// FormattingOutput contains methods for formatting log output.
type FormattingOutput interface {
	Fatalf(format string, args ...interface{})
}

// Logger unites interfaces.
type Logger struct {
	FormattingOutput
}

// NewLogger is the logger constructor.
func NewLogger() *Logger {
	return &Logger{
		FormattingOutput: log.NewCustomLogger(logrus.New()),
	}
}
