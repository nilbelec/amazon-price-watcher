package configuration

import "github.com/nilbelec/amazon-price-watcher/pkg/configuration"

type settingsJSON struct {
	WebServerPort                    int      `json:"port"`
	ProductsRefreshIntervalInMinutes int      `json:"refresh_interval"`
	TelegramBotToken                 string   `json:"telegram_bot_token"`
	TelegramChatIDs                  []int64  `json:"telegram_chat_ids"`
	SMTPHost                         string   `json:"smtp_host"`
	SMTPPort                         int      `json:"smtp_port"`
	SMTPUsername                     string   `json:"smtp_username"`
	SMTPPassword                     string   `json:"smtp_password"`
	SMTPTo                           []string `json:"smtp_to"`
}

func toJSON(s *configuration.Settings) *settingsJSON {
	return &settingsJSON{
		ProductsRefreshIntervalInMinutes: s.ProductsRefreshIntervalInMinutes,
		TelegramBotToken:                 s.TelegramBotToken,
		TelegramChatIDs:                  s.TelegramChatIDs,
		WebServerPort:                    s.WebServerPort,
		SMTPHost:                         s.SMTPHost,
		SMTPPort:                         s.SMTPPort,
		SMTPUsername:                     s.SMTPUsername,
		SMTPPassword:                     s.SMTPPassword,
		SMTPTo:                           s.SMTPTo,
	}
}

func fromJSON(j *settingsJSON) *configuration.Settings {
	return &configuration.Settings{
		WebServerPort:                    j.WebServerPort,
		ProductsRefreshIntervalInMinutes: j.ProductsRefreshIntervalInMinutes,
		TelegramBotToken:                 j.TelegramBotToken,
		TelegramChatIDs:                  j.TelegramChatIDs,
		SMTPHost:                         j.SMTPHost,
		SMTPPort:                         j.SMTPPort,
		SMTPUsername:                     j.SMTPUsername,
		SMTPPassword:                     j.SMTPPassword,
		SMTPTo:                           j.SMTPTo,
	}
}
