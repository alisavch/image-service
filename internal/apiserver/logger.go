package apiserver

import (
	"github.com/alisavch/image-service/internal/log"
	"github.com/sirupsen/logrus"
)

// DisplayLog contains methods for log display.
type DisplayLog interface {
	Info(args ...interface{})
	Printf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

// Logger unites interfaces.
type Logger struct {
	DisplayLog
}

// NewLogger is the logger constructor.
func NewLogger() *Logger {
	return &Logger{
		DisplayLog: log.NewCustomLogger(logrus.New()),
	}
}
