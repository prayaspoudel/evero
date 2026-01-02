package logger

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type logrusLogger struct {
	logger *logrus.Logger
}

func NewLogrusLogger() (Logger, error) {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return &logrusLogger{logger: log}, nil
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *logrusLogger) Fatalln(args ...interface{}) {
	l.logger.Fatalln(args...)
}

func (l *logrusLogger) WithFields(fields Fields) Logger {
	return &logrusLogEntry{
		entry: l.logger.WithFields(convertToLogrusFields(fields)),
	}
}

func (l *logrusLogger) WithError(err error) Logger {
	return &logrusLogEntry{
		entry: l.logger.WithError(err),
	}
}

type logrusLogEntry struct {
	entry *logrus.Entry
}

func (l *logrusLogEntry) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logrusLogEntry) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *logrusLogEntry) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *logrusLogEntry) Fatalln(args ...interface{}) {
	l.entry.Fatalln(args...)
}

func (l *logrusLogEntry) WithFields(fields Fields) Logger {
	return &logrusLogEntry{
		entry: l.entry.WithFields(convertToLogrusFields(fields)),
	}
}

func (l *logrusLogEntry) WithError(err error) Logger {
	return &logrusLogEntry{
		entry: l.entry.WithError(err),
	}
}

func convertToLogrusFields(fields Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	for index, field := range fields {
		logrusFields[index] = field
	}

	return logrusFields
}

// NewLogger creates a new Logrus logger instance based on configuration
// This creates a configured logrus.Logger for backwards compatibility
func NewLogger(viper *viper.Viper) *logrus.Logger {
	log := logrus.New()
	log.SetLevel(logrus.Level(viper.GetInt32("log.level")))
	log.SetFormatter(&logrus.JSONFormatter{})
	return log
}

// NewStructuredLogger creates a new structured logger using the infrastructure logger
func NewStructuredLogger() (Logger, error) {
	return NewLoggerFactory(InstanceLogrusLogger)
}
