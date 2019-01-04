package web

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/amazon-price-watcher/pkg/product"
)

func (s *Server) handleProducts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handleProductsGet(s.ps, w)
			return
		} else if r.Method == "POST" {
			handleProductsPost(s.ps, w, r)
			return
		} else if r.Method == "DELETE" {
			handleProductsDelete(s.ps, w, r)
			return
		}
		http.NotFound(w, r)
	}
}

func handleProductsGet(ps *product.Service, w http.ResponseWriter) {
	list, _ := ps.ListProducts()
	jsons := ToSliceOfJSONS(list)
	response, _ := json.Marshal(jsons)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func handleProductsPost(ps *product.Service, w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	_, err := ps.AddProductByURL(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func handleProductsDelete(ps *product.Service, w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")
	if url == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}
	_, err := ps.DeleteProductByURL(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
