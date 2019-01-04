package main

import (
	"log"

	cf "github.com/nilbelec/amazon-price-watcher/pkg/configuration/file"
	"github.com/nilbelec/amazon-price-watcher/pkg/crawler/amazon"
	"github.com/nilbelec/amazon-price-watcher/pkg/notifier/telegram"
	"github.com/nilbelec/amazon-price-watcher/pkg/product"
	"github.com/nilbelec/amazon-price-watcher/pkg/storage/file"
	"github.com/nilbelec/amazon-price-watcher/pkg/web"
)

const configFile = "config.json"
const productsFile = "products.json"

func main() {
	config, err := cf.New(configFile)
	if err != nil {
		log.Fatalln("Error loading configuration file: " + err.Error())
	}
	notifiers := make([]product.Notifier, 0)
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
	ps := product.New(repo, crawler, config, notifiers)

	web := web.NewServer(config, ps, config)
	web.Start()
}
