package proxy

import (
	"fmt"
	"net/http"

	"github.com/JojiiOfficial/ReverseProxy/models"
)

// HTTPServer http server
type HTTPServer struct {
	SSL    bool
	Routes []*models.Route
	Server *http.Server
}

// Start starts the server
func (httpServer *HTTPServer) Start() {
	go httpServer.run()
}

func (httpServer *HTTPServer) run() {
	for i := range httpServer.Routes {
		fmt.Printf("Route: %p\n", httpServer.Routes[i])
	}
}
