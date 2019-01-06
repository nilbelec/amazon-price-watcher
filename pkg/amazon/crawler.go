package amazon

import (
	"crypto/tls"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/nilbelec/amazon-price-watcher/pkg/product"
	"golang.org/x/net/html"
)

// Crawler is an Amazon crawler
type Crawler struct {
}

// NewCrawler creates a new Amazon product crawler
func NewCrawler() *Crawler {
	return &Crawler{}
}

// Extract extracts Product information from an Amazon product URL
func (c *Crawler) Extract(inputURL string) (p *product.Product, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
		}
	}()
	_, err = url.ParseRequestURI(inputURL)
	if err != nil {
		err = errors.New("That doesn't seem like a valid URL")
		return
	}
	doc, err := loadURL(inputURL)
	if err != nil {
		err = errors.New("Cannot extract the product information. Are you sure it is a valid Amazon product URL?")
		return
	}
	url, err := findURL(doc)
	if err != nil {
		err = errors.New("Cannot extract the product information. Are you sure it is a valid Amazon product URL?")
		return
	}
	title, err := findTitle(doc)
	if err != nil {
		err = errors.New("Cannot extract the product information. Are you sure it is a valid Amazon product URL?")
		return
	}
	imageURL, err := findImageURL(doc)
	if err != nil {
		err = errors.New("Cannot extract the product information. Are you sure it is a valid Amazon product URL?")
		return
	}
	price, currency, _ := findPriceAndCurrency(doc)
	if err != nil {
		price = 0.0
		currency = "-"
	}
	p = &product.Product{
		URL:        url,
		Title:      title,
		ImageURL:   imageURL,
		Price:      float32(price),
		Currency:   currency,
		LastUpdate: time.Now()}
	return
}

func loadURL(url string) (*html.Node, error) {
	resp, err := doRequest(url)
	if err != nil {
		return nil, errors.New("Error while requesting the page: " + err.Error())
	}
	return parseResponse(resp)
}

func parseResponse(resp *http.Response) (*html.Node, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	textBody := string(body)
	return html.Parse(strings.NewReader(textBody))
}

func doRequest(url string) (*http.Response, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}
	transport := &http.Transport{
		TLSClientConfig: tlsConfig,
	}
	client := http.Client{Transport: transport}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)")
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func findURL(doc *html.Node) (url string, err error) {
	node := htmlquery.FindOne(doc, "/html/head/link[@rel=\"canonical\"]")
	if node == nil {
		err = errors.New("URL element not found")
		return
	}
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
		err = errors.New("Image src not found")
	}
	return
}
