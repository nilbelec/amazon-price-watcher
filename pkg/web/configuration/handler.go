package configuration

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/amazon-price-watcher/pkg/configuration"
)

// Data contains the json representation of the app configuration
type Data struct {
	WebServerPort                    int     `json:"port"`
	ProductsRefreshIntervalInMinutes int     `json:"refresh_interval"`
	TelegramBotToken                 string  `json:"telegram_bot_token"`
	TelegramChatIDs                  []int64 `json:"telegram_chat_ids"`
}

// Handler handles the configurations requests
type Handler struct {
	config configuration.Configuration
}

// NewHandler creates a new products handler
func NewHandler(config configuration.Configuration) *Handler {
	return &Handler{config}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		h.handleGet(w)
		return
	} else if r.Method == "POST" {
		h.handlePost(w, r)
		return
	}
	http.NotFound(w, r)
}

func (h *Handler) handleGet(w http.ResponseWriter) {
	settings := h.config.Settings()
	configData := &Data{
		ProductsRefreshIntervalInMinutes: settings.ProductsRefreshIntervalInMinutes,
		TelegramBotToken:                 settings.TelegramBotToken,
		TelegramChatIDs:                  settings.TelegramChatIDs,
		WebServerPort:                    settings.WebServerPort,
	}
	response, _ := json.Marshal(configData)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *Handler) handlePost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var data Data
	err := decoder.Decode(&data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	settings := &configuration.Settings{
		WebServerPort:                    data.WebServerPort,
		ProductsRefreshIntervalInMinutes: data.ProductsRefreshIntervalInMinutes,
		TelegramBotToken:                 data.TelegramBotToken,
		TelegramChatIDs:                  data.TelegramChatIDs,
	}
	err = h.config.Update(settings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
