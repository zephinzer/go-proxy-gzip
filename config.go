package main

import (
	config "github.com/spf13/viper"
)

const (
	// DefaultAddress lets you change the interface we listen on (0.0.0.0 works for *most* systems)
	DefaultAddress = "0.0.0.0"
	// DefaultPort lets you change the port incase of port conflicts
	DefaultPort = "1337"
	// DefaultForwardTo lets you add an address to forward to, empty for an echoserver so you can test it more easily
	DefaultForwardTo = ""
	// DefaultContentType lets you assign a custom content type, useful for when http.detectContentType fails
	DefaultContentType = ""
	// DefaultLogFormat lets you change the format of logs depending on whether its development or production
	DefaultLogFormat = "text"
)

func initConfiguration() {
	config.SetDefault("addr", DefaultAddress)
	config.BindEnv("addr")
	config.SetDefault("port", DefaultPort)
	config.BindEnv("port")
	config.SetDefault("forward_to", DefaultForwardTo)
	config.BindEnv("forward_to")
	config.SetDefault("content_type", DefaultContentType)
	config.BindEnv("content_type")
	config.SetDefault("log_format", DefaultLogFormat)
	config.BindEnv("log_format")
}
