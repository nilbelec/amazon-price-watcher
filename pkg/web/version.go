package web

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/amazon-price-watcher/pkg/version"
)

type gitHubVersion struct {
	TagVersion string `json:"tag_name"`
}

type appVersion struct {
	Latest  string `json:"latest"`
	Current string `json:"current"`
}

type versionHandler struct {
}

func newVersionHandler() *versionHandler {
	return &versionHandler{}
}

func (h *versionHandler) handlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.NotFound(w, r)
			return
		}
		url := "https://github.com/nilbelec/amazon-price-watcher/releases/latest"

		client := &http.Client{}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Accept", "application/json")
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		lastVersion := gitHubVersion{}
		json.NewDecoder(resp.Body).Decode(&lastVersion)

		version := appVersion{
			Current: version.Version,
			Latest:  lastVersion.TagVersion,
		}
		data, err := json.Marshal(version)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(data)
	}
}
