package web

import (
	"log"
	"net/http"

	cf "github.com/nilbelec/amazon-price-watcher/pkg/configuration"
	"github.com/nilbelec/amazon-price-watcher/pkg/services"

	"github.com/nilbelec/amazon-price-watcher/pkg/web/configuration"
	"github.com/nilbelec/amazon-price-watcher/pkg/web/home"
	"github.com/nilbelec/amazon-price-watcher/pkg/web/notifications"
	"github.com/nilbelec/amazon-price-watcher/pkg/web/products"
	"github.com/nilbelec/amazon-price-watcher/pkg/web/version"
)

// Server web server
type Server struct {
	ps      *services.ProductService
	config  cf.Configuration
	version string
}

// NewServer creates a new web server
func NewServer(ps *services.ProductService, config cf.Configuration, version string) *Server {
	return &Server{ps, config, version}
}

// Start starts the web server
func (w *Server) Start() {
	http.Handle("/", home.NewHandler())
	http.Handle("/products", products.NewHandler(w.ps))
	http.Handle("/configuration", configuration.NewHandler(w.config))
	http.Handle("/notifications", notifications.NewHandler(w.ps))
	http.Handle("/version", version.NewHandler(w.version))
	log.Println("Web server started at " + w.config.WebServerAddress())
	log.Fatal(http.ListenAndServe(w.config.WebServerAddress(), nil))
}
