package notification

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/amazon-price-watcher/pkg/product"
	"github.com/nilbelec/amazon-price-watcher/pkg/web/router"
)

// Handler is the notifications handler
type Handler struct {
	ps *product.Service
}

// NewHandler creates a new notifications handler
func NewHandler(ps *product.Service) *Handler {
	return &Handler{ps}
}

// Routes return the routes the notifications handler handles
func (h *Handler) Routes() router.Routes {
	return router.Routes{
		router.Route{Path: "/notifications", Method: "POST", Accepts: "*/*", HandlerFunc: h.saveNotifications},
	}
}

func (h *Handler) saveNotifications(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	d := json.NewDecoder(r.Body)
	j := &notificationsJSON{}
	err := d.Decode(j)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	n := fromJSON(j)
	err = h.ps.UpdateNotifications(url, n)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
