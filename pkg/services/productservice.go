package services

import (
	"errors"
	"log"
	"time"

	"github.com/nilbelec/amazon-price-watcher/pkg/model"
)

//ProductRepository interface
type ProductRepository interface {
	AddProduct(product model.Product) error
	DeleteProductByURL(url string) (model.Product, error)
	ListProducts() ([]model.Product, error)
	UpdateProduct(product model.Product) error
	GetProductByURL(url string) (model.Product, error)
}

//ProductFinder interface
type ProductFinder interface {
	FindByURL(url string) (model.Product, error)
}

//ProductsConfiguration handles the products configuration
type ProductsConfiguration interface {
	GetRefreshInterval() time.Duration
}

//ProductNotifier handles the products notifications
type ProductNotifier interface {
	NotifyProductChange(product model.Product)
	IsConfigured() bool
}

// ProductService struct
type ProductService struct {
	repo      ProductRepository
	finder    ProductFinder
	conf      ProductsConfiguration
	notifiers []ProductNotifier
}

// NewProductService Creates a new Product Service
func NewProductService(repo ProductRepository, finder ProductFinder, conf ProductsConfiguration, notifiers []ProductNotifier) (ps *ProductService) {
	ps = &ProductService{repo, finder, conf, notifiers}
	go ps.refreshProducts()
	return
}

func (ps *ProductService) refreshProducts() {
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

func (ps *ProductService) refreshProduct(product model.Product) {
	actual, err := ps.finder.FindByURL(product.URL)
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

func (ps *ProductService) notifyProductChange(product model.Product) {
	for _, notifier := range ps.notifiers {
		if notifier.IsConfigured() {
			notifier.NotifyProductChange(product)
		}
	}
}

// AddProductByURL adds a new Amazon product by its URL
func (ps *ProductService) AddProductByURL(url string) (product model.Product, err error) {
	product, err = ps.finder.FindByURL(url)
	if err != nil {
		return
	}
	product.Added = time.Now()
	product.LastPrice = product.Price
	err = ps.repo.AddProduct(product)
	return
}

// UpdateProductNotifications updates the product notifications
func (ps *ProductService) UpdateProductNotifications(url string, notifications model.Notifications) (err error) {
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
func (ps *ProductService) DeleteProductByURL(url string) (model.Product, error) {
	return ps.repo.DeleteProductByURL(url)
}

// ListProducts list current saved products
func (ps *ProductService) ListProducts() ([]model.Product, error) {
	return ps.repo.ListProducts()
}
