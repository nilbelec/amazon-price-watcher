package home

import (
	"html/template"
	"log"
	"net/http"
)

// View view model for home
type View struct {
}

// Handler the home handler
type Handler struct {
}

// NewHandler creates a new home handler
func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		return
	}
	p, err := template.ParseFiles("../../pkg/web/home/home.html")
	if err != nil {
		log.Fatal(err)
		return
	}
	p.Execute(w, View{})
}
