package file

import (
	"sync"

	"github.com/nilbelec/amazon-price-watcher/pkg/configuration"
	"github.com/nilbelec/amazon-price-watcher/pkg/util/file"
)

// Repository handles the application configuration percistence using a file
type Repository struct {
	sync.Mutex
	path string
}

type settings struct {
	WebServerPort                    int     `json:"port"`
	ProductsRefreshIntervalInMinutes int     `json:"refresh_interval_minutes"`
	TelegramBotToken                 string  `json:"telegram_bot_token"`
	TelegramChatIDs                  []int64 `json:"telegram_chat_ids"`
}

// New creates the application configuration using the specified file.
func New(configRepository string) *Repository {
	return &Repository{path: configRepository}
}

// Get returns the configuration values from the file
func (r *Repository) Get() (s *configuration.Settings, err error) {
	r.Lock()
	defer r.Unlock()
	json := &settings{}
	err = file.ReadJSON(r.path, json)
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
func (r *Repository) Exists() (bool, error) {
	return file.Exists(r.path)
}

// Save stores the settings values in the file
func (r *Repository) Save(s *configuration.Settings) error {
	r.Lock()
	defer r.Unlock()
	json := &settings{
		ProductsRefreshIntervalInMinutes: s.ProductsRefreshIntervalInMinutes,
		TelegramBotToken:                 s.TelegramBotToken,
		TelegramChatIDs:                  s.TelegramChatIDs,
		WebServerPort:                    s.WebServerPort,
	}
	return file.SaveJSON(r.path, json)
}
