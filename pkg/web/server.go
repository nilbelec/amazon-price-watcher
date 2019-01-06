package web

import (
	"log"
	"net/http"

	"github.com/nilbelec/amazon-price-watcher/pkg/configuration"
	"github.com/nilbelec/amazon-price-watcher/pkg/product"
	"github.com/nilbelec/amazon-price-watcher/pkg/version"

	webconfiguration "github.com/nilbelec/amazon-price-watcher/pkg/web/configuration"
	webhome "github.com/nilbelec/amazon-price-watcher/pkg/web/home"
	webnotifications "github.com/nilbelec/amazon-price-watcher/pkg/web/notification"
	webproduct "github.com/nilbelec/amazon-price-watcher/pkg/web/product"
	"github.com/nilbelec/amazon-price-watcher/pkg/web/router"
	webversion "github.com/nilbelec/amazon-price-watcher/pkg/web/version"
)

// Server web server
type Server struct {
	productsService      *product.Service
	configurationService *configuration.Service
	versionService       *version.Service
}

// NewServer creates a new web server
func NewServer(ps *product.Service, cs *configuration.Service, vs *version.Service) *Server {
	return &Server{ps, cs, vs}
}

// Start starts the web server
func (s *Server) Start() {
	log.Println("Web server started at " + s.configurationService.Address())
	log.Fatal(http.ListenAndServe(s.configurationService.Address(), s.router()))
}

func (s *Server) router() *router.Router {
	r := router.New()

	r.AddHandler(webhome.NewHandler())
	r.AddHandler(webconfiguration.NewHandler(s.configurationService))
	r.AddHandler(webproduct.NewHandler(s.productsService))
	r.AddHandler(webnotifications.NewHandler(s.productsService))
	r.AddHandler(webversion.NewHandler(s.versionService))

	return r
}
