// Package app is a main application
package app

import (
	"github.com/art6f/goblncr/config"
	"github.com/art6f/goblncr/internal/balancer"
	"github.com/art6f/goblncr/internal/http"
)

func Main() {
	appConfig := config.GetConfig()

	resolvedStrategy, err := appConfig.Server.Strategy.ResolveStrategy()
	if err != nil {
		panic(err)
	}

	balancer := balancer.NewBalancer(&appConfig, *resolvedStrategy)

	server := http.GetServer(&appConfig, balancer)
	server.ServeHTTP()
}
