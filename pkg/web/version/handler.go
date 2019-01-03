package version

import (
	"encoding/json"
	"net/http"
)

// Handler handles the versions requests
type Handler struct {
	currentVersion string
}

type gitHubVersion struct {
	TagVersion string `json:"tag_name"`
}

type appVersion struct {
	Latest  string `json:"latest"`
	Current string `json:"current"`
}

// NewHandler creates a new products handler
func NewHandler(currentVersion string) *Handler {
	return &Handler{currentVersion}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
		Current: h.currentVersion,
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
