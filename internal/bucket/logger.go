package bucket

import (
	"github.com/alisavch/image-service/internal/log"
	"github.com/sirupsen/logrus"
)

// Logger unites interfaces.
type Logger struct {
	FormattingOutput
}

// NewLogger configures Logger.
func NewLogger() *Logger {
	return &Logger{
		FormattingOutput: log.NewCustomLogger(logrus.New()),
	}
}
