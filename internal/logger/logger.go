package logger

type Logger interface {
	Info(args ...any)
	Infof(format string, args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Fatal(v ...any)
	Fatalf(format string, v ...any)
	Debug(v ...any)
	Debugf(format string, v ...any)
	WithFields(fields Fields) Logger
}

type Fields map[string]any
