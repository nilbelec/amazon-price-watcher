package version

// Service helps to retrieve the current and latest application version
type Service struct {
	client Client
}

// Client is an interface to get the latest version of the application
type Client interface {
	LatestVersion() (string, error)
}

// NewService creates a new version Service
func NewService(c Client) *Service {
	return &Service{c}
}

// Current returns the current application version
func (s *Service) Current() string {
	return Current
}

// Latest returns the latest application version
func (s *Service) Latest() (string, error) {
	return s.client.LatestVersion()
}
