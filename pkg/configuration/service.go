package configuration

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

// WebServerPort get the web server port
func (s *Service) WebServerPort() int {
	return s.repo.Get().WebServerPort
}

// ProductsRefreshIntervalInMinutes get the products refresh interval duration in minutes
func (s *Service) ProductsRefreshIntervalInMinutes() int {
	return s.repo.Get().ProductsRefreshIntervalInMinutes
}

// TelegramBotToken returns the Telegram Bot TOKEN
func (s *Service) TelegramBotToken() string {
	return s.repo.Get().TelegramBotToken
}

// TelegramChatIDs returns the Telegram Chat IDs
func (s *Service) TelegramChatIDs() []int64 {
	return s.repo.Get().TelegramChatIDs
}

// SMTPHost returns the SMTP host value
func (s *Service) SMTPHost() string {
	return s.repo.Get().SMTPHost
}

// SMTPPort returns the SMTP port value
func (s *Service) SMTPPort() int {
	return s.repo.Get().SMTPPort
}

// SMTPUsername returns the SMTP username value
func (s *Service) SMTPUsername() string {
	return s.repo.Get().SMTPUsername
}

// SMTPPassword returns the SMTP password value
func (s *Service) SMTPPassword() string {
	return s.repo.Get().SMTPPassword
}

// SMTPTo returns the SMTP receivers
func (s *Service) SMTPTo() []string {
	return s.repo.Get().SMTPTo
}
