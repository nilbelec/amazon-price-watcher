package web

import (
	"log"
	"net/http"

	"github.com/nilbelec/amazon-price-watcher/pkg/configuration"
	"github.com/nilbelec/amazon-price-watcher/pkg/product"
)

// Server web server
type Server struct {
	webConfig ServerConfiguration
	ps        *product.Service
	cs        configuration.Service
}

// ServerConfiguration handles the web server configuration
type ServerConfiguration interface {
	Address() string
}

// NewServer creates a new web server
func NewServer(config ServerConfiguration, ps *product.Service, cs configuration.Service) *Server {
	return &Server{config, ps, cs}
}

// Start starts the web server
func (s *Server) Start() {
	s.prepareHandlers()
	log.Println("Web server started at " + s.webConfig.Address())
	log.Fatal(http.ListenAndServe(s.webConfig.Address(), nil))
}

func (s *Server) prepareHandlers() {
	http.HandleFunc("/", s.handleHome())
	http.HandleFunc("/products", s.handleProducts())
	http.HandleFunc("/configuration", s.handleConfiguration())
	http.HandleFunc("/notifications", s.handleNotifications())
	http.HandleFunc("/version", s.handleVersion())
}
