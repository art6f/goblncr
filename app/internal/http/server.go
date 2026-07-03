package http

import (
	"fmt"
	"net/http"

	"github.com/art6f/goblncr/config"
	"github.com/art6f/goblncr/internal/balancer"
)

func ServeHttp() {
	appConfig := config.GetConfig()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Servers is up and running!")
	})

	balancer := balancer.NewBalancer()

	go balancer.Run()

	http.ListenAndServe(fmt.Sprintf("%s:%d", appConfig.Address, appConfig.Port), nil)
}
