package http

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"log/slog"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/art6f/goblncr/config"
	"github.com/art6f/goblncr/internal/balancer"
)

//go:embed error_pages/*.html
var templatesFS embed.FS

func ServeHttp(config *config.AppConfig, balancer *balancer.Balancer) {
	template.Must(
		template.ParseFS(
			templatesFS,
			"error_pages/*.html",
		),
	)

	go balancer.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		pods := balancer.GetActivePodsIp()

		if len(pods) == 0 {
			errorHandler(w, r, errors.New("No services available"))
			return
		}

		selectedPod := pods[rand.Intn(len(pods))]

		proxyURL, _ := url.Parse(fmt.Sprintf("http://%s:%d", selectedPod, config.Target.Port))
		proxy := httputil.NewSingleHostReverseProxy(proxyURL)

		// handlers/hooks
		proxy.ErrorHandler = errorHandler
		proxy.ModifyResponse = func(resp *http.Response) error {
			slog.Info(fmt.Sprintf("%d - %s%s", resp.StatusCode, proxyURL.String(), r.RequestURI))
			return nil
		}

		proxy.ServeHTTP(w, r)
	})

	http.ListenAndServe(fmt.Sprintf("%s:%d", config.Server.Address, config.Server.Port), nil)
}

func errorHandler(w http.ResponseWriter, r *http.Request, e error) {
	slog.Error("Unable to handle request: " + e.Error())

	statusCode := 500
	title, description := "Error", "Unknow server error"

	tmpl, err := template.ParseFS(templatesFS, "error_pages/error_template.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("%d Error", statusCode), statusCode)
		return
	}

	if r.Response != nil {
		statusCode = r.Response.StatusCode
		title = r.Response.Status
		description = http.StatusText(statusCode)
	}

	data := map[string]interface{}{
		"StatusCode":  statusCode,
		"Title":       title,
		"Description": description,
	}

	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, fmt.Sprintf("%d Error", statusCode), statusCode)
	}
}
