package config

import (
	"log/slog"
	"os"
	"strconv"
)

type AppConfig struct {
	Address      string
	Port         int
	UseTls       bool
	Namespace    string
	PodsSelector string
}

func loadDefaults() AppConfig {
	return AppConfig{
		Address:   "",
		Port:      8080,
		UseTls:    false,
		Namespace: "default",
	}
}

var loaded bool = false
var config AppConfig

func GetConfig() AppConfig {
	if !loaded {
		config = loadDefaults()

		// network level
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

		// k8s selectors
		configErr := false
		if namespace := os.Getenv("K8S_NAMESPACE"); len(namespace) > 0 {
			config.Namespace = namespace
			slog.Info("Resolved", "namespace", namespace)
		} else {
			slog.Error("Unable to resolve namespace", "K8S_NAMESPACE", nil)
			configErr = true
		}

		if selector := os.Getenv("K8S_PODS_SELECTOR"); len(selector) > 0 {
			config.PodsSelector = selector
			slog.Info("Resolved", "selector", selector)
		} else {
			slog.Error("Unable to resolve selector ", "K8S_PODS_SELECTOR", nil)
			configErr = true
		}

		if configErr {
			panic("Unable to start load config due to errors")
		}

		loaded = true
	}

	return config
}
