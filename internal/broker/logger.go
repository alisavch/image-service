package broker

import (
	"github.com/alisavch/image-service/internal/log"

	"github.com/sirupsen/logrus"
)

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
