package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nilbelec/amazon-price-watcher/pkg/model"
)

// BotConfig provides the telegram configuration
type BotConfig interface {
	GetBotToken() string
	GetChatIDs() []int64
}

// Notifier is a telegram product notifier
type Notifier struct {
	config BotConfig
}

// New creates a new telegram product notifier
func New(config BotConfig) (n *Notifier, err error) {
	n = &Notifier{config}
	return
}

// NotifyProductChange send a telegram bot message notifying a product change
func (n *Notifier) NotifyProductChange(product model.Product) {
	if !n.IsEnabled() || !product.Changed() {
		return
	}
	bot, err := tgbotapi.NewBotAPI(n.config.GetBotToken())
	if err != nil {
		log.Printf("Error preparing Telegram Bot: " + err.Error())
		return
	}

	var text string
	if product.PriceHasDecreased() {
		text = fmt.Sprintf("Cool! '%s' is now at %.2f %s (before: %.2f %s)", product.Title, product.Price, product.Currency, product.LastPrice, product.Currency)
	} else if product.PriceHasIncrecreased() {
		text = fmt.Sprintf("Oh no! '%s' is now at %.2f %s (before: %.2f %s)", product.Title, product.Price, product.Currency, product.LastPrice, product.Currency)
	} else {
		return
	}

	for _, chatID := range n.config.GetChatIDs() {
		msg := tgbotapi.NewMessage(chatID, text)
		_, err := bot.Send(msg)
		if err != nil {
			log.Printf("Error sending message to chatID %d: %s\n", chatID, err.Error())
		}
	}
	return
}

// IsEnabled tells if the Telegram Notifier is enabled
func (n *Notifier) IsEnabled() bool {
	return n.config.GetBotToken() != "" && len(n.config.GetChatIDs()) != 0
}
