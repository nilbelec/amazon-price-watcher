package web

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/amazon-price-watcher/pkg/product"
)

type notificationsHandler struct {
	ps *product.Service
}

func newNotificationsHandler(ps *product.Service) *notificationsHandler {
	return &notificationsHandler{ps}
}

func (h *notificationsHandler) handlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
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
		var data Notifications
		err := decoder.Decode(&data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		notifications := product.Notifications{
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
}
