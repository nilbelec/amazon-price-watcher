package product

import (
	"errors"
	"log"
	"time"
)

// Repository handles the products persistence
type Repository interface {
	AddProduct(product Product) error
	DeleteProductByURL(url string) (Product, error)
	ListProducts() ([]Product, error)
	UpdateProduct(product Product) error
	GetProductByURL(url string) (Product, error)
}

// Crawler handles the products crawling
type Crawler interface {
	FindByURL(url string) (Product, error)
}

// Configuration handles the products configuration
type Configuration interface {
	GetRefreshInterval() time.Duration
}

// Notifier handles the products notifications
type Notifier interface {
	NotifyProductChange(product Product)
	IsConfigured() bool
}

// Service struct
type Service struct {
	repo      Repository
	crawler   Crawler
	conf      Configuration
	notifiers []Notifier
}

// New creates a new Product Service
func New(repo Repository, finder Crawler, conf Configuration, notifiers []Notifier) (ps *Service) {
	ps = &Service{repo, finder, conf, notifiers}
	go ps.refreshProducts()
	return
}

func (ps *Service) refreshProducts() {
	for {
		products, err := ps.repo.ListProducts()
		if err != nil {
			err = errors.New("Error while retrieving products: " + err.Error())
			time.Sleep(ps.conf.GetRefreshInterval())
			continue
		}
		for _, p := range products {
			go ps.refreshProduct(p)
		}
		time.Sleep(ps.conf.GetRefreshInterval())
	}
}

func (ps *Service) refreshProduct(product Product) {
	actual, err := ps.crawler.FindByURL(product.URL)
	if err != nil {
		log.Println("Unable to refresh data for '" + product.URL + "': " + err.Error())
		return
	}
	product.UpdateInfo(actual)
	err = ps.repo.UpdateProduct(product)
	if err != nil {
		log.Println("Unable to update product '" + product.URL + "': " + err.Error())
		return
	}
	if product.ShouldSendAnyNotification() {
		ps.notifyProductChange(product)
	}
}

func (ps *Service) notifyProductChange(product Product) {
	for _, notifier := range ps.notifiers {
		if notifier.IsConfigured() {
			notifier.NotifyProductChange(product)
		}
	}
}

// AddProductByURL adds a new Amazon product by its URL
func (ps *Service) AddProductByURL(url string) (product Product, err error) {
	product, err = ps.crawler.FindByURL(url)
	if err != nil {
		return
	}
	product.Added = time.Now()
	product.LastPrice = product.Price
	err = ps.repo.AddProduct(product)
	return
}

// UpdateProductNotifications updates the product notifications
func (ps *Service) UpdateProductNotifications(url string, notifications Notifications) (err error) {
	product, err := ps.repo.GetProductByURL(url)
	if err != nil {
		return err
	}
	product.Notifications = notifications
	err = ps.repo.UpdateProduct(product)
	if err != nil {
		return err
	}
	return nil
}

// DeleteProductByURL deletes a product by its URL
func (ps *Service) DeleteProductByURL(url string) (Product, error) {
	return ps.repo.DeleteProductByURL(url)
}

// ListProducts list current saved products
func (ps *Service) ListProducts() ([]Product, error) {
	return ps.repo.ListProducts()
}
