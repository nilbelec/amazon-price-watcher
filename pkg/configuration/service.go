package configuration

import (
	"fmt"
	"time"
)

// Repository handles the configuration persistence
type Repository interface {
	Save(settings *Settings) error
	Exists() (bool, error)
	Get() *Settings
}

// Service helps to update or retrieve the current configuration settings
type Service struct {
	repo Repository
}

// NewService creates a new configuration service
func NewService(repo Repository) (s *Service, err error) {
	s = &Service{repo: repo}
	err = s.createIfNotExists()
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
	return s.repo.Get()
}

// GetProductsRefreshIntervalInMinutes get the products refresh interval duration in minutes
func (s *Service) GetProductsRefreshIntervalInMinutes() int {
	ss := s.repo.Get()
	return ss.ProductsRefreshIntervalInMinutes
}

// GetWebServerPort get the web server port
func (s *Service) GetWebServerPort() int {
	ss := s.repo.Get()
	return ss.WebServerPort
}

// RefreshInterval get the products refresh interval duration
func (s *Service) RefreshInterval() time.Duration {
	ss := s.repo.Get()
	return time.Duration(ss.ProductsRefreshIntervalInMinutes) * time.Minute
}

// Address get the web server address
func (s *Service) Address() string {
	ss := s.repo.Get()
	return fmt.Sprintf(":%d", ss.WebServerPort)
}

// GetBotToken returns the Telegram Bot TOKEN
func (s *Service) GetBotToken() string {
	ss := s.repo.Get()
	return ss.TelegramBotToken
}

// GetChatIDs returns the Telegram Chat IDs
func (s *Service) GetChatIDs() []int64 {
	ss := s.repo.Get()
	return ss.TelegramChatIDs
}
