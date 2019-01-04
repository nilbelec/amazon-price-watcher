package web

import (
	"time"

	"github.com/nilbelec/amazon-price-watcher/pkg/product"
)

// Product JSON representation of a Product
type Product struct {
	URL           string
	Title         string
	Price         float32
	LastPrice     float32
	Currency      string
	ImageURL      string
	LastUpdate    time.Time
	Notifications Notifications `json:"notifications"`
}

// Notifications JSON representation of a Product notifications
type Notifications struct {
	PriceBelows    float32 `json:"price_below"`
	PriceOver      float32 `json:"price_over"`
	PriceIncreases bool    `json:"price_increases"`
	PriceDecreases bool    `json:"price_decreases"`
	OutOfStock     bool    `json:"out_of_stock"`
	BackInStock    bool    `json:"back_in_stock"`
	Total          int     `json:"total"`
}

// ToSliceOfJSONS Convert Products to JSON
func ToSliceOfJSONS(products []product.Product) []Product {
	jsons := make([]Product, 0, cap(products))
	for _, p := range products {
		json := toJSON(p)
		jsons = append(jsons, json)
	}
	return jsons
}

func toJSON(product product.Product) Product {
	return Product{
		URL:        product.URL,
		Title:      product.Title,
		ImageURL:   product.ImageURL,
		Price:      product.Price,
		Currency:   product.Currency,
		LastPrice:  product.LastPrice,
		LastUpdate: product.LastUpdate,
		Notifications: Notifications{
			PriceBelows:    product.Notifications.PriceBelows,
			PriceOver:      product.Notifications.PriceOver,
			PriceIncreases: product.Notifications.PriceIncreases,
			PriceDecreases: product.Notifications.PriceDecreases,
			OutOfStock:     product.Notifications.OutOfStock,
			BackInStock:    product.Notifications.BackInStock,
			Total:          product.Notifications.GetTotal(),
		}}
}
