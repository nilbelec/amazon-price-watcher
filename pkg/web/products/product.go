package products

import (
	"time"

	"github.com/nilbelec/amazon-price-watcher/pkg/model"
)

// Product JSON representation of a Product
type Product struct {
	URL        string
	Title      string
	Price      float32
	LastPrice  float32
	Currency   string
	ImageURL   string
	LastUpdate time.Time
}

// ToSliceOfJSONS Convert Products to JSON
func ToSliceOfJSONS(products []model.Product) []Product {
	jsons := make([]Product, 0, cap(products))
	for _, p := range products {
		json := toJSON(p)
		jsons = append(jsons, json)
	}
	return jsons
}

func toJSON(product model.Product) Product {
	return Product{
		URL:        product.URL,
		Title:      product.Title,
		ImageURL:   product.ImageURL,
		Price:      product.Price,
		Currency:   product.Currency,
		LastPrice:  product.LastPrice,
		LastUpdate: product.LastUpdate}
}
