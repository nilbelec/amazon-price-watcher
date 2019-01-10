package product

import (
	"errors"
	"log"
	"time"
)

// Configuration handles the products configuration
type Configuration struct {
	RefreshIntervalMinutes func() int
}

// Crawler handles the products crawling
type Crawler interface {
	Extract(url string) (*Product, error)
}

// Notifier handles the products notifications
type Notifier interface {
	NotifyChanges(p *Product)
	IsConfigured() bool
}

// Notifiers is a slice of Notifier
type Notifiers []Notifier

// Repository handles the products persistence
type Repository interface {
	Add(p *Product) error
	Delete(url string) (*Product, error)
	List() (*Products, error)
	Update(p *Product) error
	Get(url string) (*Product, error)
}

// Service struct
type Service struct {
	repo      Repository
	crawler   Crawler
	notifiers *Notifiers
	config    *Configuration
}

// NewService creates a new Product Service
func NewService(r Repository, c Crawler, ns *Notifiers, cfg *Configuration) (ps *Service) {
	ps = &Service{r, c, ns, cfg}
	go ps.refreshProducts()
	return
}

func (ps *Service) refreshProducts() {
	for {
		interval := time.Duration(ps.config.RefreshIntervalMinutes()) * time.Minute
		prs, err := ps.repo.List()
		if err != nil {
			err = errors.New("Error while retrieving products: " + err.Error())
			time.Sleep(interval)
			continue
		}
		for _, p := range *prs {
			go ps.refreshProduct(p)
		}
		time.Sleep(interval)
	}
}

func (ps *Service) refreshProduct(p *Product) {
	a, err := ps.crawler.Extract(p.URL)
	if err != nil {
		log.Println("Unable to refresh data for '" + p.URL + "': " + err.Error())
		return
	}
	p.UpdateInfo(a)
	err = ps.repo.Update(p)
	if err != nil {
		log.Println("Unable to update product '" + p.URL + "': " + err.Error())
		return
	}
	if p.ShouldSendAnyNotification() {
		ps.notifyProductChange(p)
	}
}

func (ps *Service) notifyProductChange(p *Product) {
	for _, n := range *ps.notifiers {
		if n.IsConfigured() {
			n.NotifyChanges(p)
		}
	}
}

// AddProductByURL adds a new Amazon product by its URL
func (ps *Service) AddProductByURL(url string) (p *Product, err error) {
	p, err = ps.crawler.Extract(url)
	if err != nil {
		return
	}
	p.Added = time.Now()
	p.LastPrice = p.Price
	err = ps.repo.Add(p)
	return
}

// UpdateNotifications updates the product notifications
func (ps *Service) UpdateNotifications(url string, n *Notifications) (err error) {
	p, err := ps.repo.Get(url)
	if err != nil {
		return err
	}
	p.Notifications = n
	err = ps.repo.Update(p)
	if err != nil {
		return err
	}
	return nil
}

// DeleteProductByURL deletes a product by its URL
func (ps *Service) DeleteProductByURL(url string) (*Product, error) {
	return ps.repo.Delete(url)
}

// ListProducts list current saved products
func (ps *Service) ListProducts() (*Products, error) {
	return ps.repo.List()
}
