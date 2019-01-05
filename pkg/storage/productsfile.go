package storage

import (
	"errors"
	"sort"
	"sync"

	"github.com/nilbelec/amazon-price-watcher/pkg/product"
)

// ProductsFile struct
type ProductsFile struct {
	sync.Mutex
	filename string
	products map[string]product.Product
}

// NewProductsFile creates a new ProductsFile to store the products
func NewProductsFile(filename string) (pf *ProductsFile, err error) {
	pf = &ProductsFile{filename: filename}
	err = pf.load()
	return
}

func (pf *ProductsFile) persists() (err error) {
	pf.Lock()
	defer pf.Unlock()
	return SaveJSON(pf.filename, &pf.products)
}

func (pf *ProductsFile) load() (err error) {
	pf.Lock()
	defer pf.Unlock()
	return ReadJSON(pf.filename, &pf.products)
}

// Add stores a new product
func (pf *ProductsFile) Add(p product.Product) error {
	if _, ok := pf.products[p.URL]; ok {
		return errors.New("The product is already on your watchlist")
	}
	pf.products[p.URL] = p
	err := pf.persists()
	if err != nil {
		delete(pf.products, p.URL)
		return errors.New("Error saving products: " + err.Error())
	}
	return nil
}

// Delete removes a product by using its URL
func (pf *ProductsFile) Delete(url string) (p product.Product, err error) {
	p, ok := pf.products[url]
	if !ok {
		err = errors.New("The product is not on your watchlist")
		return
	}
	delete(pf.products, url)
	err = pf.persists()
	if err != nil {
		pf.products[url] = p
		err = errors.New("Error saving products: " + err.Error())
		return
	}
	return
}

// Update updates an existing product
func (pf *ProductsFile) Update(p product.Product) error {
	old, ok := pf.products[p.URL]
	if !ok {
		return errors.New("The product is not on your watchlist")
	}
	p.Added = old.Added
	pf.products[p.URL] = p
	err := pf.persists()
	if err != nil {
		pf.products[p.URL] = old
		return errors.New("Error updating the product: " + err.Error())
	}
	return nil
}

// List lists all products
func (pf *ProductsFile) List() ([]product.Product, error) {
	p := make([]product.Product, 0, len(pf.products))
	for _, value := range pf.products {
		p = append(p, value)
	}
	sort.Slice(p, func(i, j int) bool {
		if p[i].Price < p[i].LastPrice && p[j].Price >= p[j].LastPrice {
			return true
		}
		if p[j].Price < p[j].LastPrice && p[i].Price >= p[i].LastPrice {
			return false
		}
		return p[i].Added.After(p[j].Added)
	})
	return p, nil
}

// Get gets an existing product by its URL
func (pf *ProductsFile) Get(url string) (p product.Product, err error) {
	p, ok := pf.products[url]
	if !ok {
		err = errors.New("The product is not on your watchlist")
	}
	return
}
