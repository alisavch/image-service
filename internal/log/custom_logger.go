package log

import (
	"github.com/sirupsen/logrus"
)

// Logger unites interfaces.
type Logger struct {
	*CustomLogger
}

// NewLogger is the logger constructor.
func NewLogger() *Logger {
	return &Logger{
		CustomLogger: NewCustomLogger(logrus.New()),
	}
}
