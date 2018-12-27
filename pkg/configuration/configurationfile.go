package configuration

import (
	"fmt"
	"sync"
	"time"

	"github.com/nilbelec/amazon-price-watcher/pkg/file"
)

// File handles the application configuration by file
type File struct {
	path string
	data *ConfigData
}

var mutex = &sync.Mutex{}

// ConfigData struct to handle the configuration options
type ConfigData struct {
	WebServerPort                    int     `json:"port"`
	ProductsRefreshIntervalInMinutes int     `json:"refresh_interval_minutes"`
	TelegramBotToken                 string  `json:"telegram_bot_token"`
	TelegramChatIDs                  []int64 `json:"telegram_chat_ids"`
}

var defaultConfig = ConfigData{
	WebServerPort:                    10035,
	ProductsRefreshIntervalInMinutes: 5,
	TelegramBotToken:                 "",
	TelegramChatIDs:                  make([]int64, 0),
}

// Load loads the application configuration from the specified file.
// If the file doesn't exists, it will be created with default values
func Load(configFile string) (config *File, err error) {
	config = &File{configFile, &ConfigData{}}
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
	f.data = &defaultConfig
	err = f.saveData()
	if err != nil {
		return fmt.Errorf("Error creating default config file: %s", err.Error())
	}
	return nil
}

func (f *File) loadData() error {
	return file.ReadJSON(f.path, f.data)
}

func (f *File) saveData() error {
	return file.SaveJSON(f.path, f.data)
}

// UpdateConfigurationData stores the configuration
func (f *File) UpdateConfigurationData(data *ConfigData) error {
	mutex.Lock()
	defer mutex.Unlock()
	f.data = data
	err := f.saveData()
	if err != nil {
		return fmt.Errorf("Error updating configuration file: %s", err.Error())
	}
	return nil
}

// GetProductsRefreshIntervalInMinutes get the products refresh interval duration in minutes
func (f *File) GetProductsRefreshIntervalInMinutes() int {
	mutex.Lock()
	defer mutex.Unlock()
	return f.data.ProductsRefreshIntervalInMinutes
}

// GetWebServerPort get the web server port
func (f *File) GetWebServerPort() int {
	mutex.Lock()
	defer mutex.Unlock()
	return f.data.WebServerPort
}

// GetRefreshInterval get the products refresh interval duration
func (f *File) GetRefreshInterval() time.Duration {
	mutex.Lock()
	defer mutex.Unlock()
	return time.Duration(f.data.ProductsRefreshIntervalInMinutes) * time.Minute
}

// GetAddress get the web server address
func (f *File) GetAddress() string {
	mutex.Lock()
	defer mutex.Unlock()
	return fmt.Sprintf(":%d", f.data.WebServerPort)
}

// GetBotToken returns the Telegram Bot TOKEN
func (f *File) GetBotToken() string {
	mutex.Lock()
	defer mutex.Unlock()
	return f.data.TelegramBotToken
}

// GetChatIDs returns the Telegram Bot TOKEN
func (f *File) GetChatIDs() []int64 {
	mutex.Lock()
	defer mutex.Unlock()
	return f.data.TelegramChatIDs
}
