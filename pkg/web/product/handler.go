package product

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/amazon-price-watcher/pkg/product"
	"github.com/nilbelec/amazon-price-watcher/pkg/web/router"
)

// Handler handles the products requests
type Handler struct {
	ps *product.Service
}

// NewHandler creates a new products handler
func NewHandler(ps *product.Service) *Handler {
	return &Handler{ps}
}

// Routes return the routes the products handler handles
func (h *Handler) Routes() router.Routes {
	return router.Routes{
		router.Route{Path: "/products", Method: "GET", Accepts: "application/json", HandlerFunc: h.listProducts},
		router.Route{Path: "/products", Method: "POST", Accepts: "*/*", HandlerFunc: h.saveProduct},
		router.Route{Path: "/products", Method: "DELETE", Accepts: "*/*", HandlerFunc: h.deleteProduct},
	}
}

func (h *Handler) listProducts(w http.ResponseWriter, r *http.Request) {
	ps, _ := h.ps.ListProducts()
	jsons := toJSONs(ps)
	response, _ := json.Marshal(jsons)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) saveProduct(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	_, err := h.ps.AddProductByURL(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	if url == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	_, err := h.ps.DeleteProductByURL(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
