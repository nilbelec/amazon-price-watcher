package configuration

// Settings stores the configuration settings values
type Settings struct {
	WebServerPort                    int
	ProductsRefreshIntervalInMinutes int
	TelegramBotToken                 string
	TelegramChatIDs                  []int64
	SMTPHost                         string
	SMTPPort                         int
	SMTPUsername                     string
	SMTPPassword                     string
	SMTPTo                           []string
}

// Defaults contains the default configuration settings values
var Defaults = &Settings{
	WebServerPort:                    10035,
	ProductsRefreshIntervalInMinutes: 5,
	TelegramBotToken:                 "",
	TelegramChatIDs:                  make([]int64, 0),
	SMTPHost:                         "",
	SMTPPort:                         0,
	SMTPUsername:                     "",
	SMTPPassword:                     "",
	SMTPTo:                           make([]string, 0),
}
