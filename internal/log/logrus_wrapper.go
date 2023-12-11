package log

import "github.com/sirupsen/logrus"

func New() Logger {
	return &LogrusLogger{
		logger: logrus.New(),
	}
}

type LogrusLogger struct {
	logger interface {
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
		WithFields(fields logrus.Fields) *logrus.Entry
	}
}

func (l *LogrusLogger) Info(args ...any) {
	l.logger.Info(args...)
}
func (l *LogrusLogger) Infof(format string, args ...any) {
	l.logger.Infof(format, args...)
}
func (l *LogrusLogger) Warn(args ...any) {
	l.logger.Warn(args...)
}
func (l *LogrusLogger) Warnf(format string, args ...any) {
	l.logger.Warnf(format, args...)
}
func (l *LogrusLogger) Error(args ...any) {
	l.logger.Error(args...)
}
func (l *LogrusLogger) Errorf(format string, args ...any) {
	l.logger.Errorf(format, args...)
}
func (l *LogrusLogger) Fatal(args ...any) {
	l.logger.Fatal(args...)
}
func (l *LogrusLogger) Fatalf(format string, args ...any) {
	l.logger.Fatalf(format, args...)
}
func (l *LogrusLogger) Debug(args ...any) {
	l.logger.Debug(args...)
}
func (l *LogrusLogger) Debugf(format string, args ...any) {
	l.logger.Debugf(format, args...)
}
func (l *LogrusLogger) WithFields(fields Fields) Logger {
	return &LogrusLogger{
		logger: l.logger.WithFields(logrus.Fields(fields)),
	}
}
