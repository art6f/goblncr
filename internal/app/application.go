package app

import (
	"github.com/art6f/goblncr/config"
	"github.com/art6f/goblncr/internal/balancer"
	"github.com/art6f/goblncr/internal/http"
)

func Main() {
	appConfig := config.GetConfig()
	balancer := balancer.NewBalancer(&appConfig)

	http.ServeHttp(&appConfig, balancer)
}
