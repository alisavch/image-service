package log

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		TimestampFormat:        "2006-01-02 15:04:05",
		FullTimestamp:          true,
		DisableLevelTruncation: true,
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return "", fmt.Sprintf("%s:%d", formatFilePath(f.File), f.Line)
		},
	}
	logrus.SetFormatter(formatter)
}

func formatFilePath(path string) string {
	arr := strings.Split(path, "/")
	return arr[len(arr)-1]
}

// CustomLogger contains logrus.
type CustomLogger struct {
	logger *logrus.Logger
}

// NewCustomLogger is constructor of the CustomLogger.
func NewCustomLogger(logger *logrus.Logger) *CustomLogger {
	return &CustomLogger{logger: logger}
}

// Debugf outputs debug error plus log format.
func (l *CustomLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Infof outputs info plus log format.
func (l *CustomLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Printf outputs message plus log format.
func (l *CustomLogger) Printf(format string, args ...interface{}) {
	l.logger.Printf(format, args...)
}

// Errorf outputs error plus log format.
func (l *CustomLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Fatalf outputs fatal error plus log format.
func (l *CustomLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

// Debug outputs debug to the log.
func (l *CustomLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Info outputs info to the log.
func (l *CustomLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Print outputs message to the log.
func (l *CustomLogger) Print(args ...interface{}) {
	l.logger.Print(args...)
}

// Error outputs error to the log.
func (l *CustomLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Fatal outputs fatal error to the log.
func (l *CustomLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}
