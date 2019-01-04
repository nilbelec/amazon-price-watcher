package configuration

// Settings stores the configuration settings values
type Settings struct {
	WebServerPort                    int
	ProductsRefreshIntervalInMinutes int
	TelegramBotToken                 string
	TelegramChatIDs                  []int64
}

// Defaults contains the default configuration settings values
var Defaults = &Settings{
	WebServerPort:                    10035,
	ProductsRefreshIntervalInMinutes: 5,
	TelegramBotToken:                 "",
	TelegramChatIDs:                  make([]int64, 0),
}
