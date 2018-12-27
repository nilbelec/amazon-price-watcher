package home

import (
	"log"
	"net/http"

	"github.com/gobuffalo/packr"
)

// View view model for home
type View struct {
}

// Handler the home handler
type Handler struct {
	box packr.Box
}

// NewHandler creates a new home handler
func NewHandler() *Handler {
	box := packr.NewBox(".")
	return &Handler{box}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.NotFound(w, r)
		return
	}
	b, err := h.box.Find("home.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	w.Write(b)
}
