package file

import (
	"fmt"
	"sync"
	"time"

	"github.com/nilbelec/amazon-price-watcher/pkg/configuration"
	"github.com/nilbelec/amazon-price-watcher/pkg/util/file"
)

// File handles the application configuration by file
type File struct {
	sync.Mutex
	path     string
	settings *configuration.Settings
}

type settings struct {
	WebServerPort                    int     `json:"port"`
	ProductsRefreshIntervalInMinutes int     `json:"refresh_interval_minutes"`
	TelegramBotToken                 string  `json:"telegram_bot_token"`
	TelegramChatIDs                  []int64 `json:"telegram_chat_ids"`
}

// New creates the application configuration using the specified file.
// If the file doesn't exists, it will be created with default values
func New(configFile string) (config *File, err error) {
	config = &File{path: configFile}
	err = config.initialize()
	return
}

func (f *File) initialize() error {
	err := f.createIfNotExists()
	if err != nil {
		return err
	}
	return f.loadData()
}

func (f *File) createIfNotExists() error {
	exists, err := file.Exists(f.path)
	if err != nil {
		return fmt.Errorf("Error checking if config file exists: %s", err.Error())
	}
	if exists == true {
		return nil
	}
	f.settings = &configuration.Defaults
	err = f.saveData()
	if err != nil {
		return fmt.Errorf("Error creating default config file: %s", err.Error())
	}
	return nil
}

func (f *File) loadData() (err error) {
	json := &settings{}
	err = file.ReadJSON(f.path, json)
	if err != nil {
		return
	}
	f.settings = &configuration.Settings{
		ProductsRefreshIntervalInMinutes: json.ProductsRefreshIntervalInMinutes,
		TelegramBotToken:                 json.TelegramBotToken,
		TelegramChatIDs:                  json.TelegramChatIDs,
		WebServerPort:                    json.WebServerPort,
	}
	return
}

func (f *File) saveData() error {
	json := &settings{
		ProductsRefreshIntervalInMinutes: f.settings.ProductsRefreshIntervalInMinutes,
		TelegramBotToken:                 f.settings.TelegramBotToken,
		TelegramChatIDs:                  f.settings.TelegramChatIDs,
		WebServerPort:                    f.settings.WebServerPort,
	}
	return file.SaveJSON(f.path, json)
}

// Update stores the configuration
func (f *File) Update(settings *configuration.Settings) error {
	f.Lock()
	defer f.Unlock()
	f.settings = settings
	err := f.saveData()
	if err != nil {
		return fmt.Errorf("Error updating configuration file: %s", err.Error())
	}
	return nil
}

// Settings returns the current configuration settings
func (f *File) Settings() *configuration.Settings {
	return f.settings
}

// GetProductsRefreshIntervalInMinutes get the products refresh interval duration in minutes
func (f *File) GetProductsRefreshIntervalInMinutes() int {
	f.Lock()
	defer f.Unlock()
	return f.settings.ProductsRefreshIntervalInMinutes
}

// GetWebServerPort get the web server port
func (f *File) GetWebServerPort() int {
	f.Lock()
	defer f.Unlock()
	return f.settings.WebServerPort
}

// GetRefreshInterval get the products refresh interval duration
func (f *File) GetRefreshInterval() time.Duration {
	f.Lock()
	defer f.Unlock()
	return time.Duration(f.settings.ProductsRefreshIntervalInMinutes) * time.Minute
}

// WebServerAddress get the web server address
func (f *File) WebServerAddress() string {
	f.Lock()
	defer f.Unlock()
	return fmt.Sprintf(":%d", f.settings.WebServerPort)
}

// GetBotToken returns the Telegram Bot TOKEN
func (f *File) GetBotToken() string {
	f.Lock()
	defer f.Unlock()
	return f.settings.TelegramBotToken
}

// GetChatIDs returns the Telegram Chat IDs
func (f *File) GetChatIDs() []int64 {
	f.Lock()
	defer f.Unlock()
	return f.settings.TelegramChatIDs
}
