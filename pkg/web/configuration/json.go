package configuration

import "github.com/nilbelec/amazon-price-watcher/pkg/configuration"

type settingsJSON struct {
	WebServerPort                    int     `json:"port"`
	ProductsRefreshIntervalInMinutes int     `json:"refresh_interval"`
	TelegramBotToken                 string  `json:"telegram_bot_token"`
	TelegramChatIDs                  []int64 `json:"telegram_chat_ids"`
}

func toJSON(s *configuration.Settings) *settingsJSON {
	return &settingsJSON{
		ProductsRefreshIntervalInMinutes: s.ProductsRefreshIntervalInMinutes,
		TelegramBotToken:                 s.TelegramBotToken,
		TelegramChatIDs:                  s.TelegramChatIDs,
		WebServerPort:                    s.WebServerPort,
	}
}

func fromJSON(j *settingsJSON) *configuration.Settings {
	return &configuration.Settings{
		WebServerPort:                    j.WebServerPort,
		ProductsRefreshIntervalInMinutes: j.ProductsRefreshIntervalInMinutes,
		TelegramBotToken:                 j.TelegramBotToken,
		TelegramChatIDs:                  j.TelegramChatIDs,
	}
}
