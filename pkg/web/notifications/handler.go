package notifications

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/amazon-price-watcher/pkg/model"

	"github.com/nilbelec/amazon-price-watcher/pkg/services"
	"github.com/nilbelec/amazon-price-watcher/pkg/web/products"
)

// Handler handles the notifications requests
type Handler struct {
	ps *services.ProductService
}

// NewHandler creates a new products handler
func NewHandler(ps *services.ProductService) *Handler {
	return &Handler{ps}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.NotFound(w, r)
		return
	}
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	decoder := json.NewDecoder(r.Body)
	var data products.Notifications
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	notifications := model.Notifications{
		BackInStock:    data.BackInStock,
		OutOfStock:     data.OutOfStock,
		PriceBelows:    data.PriceBelows,
		PriceDecreases: data.PriceDecreases,
		PriceIncreases: data.PriceIncreases,
		PriceOver:      data.PriceOver,
	}
	err = h.ps.UpdateProductNotifications(url, notifications)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)

}
