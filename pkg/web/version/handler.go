package version

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/amazon-price-watcher/pkg/version"
	"github.com/nilbelec/amazon-price-watcher/pkg/web/router"
)

type appVersion struct {
	Latest  string `json:"latest"`
	Current string `json:"current"`
}

// Handler handles the version related requests
type Handler struct {
	versionService *version.Service
}

// NewHandler creates a new VersionHandler
func NewHandler(vs *version.Service) *Handler {
	return &Handler{vs}
}

// Routes return the routes VersionHandler handles
func (h *Handler) Routes() router.Routes {
	return router.Routes{
		router.Route{Path: "/version", Method: "GET", Accepts: "application/json", HandlerFunc: h.getVersion},
	}
}

func (h *Handler) getVersion(w http.ResponseWriter, r *http.Request) {
	l, err := h.versionService.Latest()
	if err != nil {
		http.Error(w, "Error while checking latest version", http.StatusInternalServerError)
		return
	}
	version := appVersion{
		Current: h.versionService.Current(),
		Latest:  l,
	}
	data, err := json.Marshal(version)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
