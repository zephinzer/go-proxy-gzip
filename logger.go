package main

import (
	"fmt"

	"github.com/evalphobia/logrus_fluent"
	"github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
)

var fluentLevels = []logrus.Level{
	logrus.TraceLevel,
	logrus.DebugLevel,
	logrus.InfoLevel,
	logrus.WarnLevel,
	logrus.ErrorLevel,
	logrus.PanicLevel,
}

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
	logger.WithFields(logrus.Fields{
		"module": "logger",
	}).Debugf("initialising logger for %s\n", moduleName)
	fluentDHost := config.GetString("fluentd_host")
	fluentDPort := config.GetInt("fluentd_port")
	if len(fluentDHost) == 0 || fluentDPort == 0 {
		logger.WithFields(logrus.Fields{
			"module": "logger",
		}).Warnf("fluentd logs streaming is DISABLED for module '%s'\n", moduleName)
	} else if fluentHook, err := logrus_fluent.NewWithConfig(logrus_fluent.Config{
		Host: config.GetString("fluentd_host"),
		Port: config.GetInt("fluentd_port"),
	}); err != nil {
		logger.WithFields(logrus.Fields{
			"module": "logger",
		}).Error(err)
	} else {
		logger.WithFields(logrus.Fields{
			"module": "logger",
		}).Infof("fluentd logs streaming is ENABLED for module '%s', initialising for levels '%v'\n", moduleName, fluentLevels)
		fluentHook.SetLevels(fluentLevels)
		fluentHook.SetTag(fmt.Sprintf("%s.%s", config.GetString("app_id"), moduleName))
		logger.AddHook(fluentHook)
	}
	return logger.WithFields(logrus.Fields{"module": moduleName})
}
