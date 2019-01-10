package main

import (
	"log"

	"github.com/nilbelec/amazon-price-watcher/pkg/smtp"
	"github.com/nilbelec/amazon-price-watcher/pkg/telegram"

	"github.com/nilbelec/amazon-price-watcher/pkg/github"
	"github.com/nilbelec/amazon-price-watcher/pkg/version"

	"github.com/nilbelec/amazon-price-watcher/pkg/amazon"
	"github.com/nilbelec/amazon-price-watcher/pkg/configuration"
	"github.com/nilbelec/amazon-price-watcher/pkg/product"
	"github.com/nilbelec/amazon-price-watcher/pkg/storage"
	"github.com/nilbelec/amazon-price-watcher/pkg/web"
)

const (
	configFile   = "config.json"
	productsFile = "products.json"
	gitHubUser   = "nilbelec"
	gitHubRepo   = "amazon-price-watcher"
)

func main() {
	cs := configurationService()
	ps := productsService(cs)
	vs := versionService()
	s := web.NewServer(ps, cs, vs)
	s.Start()
}

func configurationService() *configuration.Service {
	cf, err := storage.NewConfigurationFile(configFile)
	if err != nil {
		log.Fatalln("Error loading configuration file: " + err.Error())
	}
	cs, err := configuration.NewService(cf)
	if err != nil {
		log.Fatalln("Error loading configuration service: " + err.Error())
	}
	return cs
}

func versionService() *version.Service {
	gh := github.NewClient(gitHubUser, gitHubRepo)
	return version.NewService(gh)
}

func notifiers(cs *configuration.Service) *product.Notifiers {
	tn := telegramNotifier(cs)
	mn := smtpNotifier(cs)
	return &product.Notifiers{tn, mn}
}

func productsService(cs *configuration.Service) *product.Service {
	ns := notifiers(cs)
	ac := amazon.NewCrawler()

	pf, err := storage.NewProductsFile(productsFile)
	if err != nil {
		log.Fatalln("Error loading products file: " + err.Error())
	}
	pc := &product.Configuration{RefreshIntervalMinutes: cs.ProductsRefreshIntervalInMinutes}
	ps := product.NewService(pf, ac, ns, pc)
	return ps
}

func telegramNotifier(cs *configuration.Service) *telegram.Notifier {
	c := &telegram.Configuration{
		BotToken: cs.TelegramBotToken,
		ChatIDs:  cs.TelegramChatIDs,
	}
	return telegram.NewNotifier(c)
}

func smtpNotifier(cs *configuration.Service) *smtp.Notifier {
	c := &smtp.Configuration{
		Host:     cs.SMTPHost,
		Port:     cs.SMTPPort,
		Username: cs.SMTPUsername,
		Password: cs.SMTPPassword,
		To:       cs.SMTPTo,
	}
	return smtp.NewNotifier(c)
}
