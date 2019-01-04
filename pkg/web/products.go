package web

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/amazon-price-watcher/pkg/product"
)

type productsHandler struct {
	ps *product.Service
}

func newProductsHandler(ps *product.Service) *productsHandler {
	return &productsHandler{ps}
}

func (h *productsHandler) handlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h.handleProductsGet(w)
			return
		} else if r.Method == "POST" {
			h.handleProductsPost(w, r)
			return
		} else if r.Method == "DELETE" {
			h.handleProductsDelete(w, r)
			return
		}
		http.NotFound(w, r)
	}
}

func (h *productsHandler) handleProductsGet(w http.ResponseWriter) {
	list, _ := h.ps.ListProducts()
	jsons := ToSliceOfJSONS(list)
	response, _ := json.Marshal(jsons)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *productsHandler) handleProductsPost(w http.ResponseWriter, r *http.Request) {
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

func (h *productsHandler) handleProductsDelete(w http.ResponseWriter, r *http.Request) {
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
