package smtp

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/smtp"
	"strings"

	"github.com/gobuffalo/packr"
	"github.com/nilbelec/amazon-price-watcher/pkg/product"
)

const (
	mime = "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
)

// Notifier is a SMTP Notifier
type Notifier struct {
	config *Configuration
	b      packr.Box
}

// Configuration handles the required parameters to use the client
type Configuration struct {
	Host     func() string
	Port     func() int
	Username func() string
	Password func() string
	To       func() []string
}

// NewNotifier creates a new SMTP notifier
func NewNotifier(c *Configuration) *Notifier {
	b := packr.NewBox(".")
	return &Notifier{c, b}
}

// IsConfigured returns true if the notifier can be used
func (n *Notifier) IsConfigured() bool {
	return n.config.Host() != "" &&
		n.config.Port() > 0 &&
		n.config.Username() != "" &&
		n.config.Password() != "" &&
		n.config.To() != nil &&
		len(n.config.To()) > 0
}

// NotifyChanges sends a email message notifying a product change
func (n *Notifier) NotifyChanges(p *product.Product) {
	if !n.IsConfigured() || !p.ShouldSendAnyNotification() {
		return
	}
	msg, err := n.prepareMessage(p)
	if err != nil {
		log.Printf("Error preparing notification mail: %s\n", err.Error())
		return
	}
	err = n.sendMail(msg)
	if err != nil {
		log.Printf("Error sending notification mail: %s\n", err.Error())
	}
}

func (n *Notifier) prepareMessage(p *product.Product) (b []byte, err error) {
	subject := "[amazon-price-watcher] - There is a change in one of your products!"
	body, err := n.prepareBody(p)
	if err != nil {
		return
	}
	msg := "To: " + strings.Join(n.config.To(), ", ") + "\r\nSubject: " + subject + "\r\n" + mime + "\r\n" + body
	b = []byte(msg)
	return
}

func (n *Notifier) prepareBody(p *product.Product) (body string, err error) {
	b, err := n.b.Find("mailbody.html")
	if err != nil {
		return
	}
	t := template.New("body")
	t.Parse(string(b))
	buf := new(bytes.Buffer)
	t.Execute(buf, nil)
	bs, err := ioutil.ReadAll(buf)
	body = string(bs)
	return
}

func (n *Notifier) sendMail(msg []byte) error {
	auth := smtp.PlainAuth("", n.config.Username(), n.config.Password(), n.config.Host())
	addr := fmt.Sprintf("%s:%d", n.config.Host(), n.config.Port())
	return smtp.SendMail(addr, auth, n.config.Username(), n.config.To(), msg)
}
