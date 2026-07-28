package config

import (
	"log/slog"
	"os"
	"strconv"
)

type ServerConfig struct {
	Address string
	Port    int
	Tls     bool
}

type TargetConfig struct {
	Namespace string
	Selector  string
	Port      int
}

type AppConfig struct {
	Server ServerConfig
	Target TargetConfig
}

func loadDefaults() AppConfig {
	return AppConfig{
		ServerConfig{
			Address: "",
			Port:    8080,
			Tls:     false,
		},
		TargetConfig{
			Namespace: "default",
			Selector:  "",
			Port:      80,
		},
	}
}

var loaded bool = false
var config AppConfig

func GetConfig() AppConfig {
	if !loaded {
		config = loadDefaults()

		// LB Sever Settings
		if address, ok := os.LookupEnv("BALANCER_ADDRESS"); ok {
			config.Server.Address = address
			slog.Info("[BALANCER] Resolved", "address", address)
		} else {
			slog.Warn("[BALANCER] Unable to resolve address, failing back to empty string")
		}

		if port, err := strconv.Atoi(os.Getenv("BALANCER_PORT")); err == nil {
			config.Server.Port = port
			slog.Info("[BALANCER] Resolved", "port", config.Server.Port)
		} else {
			slog.Warn("[BALANCER] Unable to resolve port, failing back to", "port", config.Server.Port)
		}

		if tls, err := strconv.ParseBool(os.Getenv("BALANCER_TLS")); err == nil {
			config.Server.Tls = tls
			slog.Info("[BALANCER] Resolved", "TSL", config.Server.Tls)
		} else {
			slog.Warn("[BALANCER] Unable to resolve TLS flag, failing back to", "TSL", config.Server.Tls)
		}

		// Target Settings
		configErr := false
		if namespace := os.Getenv("TARGET_NAMESPACE"); len(namespace) > 0 {
			config.Target.Namespace = namespace
			slog.Info("[TARGET] Resolved", "namespace", namespace)
		} else {
			slog.Error("[TARGET] Unable to resolve namespace", "TARGET_NAMESPACE", nil)
			configErr = true
		}

		if selector := os.Getenv("TARGET_SELECTOR"); len(selector) > 0 {
			config.Target.Selector = selector
			slog.Info("[TARGET] Resolved", "selector", selector)
		} else {
			slog.Error("[TARGET] Unable to resolve selector ", "TARGET_SELECTOR", nil)
			configErr = true
		}

		if port, err := strconv.Atoi(os.Getenv("TARGET_PORT")); err == nil {
			config.Target.Port = port
			slog.Info("[TARGET] Resolved", "port", config.Target.Port)
		} else {
			slog.Warn("[TARGET] Unable to resolve port, failing back to", "port", config.Target.Port)
		}

		if configErr {
			panic("Unable to start load config due to errors")
		}

		loaded = true
	}

	return config
}
