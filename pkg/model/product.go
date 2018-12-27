package model

import "time"

// Product model
type Product struct {
	URL        string
	Title      string
	ImageURL   string
	LastPrice  float32
	Price      float32
	Currency   string
	Added      time.Time
	LastUpdate time.Time
	changed    bool
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

// Changed returns true if there is a price change in the product
func (p *Product) Changed() bool {
	return p.changed
}

// PriceHasDecreased returns true if the LastPrice
// is greater than the current Price
func (p *Product) PriceHasDecreased() bool {
	return p.changed && p.LastPrice > p.Price
}

// PriceHasIncrecreased returns true if the LastPrice
// is smaller than the current Price
func (p *Product) PriceHasIncrecreased() bool {
	return p.changed && p.LastPrice < p.Price
}
