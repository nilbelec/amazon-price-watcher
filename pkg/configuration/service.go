package configuration

import (
	"fmt"
	"time"
)

// Repository handles the configuration persistence
type Repository interface {
	Save(settings *Settings) error
	Exists() (bool, error)
	Get() (*Settings, error)
}

// Service helps to update or retrieve the current configuration settings
type Service struct {
	repo     Repository
	settings *Settings
}

// New creates a new configuration service
func New(repo Repository) (s *Service, err error) {
	s = &Service{repo: repo}
	err = s.init()
	return
}

func (s *Service) init() (err error) {
	err = s.createIfNotExists()
	if err != nil {
		return
	}
	return s.loadSettings()
}

func (s *Service) loadSettings() (err error) {
	s.settings, err = s.repo.Get()
	return
}

func (s *Service) createIfNotExists() (err error) {
	exists, err := s.repo.Exists()
	if exists || err != nil {
		return
	}
	return s.repo.Save(Defaults)
}

// Save saves/overwrites the settings values
func (s *Service) Save(settings *Settings) error {
	return s.repo.Save(settings)
}

// Settings returns the current settings values
func (s *Service) Settings() *Settings {
	return s.settings
}

// GetProductsRefreshIntervalInMinutes get the products refresh interval duration in minutes
func (s *Service) GetProductsRefreshIntervalInMinutes() int {
	return s.settings.ProductsRefreshIntervalInMinutes
}

// GetWebServerPort get the web server port
func (s *Service) GetWebServerPort() int {
	return s.settings.WebServerPort
}

// RefreshInterval get the products refresh interval duration
func (s *Service) RefreshInterval() time.Duration {
	return time.Duration(s.settings.ProductsRefreshIntervalInMinutes) * time.Minute
}

// Address get the web server address
func (s *Service) Address() string {
	return fmt.Sprintf(":%d", s.settings.WebServerPort)
}

// GetBotToken returns the Telegram Bot TOKEN
func (s *Service) GetBotToken() string {
	return s.settings.TelegramBotToken
}

// GetChatIDs returns the Telegram Chat IDs
func (s *Service) GetChatIDs() []int64 {
	return s.settings.TelegramChatIDs
}
