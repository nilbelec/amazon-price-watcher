package products

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/amazon-price-watcher/pkg/services"
)

// Handler the home handler
type Handler struct {
	ps *services.ProductService
}

// NewHandler creates a new products handler
func NewHandler(ps *services.ProductService) *Handler {
	return &Handler{ps}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.handleGet(w)
		return
	} else if r.Method == "POST" {
		h.handlePost(w, r)
		return
	} else if r.Method == "DELETE" {
		h.handleDelete(w, r)
		return
	}
	http.NotFound(w, r)
}

func (h *Handler) handleGet(w http.ResponseWriter) {
	list, _ := h.ps.ListProducts()
	jsons := ToSliceOfJSONS(list)
	response, _ := json.Marshal(jsons)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {
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

func (h *Handler) handleDelete(w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
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
