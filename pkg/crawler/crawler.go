package crawler

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/nilbelec/amazon-price-watcher/pkg/model"
	"golang.org/x/net/html"
)

// ProductCrawler Amazon crawler
type ProductCrawler struct {
}

// NewProductCrawler creates a new ProductCrawler
func NewProductCrawler() *ProductCrawler {
	return &ProductCrawler{}
}

// FindByURL finds Product by Amazon URL
func (pc *ProductCrawler) FindByURL(inputURL string) (product model.Product, err error) {
	for i := 0; i < 3; i++ {
		product, err = pc.findByURL(inputURL)
		if err != nil {
			continue
		}
		return
	}
	return
}

func (pc *ProductCrawler) findByURL(inputURL string) (product model.Product, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	doc, err := htmlquery.LoadURL(inputURL)
	if err != nil {
		return
	}
	url, err := findURL(doc)
	if err != nil {
		return
	}
	title, err := findTitle(doc)
	if err != nil {
		return
	}
	imageURL, err := findImageURL(doc)
	if err != nil {
		return
	}
	price, currency, _ := findPriceAndCurrency(doc)
	if err != nil {
		price = 0.0
		currency = "-"
	}
	product = model.Product{
		URL:        url,
		Title:      title,
		ImageURL:   imageURL,
		Price:      float32(price),
		Currency:   currency,
		LastUpdate: time.Now()}
	return
}

func findURL(doc *html.Node) (url string, err error) {
	node := htmlquery.FindOne(doc, "/html/head/link[@rel=\"canonical\"]")
	url = strings.TrimSpace(htmlquery.SelectAttr(node, "href"))
	if url == "" {
		err = errors.New("URL href attr not found")
	}
	return
}

func findPriceAndCurrency(doc *html.Node) (price float64, currency string, err error) {
	span := htmlquery.FindOne(doc, "//span[@id='priceblock_ourprice']")
	if span == nil {
		span = htmlquery.FindOne(doc, "//span[@id='priceblock_dealprice']")
	}
	if span == nil {
		err = errors.New("Price element not found")
		return
	}
	text := htmlquery.InnerText(span)
	if text == "" {
		err = errors.New("Price and Currency text not found")
		return
	}
	text = strings.TrimSpace(text)
	var priceString string
	if text[0] == '$' {
		currency = string(text[0])
		priceString = text[1:len(text)]
	} else {
		textSplit := strings.Split(text, " ")
		currency = textSplit[0]
		priceString = textSplit[1]
	}
	priceString = strings.Replace(priceString, ",", ".", -1)
	price, err = strconv.ParseFloat(priceString, 32)
	return
}

func findTitle(doc *html.Node) (title string, err error) {
	titleEl := htmlquery.FindOne(doc, "//span[@id='productTitle']")
	if titleEl == nil {
		err = errors.New("Title element not found")
		return
	}
	title = htmlquery.InnerText(titleEl)
	if title == "" {
		err = errors.New("Title text not found")
		return
	}
	title = strings.TrimSpace(title)
	return
}

func findImageURL(doc *html.Node) (imageURL string, err error) {
	imgEl := htmlquery.FindOne(doc, "//img[@id='landingImage']")
	if imgEl == nil {
		err = errors.New("ImageURL element not found")
		return
	}
	imageURL = htmlquery.SelectAttr(imgEl, "src")
	if imageURL == "" {
		err = errors.New("Price text not found")
	}
	return
}
