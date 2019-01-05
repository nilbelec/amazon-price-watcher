package storage

import (
	"sync"

	"github.com/nilbelec/amazon-price-watcher/pkg/configuration"
)

// ConfigurationFile handles the application configuration persistence using a json file
type ConfigurationFile struct {
	sync.Mutex
	path string
}

type settingsJSON struct {
	WebServerPort                    int     `json:"port"`
	ProductsRefreshIntervalInMinutes int     `json:"refresh_interval_minutes"`
	TelegramBotToken                 string  `json:"telegram_bot_token"`
	TelegramChatIDs                  []int64 `json:"telegram_chat_ids"`
}

// NewConfigurationFile creates the application configuration using the specified file.
func NewConfigurationFile(path string) *ConfigurationFile {
	return &ConfigurationFile{path: path}
}

// Get returns the configuration values from the file
func (cf *ConfigurationFile) Get() (s *configuration.Settings, err error) {
	cf.Lock()
	defer cf.Unlock()
	json := &settingsJSON{}
	err = ReadJSON(cf.path, json)
	if err != nil {
		return
	}
	s = &configuration.Settings{
		ProductsRefreshIntervalInMinutes: json.ProductsRefreshIntervalInMinutes,
		TelegramBotToken:                 json.TelegramBotToken,
		TelegramChatIDs:                  json.TelegramChatIDs,
		WebServerPort:                    json.WebServerPort,
	}
	return
}

// Exists checks if the configuration file exists
func (cf *ConfigurationFile) Exists() (bool, error) {
	return Exists(cf.path)
}

// Save stores the settings values in the file
func (cf *ConfigurationFile) Save(s *configuration.Settings) error {
	cf.Lock()
	defer cf.Unlock()
	json := &settingsJSON{
		ProductsRefreshIntervalInMinutes: s.ProductsRefreshIntervalInMinutes,
		TelegramBotToken:                 s.TelegramBotToken,
		TelegramChatIDs:                  s.TelegramChatIDs,
		WebServerPort:                    s.WebServerPort,
	}
	return SaveJSON(cf.path, json)
}
