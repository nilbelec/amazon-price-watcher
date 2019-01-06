package product

import (
	"time"

	"github.com/nilbelec/amazon-price-watcher/pkg/product"
)

type productJSON struct {
	URL           string
	Title         string
	Price         float32
	LastPrice     float32
	Currency      string
	ImageURL      string
	LastUpdate    time.Time
	Notifications *notificationsJSON `json:"notifications"`
}

type productsJSON []*productJSON

type notificationsJSON struct {
	PriceBelows    float32 `json:"price_below"`
	PriceOver      float32 `json:"price_over"`
	PriceIncreases bool    `json:"price_increases"`
	PriceDecreases bool    `json:"price_decreases"`
	OutOfStock     bool    `json:"out_of_stock"`
	BackInStock    bool    `json:"back_in_stock"`
	Total          int     `json:"total"`
}

func toJSONs(ps *product.Products) *productsJSON {
	jsons := make(productsJSON, 0, cap(*ps))
	for _, p := range *ps {
		jsons = append(jsons, toJSON(p))
	}
	return &jsons
}

func toJSON(p *product.Product) *productJSON {
	return &productJSON{
		URL:        p.URL,
		Title:      p.Title,
		ImageURL:   p.ImageURL,
		Price:      p.Price,
		Currency:   p.Currency,
		LastPrice:  p.LastPrice,
		LastUpdate: p.LastUpdate,
		Notifications: &notificationsJSON{
			PriceBelows:    p.Notifications.PriceBelows,
			PriceOver:      p.Notifications.PriceOver,
			PriceIncreases: p.Notifications.PriceIncreases,
			PriceDecreases: p.Notifications.PriceDecreases,
			OutOfStock:     p.Notifications.OutOfStock,
			BackInStock:    p.Notifications.BackInStock,
			Total:          p.Notifications.GetTotal(),
		}}
}
