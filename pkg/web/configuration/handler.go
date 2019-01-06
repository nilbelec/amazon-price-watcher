package configuration

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/amazon-price-watcher/pkg/configuration"
	"github.com/nilbelec/amazon-price-watcher/pkg/web/router"
)

// Handler is the configuration handler
type Handler struct {
	cs *configuration.Service
}

// NewHandler creates a new configuration handler
func NewHandler(cs *configuration.Service) *Handler {
	return &Handler{cs}
}

// Routes return the routes the configuration handler handles
func (h *Handler) Routes() router.Routes {
	return router.Routes{
		router.Route{Path: "/configuration", Method: "GET", Accepts: "application/json", HandlerFunc: h.getConfiguration},
		router.Route{Path: "/configuration", Method: "POST", Accepts: "*/*", HandlerFunc: h.saveConfiguration},
	}
}

func (h *Handler) getConfiguration(w http.ResponseWriter, r *http.Request) {
	s := h.cs.Settings()
	j := toJSON(s)
	b, _ := json.Marshal(j)
	w.Header().Set("Content-Type", "application/json")
	w.Write(b)
}

func (h *Handler) saveConfiguration(w http.ResponseWriter, r *http.Request) {
	d := json.NewDecoder(r.Body)
	j := &settingsJSON{}
	err := d.Decode(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	s := fromJSON(j)
	err = h.cs.Save(s)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
