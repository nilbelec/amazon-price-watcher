package web

import (
	"encoding/json"
	"net/http"

	"github.com/nilbelec/amazon-price-watcher/pkg/configuration"
)

type jsonSettings struct {
	WebServerPort                    int     `json:"port"`
	ProductsRefreshIntervalInMinutes int     `json:"refresh_interval"`
	TelegramBotToken                 string  `json:"telegram_bot_token"`
	TelegramChatIDs                  []int64 `json:"telegram_chat_ids"`
}

type configurationHandler struct {
	cs configuration.Service
}

func newConfigurationHandler(cs configuration.Service) *configurationHandler {
	return &configurationHandler{cs}
}

func (h *configurationHandler) handlerFunc() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			h.handleConfigurationGet(w)
			return
		} else if r.Method == "POST" {
			h.handleConfigurationPost(w, r)
			return
		}
		http.NotFound(w, r)
	}
}

func (h *configurationHandler) handleConfigurationGet(w http.ResponseWriter) {
	s := h.cs.Settings()
	jsonSettings := &jsonSettings{
		ProductsRefreshIntervalInMinutes: s.ProductsRefreshIntervalInMinutes,
		TelegramBotToken:                 s.TelegramBotToken,
		TelegramChatIDs:                  s.TelegramChatIDs,
		WebServerPort:                    s.WebServerPort,
	}
	response, _ := json.Marshal(jsonSettings)

	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (h *configurationHandler) handleConfigurationPost(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var jsonSettings jsonSettings
	err := decoder.Decode(&jsonSettings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	settings := &configuration.Settings{
		WebServerPort:                    jsonSettings.WebServerPort,
		ProductsRefreshIntervalInMinutes: jsonSettings.ProductsRefreshIntervalInMinutes,
		TelegramBotToken:                 jsonSettings.TelegramBotToken,
		TelegramChatIDs:                  jsonSettings.TelegramChatIDs,
	}
	err = h.cs.Update(settings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
