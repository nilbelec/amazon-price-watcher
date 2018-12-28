package file

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"sync"

	"github.com/nilbelec/amazon-price-watcher/pkg/file"
	"github.com/nilbelec/amazon-price-watcher/pkg/model"
)

// ProductRepository struct
type ProductRepository struct {
	products     map[string]model.Product
	productsFile string
}

var mutex = &sync.Mutex{}

// New creates a new file-based ProductRepository
func New(productsFile string) (ps *ProductRepository, err error) {
	products, err := loadProducts(productsFile)
	if err != nil {
		err = fmt.Errorf("Error loading products from file: %s", err.Error())
		return
	}
	ps = &ProductRepository{products, productsFile}
	return
}

func (r *ProductRepository) saveProducts() (err error) {
	data, err := json.MarshalIndent(&r.products, "", " ")
	if err != nil {
		return
	}
	err = ioutil.WriteFile(r.productsFile, data, 0644)
	if err != nil {
		return
	}
	return
}

func loadProducts(productsFile string) (products map[string]model.Product, err error) {
	exists, err := file.Exists(productsFile)
	if err != nil {
		err = fmt.Errorf("Error checking if config file exists: %s", err.Error())
		return
	}
	if exists == false {
		products = make(map[string]model.Product)
		return
	}
	file, err := os.Open(productsFile)
	defer file.Close()
	if err != nil {
		err = fmt.Errorf("Error opening products file: %s", err.Error())
		return
	}
	decoder := json.NewDecoder(file)
	products = make(map[string]model.Product)
	err = decoder.Decode(&products)
	if err != nil {
		err = fmt.Errorf("Error parsing json from products file: %s", err.Error())
		return
	}
	return
}

// AddProduct stores a new Product
func (r *ProductRepository) AddProduct(product model.Product) error {
	mutex.Lock()
	defer mutex.Unlock()
	if _, ok := r.products[product.URL]; ok {
		return errors.New("Product already exists")
	}
	r.products[product.URL] = product
	err := r.saveProducts()
	if err != nil {
		delete(r.products, product.URL)
		return errors.New("Error saving products: " + err.Error())
	}
	return nil
}

// DeleteProductByURL removes a product by its URL
func (r *ProductRepository) DeleteProductByURL(url string) (product model.Product, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	product, ok := r.products[url]
	if !ok {
		err = errors.New("Product doesnt exist")
		return
	}
	delete(r.products, url)
	err = r.saveProducts()
	if err != nil {
		r.products[url] = product
		err = errors.New("Error saving products: " + err.Error())
		return
	}
	return
}

// UpdateProduct updates an existing product
func (r *ProductRepository) UpdateProduct(product model.Product) error {
	mutex.Lock()
	defer mutex.Unlock()
	old, ok := r.products[product.URL]
	if !ok {
		return errors.New("Product is not stored")
	}
	product.Added = old.Added
	r.products[product.URL] = product
	err := r.saveProducts()
	if err != nil {
		r.products[product.URL] = old
		return errors.New("Error saving products: " + err.Error())
	}
	return nil
}

// ListProducts list all products
func (r *ProductRepository) ListProducts() ([]model.Product, error) {
	mutex.Lock()
	defer mutex.Unlock()
	v := make([]model.Product, 0, len(r.products))
	for _, value := range r.products {
		v = append(v, value)
	}
	sort.Slice(v, func(i, j int) bool {
		return (v[i].Price < v[i].LastPrice && v[j].Price >= v[j].LastPrice) || v[i].Added.After(v[j].Added)
	})
	return v, nil
}

// GetProductByURL gets an existing product by its URL
func (r *ProductRepository) GetProductByURL(url string) (product model.Product, err error) {
	mutex.Lock()
	defer mutex.Unlock()
	product, ok := r.products[url]
	if !ok {
		err = errors.New("Product is not stored")
	}
	return
}
