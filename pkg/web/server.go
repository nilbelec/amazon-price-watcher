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

type handler interface {
	handlerFunc() http.HandlerFunc
}

// NewServer creates a new web server
func NewServer(config ServerConfiguration, ps *product.Service, cs configuration.Service) *Server {
	return &Server{config, ps, cs}
}

// Start starts the web server
func (s *Server) Start() {
	s.addHandler("/", newHomeHandler())
	s.addHandler("/products", newProductsHandler(s.ps))
	s.addHandler("/configuration", newConfigurationHandler(s.cs))
	s.addHandler("/notifications", newNotificationsHandler(s.ps))
	s.addHandler("/version", newVersionHandler())

	log.Println("Web server started at " + s.webConfig.Address())
	log.Fatal(http.ListenAndServe(s.webConfig.Address(), nil))
}

func (s *Server) addHandler(pattern string, h handler) {
	http.HandleFunc(pattern, h.handlerFunc())
}
