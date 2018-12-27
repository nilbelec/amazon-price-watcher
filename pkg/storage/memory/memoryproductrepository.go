package memory

import (
	"errors"
	"sort"
	"sync"
	"time"

	"github.com/nilbelec/amazon-price-watcher/pkg/model"
)

// InMemoryProductRepository struct
type InMemoryProductRepository struct {
	products map[string]model.Product
}

var mutex = &sync.Mutex{}

// NewInMemoryProductRepository creates a new InMemoryProductRepository
func NewInMemoryProductRepository() *InMemoryProductRepository {
	products := make(map[string]model.Product)
	return &InMemoryProductRepository{products}
}

// AddProduct stores a new Product
func (r *InMemoryProductRepository) AddProduct(product model.Product) error {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := r.products[product.URL]; ok {
		return errors.New("Product already saved")
	}
	product.Added = time.Now()
	r.products[product.URL] = product
	return nil
}

// DeleteProductByURL removes a product by its URL
func (r *InMemoryProductRepository) DeleteProductByURL(url string) (product model.Product, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	product, ok := r.products[url]
	if !ok {
		err = errors.New("Product is not stored")
		return
	}
	delete(r.products, url)
	return
}

// UpdateProduct updates an existing product
func (r *InMemoryProductRepository) UpdateProduct(product model.Product) error {
	mutex.Lock()
	defer mutex.Unlock()
	_, ok := r.products[product.URL]
	if !ok {
		return errors.New("Product is not stored")
	}
	r.products[product.URL] = product
	return nil
}

// ListProducts list all products
func (r *InMemoryProductRepository) ListProducts() ([]model.Product, error) {
	mutex.Lock()
	defer mutex.Unlock()
	v := make([]model.Product, 0, len(r.products))
	for _, value := range r.products {
		v = append(v, value)
	}
	sort.Slice(v, func(i, j int) bool {
		return v[i].Added.After(v[j].Added)
	})
	return v, nil
}
