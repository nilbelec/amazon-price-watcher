package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nilbelec/amazon-price-watcher/pkg/product"
)

// BotConfig provides the Telegram configuration
type BotConfig interface {
	GetBotToken() string
	GetChatIDs() []int64
}

// Notifier is the Telegram notifier
type Notifier struct {
	config BotConfig
}

// New creates a new Telegram product notifier
func New(config BotConfig) (n *Notifier, err error) {
	n = &Notifier{config}
	return
}

// NotifyChanges sends a Telegram bot message notifying a product change
func (n *Notifier) NotifyChanges(product product.Product) {
	if !n.IsConfigured() || !product.ShouldSendAnyNotification() {
		return
	}
	bot, err := tgbotapi.NewBotAPI(n.config.GetBotToken())
	if err != nil {
		log.Printf("Error preparing Telegram Bot: " + err.Error())
		return
	}

	messages := prepareNotificationMessages(product)
	sendNotifications(bot, messages, n.config.GetChatIDs())
	return
}

func prepareNotificationMessages(product product.Product) []string {
	messages := make([]string, 0)
	if product.ShouldSendPriceDecreasesNotification() {
		messages = append(messages, fmt.Sprintf("'%s' is now at %.2f %s (before: %.2f %s)", product.Title, product.Price, product.Currency, product.LastPrice, product.Currency))
	}
	if product.ShouldSendPriceIncreasesNotification() {
		messages = append(messages, fmt.Sprintf("'%s' is now at %.2f %s (before: %.2f %s)", product.Title, product.Price, product.Currency, product.LastPrice, product.Currency))
	}
	if product.ShouldSendPriceBelowsNotification() {
		messages = append(messages, fmt.Sprintf("Cool! '%s' is below %.2f %s! It's now at %.2f %s (before: %.2f %s)", product.Title, product.Notifications.PriceBelows, product.Currency, product.Price, product.Currency, product.LastPrice, product.Currency))
	}
	if product.ShouldSendPriceOverNotification() {
		messages = append(messages, fmt.Sprintf("'%s' is over %.2f %s... It's at %.2f %s (before: %.2f %s)", product.Title, product.Notifications.PriceOver, product.Currency, product.Price, product.Currency, product.LastPrice, product.Currency))
	}
	if product.ShouldSendBackInStockNotification() {
		messages = append(messages, fmt.Sprintf("'%s' is back in stock! It's now at %.2f %s", product.Title, product.Price, product.Currency))
	}
	if product.ShouldSendOutOfStockNotification() {
		messages = append(messages, fmt.Sprintf("'%s' is out of stock... Before it was at %.2f %s ", product.Title, product.LastPrice, product.Currency))
	}
	return messages
}

func sendNotifications(bot *tgbotapi.BotAPI, messages []string, chatIDs []int64) {
	if len(messages) == 0 || len(chatIDs) == 0 {
		return
	}
	for _, text := range messages {
		for _, chatID := range chatIDs {
			msg := tgbotapi.NewMessage(chatID, text)
			_, err := bot.Send(msg)
			if err != nil {
				log.Printf("Error sending message to chatID %d: %s\n", chatID, err.Error())
			}
		}
	}
}

// IsConfigured tells if the Telegram Notifier is configured
func (n *Notifier) IsConfigured() bool {
	return n.config.GetBotToken() != "" && len(n.config.GetChatIDs()) != 0
}
