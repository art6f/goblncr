// Package config is a main config
package config

import (
	"log/slog"
	"os"
	"strconv"

	"github.com/art6f/goblncr/internal/balancer/strategies"
	"go.yaml.in/yaml/v3"
)

type ServerConfig struct {
	Address  string                       `yaml:"address"`
	Port     int                          `yaml:"port"`
	TLS      bool                         `yaml:"tls"`
	Strategy strategies.BalancingStrategy `yaml:"strategy"`
}

type HashringConfig struct {
	Vnodes int `yaml:"vnodes"`
}

type TargetConfig struct {
	Namespace string `yaml:"namespace"`
	Selector  string `yaml:"selector"`
	Port      int    `yaml:"port"`
}

type AppConfig struct {
	Server   ServerConfig   `yaml:"server"`
	Target   TargetConfig   `yaml:"target"`
	Hashring HashringConfig `yaml:"hashring"`
}

func loadConfig() AppConfig {
	// defaults
	config := AppConfig{
		ServerConfig{
			Address:  "",
			Port:     8080,
			TLS:      false,
			Strategy: strategies.StrategyRandom,
		},
		TargetConfig{
			Namespace: "default",
			Selector:  "",
			Port:      80,
		},
		HashringConfig{
			Vnodes: 32,
		},
	}

	configData, errFile := os.ReadFile("config.yaml")
	if errFile != nil {
		slog.Error("Config file not found, using default")
	}

	if err := yaml.Unmarshal(configData, &config); err != nil {
		slog.Error("Unable to load config file, using default", "error", err)
	}

	return config
}

var config AppConfig

func GetConfig() AppConfig {
	config = loadConfig()
	updateConfigFromEnv(&config)

	slog.Info("[BALANCER] Config",
		"address", config.Server.Address,
		"port", config.Server.Port,
		"TSL", config.Server.TLS,
		"strategy", config.Server.Strategy)

	slog.Info("[TARGET] Config", "namespace", config.Target.Namespace, "selector", config.Target.Selector, "port", config.Target.Port)

	if len(config.Target.Selector) == 0 || len(config.Target.Namespace) == 0 {
		panic("Target Namespace and Selector cannot be empty")
	}

	return config
}

func updateConfigFromEnv(config *AppConfig) {
	slog.Info("Checking environment variables for addition config overrides...")

	// LB Sever Settings
	if address, ok := os.LookupEnv("BALANCER_ADDRESS"); ok {
		config.Server.Address = address
	}

	if port, err := strconv.Atoi(os.Getenv("BALANCER_PORT")); err == nil {
		config.Server.Port = port
	}

	if tls, err := strconv.ParseBool(os.Getenv("BALANCER_TLS")); err == nil {
		config.Server.TLS = tls
	}

	// Targets
	if namespace := os.Getenv("TARGET_NAMESPACE"); len(namespace) > 0 {
		config.Target.Namespace = namespace
	}

	if selector := os.Getenv("TARGET_SELECTOR"); len(selector) > 0 {
		config.Target.Selector = selector
	}

	if port, err := strconv.Atoi(os.Getenv("TARGET_PORT")); err == nil {
		config.Target.Port = port
	}
}
