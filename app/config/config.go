package config

import (
	"log/slog"
	"os"
	"strconv"
)

type AppConfig struct {
	Address string
	Port    int
	UseTls  bool
}

func loadDefaults() AppConfig {
	return AppConfig{
		Address: "",
		Port:    8080,
		UseTls:  false,
	}
}

var loaded bool = false
var config AppConfig

func GetConfig() AppConfig {
	if !loaded {
		config = loadDefaults()

		if address, ok := os.LookupEnv("HTTP_ADDRESS"); ok {
			config.Address = address
			slog.Info("Resolved", "address", address)
		} else {
			slog.Warn("Unable to resolve address, failing back to empty string")
		}

		if port, err := strconv.Atoi(os.Getenv("HTTP_PORT")); err == nil {
			config.Port = port
			slog.Info("Resolved", "port", config.Port)
		} else {
			slog.Warn("Unable to resolve port, failing back to", "port", config.Port)
		}

		if tls, err := strconv.ParseBool(os.Getenv("HTTP_TLS")); err == nil {
			config.UseTls = tls
			slog.Info("Resolved", "TSL", config.UseTls)
		} else {
			slog.Warn("Unable to resolve TLS flag, failing back to", "TSL", config.UseTls)
		}

		loaded = true
	}

	return config
}
