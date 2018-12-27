package main

import (
	"log"

	"github.com/nilbelec/amazon-price-watcher/pkg/configuration"
	"github.com/nilbelec/amazon-price-watcher/pkg/crawler"
	"github.com/nilbelec/amazon-price-watcher/pkg/services"
	"github.com/nilbelec/amazon-price-watcher/pkg/storage/file"
	"github.com/nilbelec/amazon-price-watcher/pkg/telegram"
	"github.com/nilbelec/amazon-price-watcher/pkg/web"
)

const configFile = "config.json"
const productsFile = "products.json"

func main() {
	config, err := configuration.Load(configFile)
	if err != nil {
		log.Fatalln("Error loading configuration file: " + err.Error())
	}
	notifiers := make([]services.ProductNotifier, 0)
	telegram, err := telegram.New(config)
	if err != nil {
		log.Fatalln("Error preparing the Telegram notifier: " + err.Error())
	}
	notifiers = append(notifiers, telegram)
	repo, err := file.New(productsFile)
	if err != nil {
		log.Fatalln("Error loading products file: " + err.Error())
	}
	finder := crawler.NewProductCrawler()
	ps := services.NewProductService(repo, finder, config, notifiers)

	web := web.NewServer(ps, config)
	web.Start()
}
