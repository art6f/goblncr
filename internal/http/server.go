// Package http implements the HTTP reverse proxy server
package http

import (
	"embed"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/art6f/goblncr/config"
	"github.com/art6f/goblncr/internal/balancer"
)

//go:embed error_pages/*.html
var templatesFS embed.FS

type HTTPServer struct {
	serverAddress string
	serverPort    int
	targetPort    int
	balancer      *balancer.Balancer
}

var singleInstance *HTTPServer = nil

func GetServer(config *config.AppConfig, balancer *balancer.Balancer) *HTTPServer {
	if singleInstance == nil {
		singleInstance = &HTTPServer{
			serverAddress: config.Server.Address,
			serverPort:    config.Server.Port,
			targetPort:    config.Target.Port,
			balancer:      balancer,
		}
	}

	return singleInstance
}

func (server *HTTPServer) ServeHTTP() {
	template.Must(
		template.ParseFS(
			templatesFS,
			"error_pages/*.html",
		),
	)

	go server.balancer.Run()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		selectedPod, err := server.balancer.SelectServer(r)

		if err != nil {
			server.errorHandler(w, r, err)
			return
		}

		proxyURL, _ := url.Parse(fmt.Sprintf("http://%s:%d", selectedPod, server.targetPort))
		proxy := httputil.NewSingleHostReverseProxy(proxyURL)

		// handlers/hooks
		proxy.ErrorHandler = server.errorHandler
		proxy.ModifyResponse = func(resp *http.Response) error {
			slog.Info(fmt.Sprintf("%d - %s%s", resp.StatusCode, proxyURL.String(), r.RequestURI))
			return nil
		}

		proxy.ServeHTTP(w, r)
	})

	http.ListenAndServe(fmt.Sprintf("%s:%d", server.serverAddress, server.serverPort), nil)
}

func (server *HTTPServer) errorHandler(w http.ResponseWriter, r *http.Request, e error) {
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
