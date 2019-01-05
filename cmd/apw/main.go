package main

import (
	"log"

	"github.com/nilbelec/amazon-price-watcher/pkg/amazon"
	"github.com/nilbelec/amazon-price-watcher/pkg/configuration"
	"github.com/nilbelec/amazon-price-watcher/pkg/product"
	"github.com/nilbelec/amazon-price-watcher/pkg/storage"
	"github.com/nilbelec/amazon-price-watcher/pkg/telegram"
	"github.com/nilbelec/amazon-price-watcher/pkg/web"
)

const configFile = "config.json"
const productsFile = "products.json"

func main() {
	cs := configurationService()
	ps := productsService(cs)
	s := web.NewServer(cs, ps, cs)
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

func notifiers(bc telegram.BotConfig) []product.Notifier {
	ns := make([]product.Notifier, 0)
	tn, err := telegram.NewNotifier(bc)
	if err != nil {
		log.Fatalln("Error preparing the Telegram notifier: " + err.Error())
	}
	ns = append(ns, tn)
	return ns
}

func productsService(cs *configuration.Service) *product.Service {
	ns := notifiers(cs)
	ac := amazon.NewCrawler()

	pf, err := storage.NewProductsFile(productsFile)
	if err != nil {
		log.Fatalln("Error loading products file: " + err.Error())
	}
	ps := product.NewService(pf, ac, cs, ns)
	return ps
}
