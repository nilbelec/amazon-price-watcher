package storage

import (
	"sync"

	"github.com/nilbelec/amazon-price-watcher/pkg/configuration"
)

// ConfigurationFile handles the application configuration persistence using a JSON file
type ConfigurationFile struct {
	sync.Mutex
	filename string
	settings *configuration.Settings
}

type settingsJSON struct {
	WebServerPort                    int     `json:"port"`
	ProductsRefreshIntervalInMinutes int     `json:"refresh_interval_minutes"`
	TelegramBotToken                 string  `json:"telegram_bot_token"`
	TelegramChatIDs                  []int64 `json:"telegram_chat_ids"`
}

// NewConfigurationFile creates a new ConfigurationFile to handle the configuration persistence
func NewConfigurationFile(filename string) (cf *ConfigurationFile, err error) {
	cf = &ConfigurationFile{filename: filename}
	exists, err := cf.Exists()
	if !exists || err != nil {
		return
	}
	err = cf.load()
	return
}

func (cf *ConfigurationFile) persist() (err error) {
	cf.Lock()
	defer cf.Unlock()
	json := &settingsJSON{
		ProductsRefreshIntervalInMinutes: cf.settings.ProductsRefreshIntervalInMinutes,
		TelegramBotToken:                 cf.settings.TelegramBotToken,
		TelegramChatIDs:                  cf.settings.TelegramChatIDs,
		WebServerPort:                    cf.settings.WebServerPort,
	}
	return SaveJSON(cf.filename, &json)
}

func (cf *ConfigurationFile) load() (err error) {
	cf.Lock()
	defer cf.Unlock()
	json := &settingsJSON{}
	err = ReadJSON(cf.filename, json)
	if err != nil {
		return
	}
	cf.settings = &configuration.Settings{
		ProductsRefreshIntervalInMinutes: json.ProductsRefreshIntervalInMinutes,
		TelegramBotToken:                 json.TelegramBotToken,
		TelegramChatIDs:                  json.TelegramChatIDs,
		WebServerPort:                    json.WebServerPort,
	}
	return
}

// Get returns the current configuration values
func (cf *ConfigurationFile) Get() *configuration.Settings {
	return cf.settings
}

// Exists checks if the configuration file exists
func (cf *ConfigurationFile) Exists() (bool, error) {
	return Exists(cf.filename)
}

// Save stores the settings values in the file
func (cf *ConfigurationFile) Save(s *configuration.Settings) error {
	cf.settings = s
	return cf.persist()
}
