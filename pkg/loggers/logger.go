package loggers

import (
	"github.com/sirupsen/logrus"
)

func Logger() *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.InfoLevel)
	return logger
}

var _ Advanced = (*logrus.Logger)(nil)
