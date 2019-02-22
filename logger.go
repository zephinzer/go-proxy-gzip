package main

import (
	"github.com/sirupsen/logrus"
)

func createLogger(moduleName string, logFormat string) *logrus.Entry {
	logFormatterText := &logrus.TextFormatter{
		ForceColors:   true,
		FullTimestamp: true,
	}
	logFormatterJSON := &logrus.JSONFormatter{}
	logger := logrus.New()
	logger.SetLevel(logrus.TraceLevel)
	if logFormat == "json" {
		logger.SetFormatter(logFormatterJSON)
	} else {
		logger.SetFormatter(logFormatterText)
	}
	return logger.WithFields(logrus.Fields{"module": moduleName})
}
