package main

import (
	"github.com/sirupsen/logrus"
	config "github.com/spf13/viper"
	"github.com/usvc/go-log/pkg/constants"
	"github.com/usvc/go-log/pkg/hooks/fluentd"
	"github.com/usvc/go-log/pkg/logger"
)

func createLogger(moduleName string, logFormat string) *logrus.Entry {
	log := logger.New(logFormat)
	fluentHook := fluentd.NewHook(&fluentd.HookConfig{
		Host:                    config.GetString("fluentd_host"),
		InitializeRetryCount:    config.GetInt("fluentd_init_retry_count"),
		InitializeRetryInterval: config.GetDuration("fluentd_init_retry_interval"),
		Levels:                  constants.DefaultHookLevels,
		Port:                    config.GetInt("fluentd_port"),
		Tag:                     config.GetString("app_id"),
	})
	log.AddHook(fluentHook)
	return log.WithFields(logrus.Fields{"module": moduleName})
}
