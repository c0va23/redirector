package log

import (
	"os"

	logrus "github.com/sirupsen/logrus"
)

type loggerHook string

func (logger loggerHook) Fire(entry *logrus.Entry) error {
	entry.WithField("logger", logger)
	return nil
}

func (logger loggerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

const timeFormat = "2006-01-02T15:04:05Z"

// NewLogger create new configured logger
func NewLogger(
	logger string,
	level logrus.Level,
) *logrus.Logger {
	l := &logrus.Logger{
		Out: os.Stderr,
		Formatter: &logrus.TextFormatter{
			DisableColors:   true,
			TimestampFormat: timeFormat,
		},
		Hooks: make(logrus.LevelHooks),
		Level: level,
	}

	l.AddHook(loggerHook(logger))

	return l
}

var logger = NewLogger("logger", logrus.InfoLevel)

const defaultLoggerLevel = logrus.InfoLevel
const levelEnvvar = "LOG_LEVEL"

func loggerLevel() logrus.Level {
	loggerLevel, configured := os.LookupEnv(levelEnvvar)
	if !configured {
		return defaultLoggerLevel
	}

	level, err := logrus.ParseLevel(loggerLevel)
	if nil != err {
		logger.WithError(err).Fatalf("Invalid level error: %s", err)
	}
	return level
}

// NewLeveledLogger build logger with level configured with ennvar LOG_LEVEL
func NewLeveledLogger(logger string) *logrus.Logger {
	return NewLogger(logger, loggerLevel())
}
