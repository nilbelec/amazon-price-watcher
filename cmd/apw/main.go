package main

import (
	"log"

	"github.com/nilbelec/amazon-price-watcher/pkg/github"
	"github.com/nilbelec/amazon-price-watcher/pkg/version"

	"github.com/nilbelec/amazon-price-watcher/pkg/amazon"
	"github.com/nilbelec/amazon-price-watcher/pkg/configuration"
	"github.com/nilbelec/amazon-price-watcher/pkg/product"
	"github.com/nilbelec/amazon-price-watcher/pkg/storage"
	"github.com/nilbelec/amazon-price-watcher/pkg/telegram"
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
	gh := github.NewClient(gitHubUser, gitHubRepo)
	vs := version.NewService(gh)
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

func notifiers(bc telegram.Configuration) *product.Notifiers {
	ns := make(product.Notifiers, 0)
	tn := telegram.NewNotifier(bc)
	ns = append(ns, tn)
	return &ns
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
