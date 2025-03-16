package loggers

import (
	"github.com/sirupsen/logrus"
)

func Logger() *logrus.Logger {
	logger := logrus.New()
	logger.SetReportCaller(true)
	return logger
}

var _ Advanced = (*logrus.Logger)(nil)
