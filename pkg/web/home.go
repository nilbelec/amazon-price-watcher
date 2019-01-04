package web

import (
	"log"
	"net/http"

	"github.com/gobuffalo/packr"
)

func (s *Server) handleHome() http.HandlerFunc {
	box := packr.NewBox(".")
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.NotFound(w, r)
			return
		}
		b, err := box.Find("home.html")
		if err != nil {
			log.Fatal(err)
			return
		}
		w.Write(b)
	}
}
