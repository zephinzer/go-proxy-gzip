package main

import (
	config "github.com/spf13/viper"
)

const (
	// DefaultAddress lets you change the interface we listen on (0.0.0.0 works for *most* systems)
	DefaultAddress = "0.0.0.0"
	// DefaultAppID is the default ID to use for the application
	DefaultAppID = "go_proxy_gzip_default"
	// DefaultContentType lets you assign a custom content type, useful for when http.detectContentType fails
	DefaultContentType = ""
	// DefaultFluentDHost defines the default place where a FluentD service is expected
	DefaultFluentDHost = ""
	// DefaultFluentDPort defines the default port on which the FluentD service is listening on
	DefaultFluentDPort = 0
	// DefaultForwardTo lets you add an address to forward to, empty for an echoserver so you can test it more easily
	DefaultForwardTo = ""
	// DefaultLogFormat lets you change the format of logs depending on whether its development or production
	DefaultLogFormat = "text"
	// DefaultPort lets you change the port incase of port conflicts
	DefaultPort = "1337"
)

func initConfiguration() {
	config.SetDefault("addr", DefaultAddress)
	config.SetDefault("app_id", DefaultAppID)
	config.SetDefault("content_type", DefaultContentType)
	config.SetDefault("fluentd_host", DefaultFluentDHost)
	config.SetDefault("fluentd_port", DefaultFluentDPort)
	config.SetDefault("forward_to", DefaultForwardTo)
	config.SetDefault("log_format", DefaultLogFormat)
	config.SetDefault("port", DefaultPort)
	config.AutomaticEnv()
}
