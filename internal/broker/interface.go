package broker

// FormattingOutput contains methods for formatting log output.
type FormattingOutput interface {
	Printf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}
