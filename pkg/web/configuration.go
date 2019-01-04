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

func (s *Server) handleConfiguration() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" {
			handleConfigurationGet(s.cs, w)
			return
		} else if r.Method == "POST" {
			handleConfigurationPost(s.cs, w, r)
			return
		}
		http.NotFound(w, r)
	}
}

func handleConfigurationGet(cs configuration.Service, w http.ResponseWriter) {
	s := cs.Settings()
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

func handleConfigurationPost(cs configuration.Service, w http.ResponseWriter, r *http.Request) {
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
	err = cs.Update(settings)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}
