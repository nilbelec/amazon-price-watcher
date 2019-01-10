package telegram

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/nilbelec/amazon-price-watcher/pkg/product"
)

// Configuration provides the Telegram configuration
type Configuration struct {
	BotToken func() string
	ChatIDs  func() []int64
}

// Notifier is the Telegram notifier
type Notifier struct {
	config *Configuration
}

// NewNotifier creates a new Telegram product notifier
func NewNotifier(c *Configuration) (n *Notifier) {
	n = &Notifier{c}
	return
}

// NotifyChanges sends a Telegram bot message notifying a product change
func (n *Notifier) NotifyChanges(p *product.Product) {
	if !n.IsConfigured() || !p.ShouldSendAnyNotification() {
		return
	}
	bot, err := tgbotapi.NewBotAPI(n.config.BotToken())
	if err != nil {
		log.Printf("Error preparing Telegram Bot: " + err.Error())
		return
	}
	msgs := n.prepareNotificationMessages(p)
	n.sendNotifications(bot, msgs, n.config.ChatIDs())
	return
}

func (n *Notifier) prepareNotificationMessages(p *product.Product) []string {
	msgs := make([]string, 0)
	if p.ShouldSendPriceDecreasesNotification() {
		msgs = append(msgs, fmt.Sprintf("'%s' is now at %.2f %s (before: %.2f %s)", p.Title, p.Price, p.Currency, p.LastPrice, p.Currency))
	}
	if p.ShouldSendPriceIncreasesNotification() {
		msgs = append(msgs, fmt.Sprintf("'%s' is now at %.2f %s (before: %.2f %s)", p.Title, p.Price, p.Currency, p.LastPrice, p.Currency))
	}
	if p.ShouldSendPriceBelowsNotification() {
		msgs = append(msgs, fmt.Sprintf("Cool! '%s' is below %.2f %s! It's now at %.2f %s (before: %.2f %s)", p.Title, p.Notifications.PriceBelows, p.Currency, p.Price, p.Currency, p.LastPrice, p.Currency))
	}
	if p.ShouldSendPriceOverNotification() {
		msgs = append(msgs, fmt.Sprintf("'%s' is over %.2f %s... It's at %.2f %s (before: %.2f %s)", p.Title, p.Notifications.PriceOver, p.Currency, p.Price, p.Currency, p.LastPrice, p.Currency))
	}
	if p.ShouldSendBackInStockNotification() {
		msgs = append(msgs, fmt.Sprintf("'%s' is back in stock! It's now at %.2f %s", p.Title, p.Price, p.Currency))
	}
	if p.ShouldSendOutOfStockNotification() {
		msgs = append(msgs, fmt.Sprintf("'%s' is out of stock... Before it was at %.2f %s ", p.Title, p.LastPrice, p.Currency))
	}
	return msgs
}

func (n *Notifier) sendNotifications(bot *tgbotapi.BotAPI, msgs []string, ids []int64) {
	if len(msgs) == 0 || len(ids) == 0 {
		return
	}
	for _, msg := range msgs {
		for _, id := range ids {
			m := tgbotapi.NewMessage(id, msg)
			_, err := bot.Send(m)
			if err != nil {
				log.Printf("Error sending message to chatID %d: %s\n", id, err.Error())
			}
		}
	}
}

// IsConfigured tells if the Telegram Notifier is configured
func (n *Notifier) IsConfigured() bool {
	return n.config.BotToken() != "" && len(n.config.ChatIDs()) != 0
}
