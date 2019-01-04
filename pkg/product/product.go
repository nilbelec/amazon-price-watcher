package product

import "time"

// Product model
type Product struct {
	URL           string
	Title         string
	ImageURL      string
	LastPrice     float32
	Price         float32
	Currency      string
	Added         time.Time
	LastUpdate    time.Time
	changed       bool
	Notifications Notifications
}

// Notifications representation of a Product notifications
type Notifications struct {
	PriceBelows    float32
	PriceOver      float32
	PriceIncreases bool
	PriceDecreases bool
	OutOfStock     bool
	BackInStock    bool
}

// UpdateInfo updates the product values
func (p *Product) UpdateInfo(other Product) {
	p.changed = false
	p.Title = other.Title
	p.ImageURL = other.ImageURL
	if p.Price != other.Price {
		p.LastPrice = p.Price
		p.Price = other.Price
		p.changed = true
	}
	p.Currency = other.Currency
	p.LastUpdate = time.Now()
}

// GetTotal returns the number of enabled notifications
func (n *Notifications) GetTotal() int {
	total := 0
	if n.PriceBelows > 0 {
		total++
	}
	if n.PriceOver > 0 {
		total++
	}
	if n.PriceIncreases {
		total++
	}
	if n.PriceDecreases {
		total++
	}
	if n.OutOfStock {
		total++
	}
	if n.BackInStock {
		total++
	}
	return total
}

// ShouldSendAnyNotification returns true if the product should trigger any notifications
func (p *Product) ShouldSendAnyNotification() bool {
	return p.ShouldSendPriceDecreasesNotification() ||
		p.ShouldSendPriceIncreasesNotification() ||
		p.ShouldSendPriceBelowsNotification() ||
		p.ShouldSendPriceOverNotification() ||
		p.ShouldSendBackInStockNotification() ||
		p.ShouldSendOutOfStockNotification()
}

// ShouldSendPriceDecreasesNotification returns true if should trigger a "price decreases" notification
func (p *Product) ShouldSendPriceDecreasesNotification() bool {
	return p.Notifications.PriceDecreases && p.changed && p.LastPrice > p.Price
}

// ShouldSendPriceIncreasesNotification returns true if should trigger a "price increases" notification
func (p *Product) ShouldSendPriceIncreasesNotification() bool {
	return p.Notifications.PriceIncreases && p.changed && p.LastPrice < p.Price
}

// ShouldSendPriceBelowsNotification returns true if should trigger a "price belows" notification
func (p *Product) ShouldSendPriceBelowsNotification() bool {
	return p.Notifications.PriceBelows > 0 && p.changed && p.Price <= p.Notifications.PriceBelows
}

// ShouldSendPriceOverNotification returns true if should trigger a "price over" notification
func (p *Product) ShouldSendPriceOverNotification() bool {
	return p.Notifications.PriceOver > 0 && p.changed && p.Price >= p.Notifications.PriceOver
}

// ShouldSendBackInStockNotification returns true if should trigger a "back in stock" notification
func (p *Product) ShouldSendBackInStockNotification() bool {
	return p.Notifications.BackInStock && p.changed && p.LastPrice == 0 && p.Price > 0
}

// ShouldSendOutOfStockNotification returns true if should trigger a "out of stock" notification
func (p *Product) ShouldSendOutOfStockNotification() bool {
	return p.Notifications.OutOfStock && p.changed && p.LastPrice > 0 && p.Price == 0
}
