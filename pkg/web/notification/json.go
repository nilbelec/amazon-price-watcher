package notification

import "github.com/nilbelec/amazon-price-watcher/pkg/product"

type notificationsJSON struct {
	PriceBelows    float32 `json:"price_below"`
	PriceOver      float32 `json:"price_over"`
	PriceIncreases bool    `json:"price_increases"`
	PriceDecreases bool    `json:"price_decreases"`
	OutOfStock     bool    `json:"out_of_stock"`
	BackInStock    bool    `json:"back_in_stock"`
}

func fromJSON(json *notificationsJSON) *product.Notifications {
	return &product.Notifications{
		BackInStock:    json.BackInStock,
		OutOfStock:     json.OutOfStock,
		PriceBelows:    json.PriceBelows,
		PriceDecreases: json.PriceDecreases,
		PriceIncreases: json.PriceIncreases,
		PriceOver:      json.PriceOver,
	}
}
