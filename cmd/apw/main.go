package main

import (
	"log"

	cf "github.com/nilbelec/amazon-price-watcher/pkg/configuration/file"
	amazon "github.com/nilbelec/amazon-price-watcher/pkg/crawler/amazon"
	"github.com/nilbelec/amazon-price-watcher/pkg/services"
	"github.com/nilbelec/amazon-price-watcher/pkg/storage/file"
	"github.com/nilbelec/amazon-price-watcher/pkg/telegram"
	"github.com/nilbelec/amazon-price-watcher/pkg/web"
)

const configFile = "config.json"
const productsFile = "products.json"

func main() {
	config, err := cf.New(configFile)
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
	crawler := amazon.New()
	ps := services.NewProductService(repo, crawler, config, notifiers)

	web := web.NewServer(ps, config, version)
	web.Start()
}
