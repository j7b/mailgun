// Package mailgun consumes the mailgun API.
/*
The client has simple methods for common operations.
More sophisticated operations and types allowing fine-grained control
are found in subdirectories of this package.
*/
package mailgun

import (
	"fmt"
	"os"

	"github.com/j7b/mailgun/client"
	"github.com/j7b/mailgun/message"
)

var prefix = `https://api.mailgun.net/v3/`

// Client is a mailgun client.
type Client struct {
	client.Caller
}

// Send sends an HTML email, returning id.
func (c *Client) Send(from, subject, html string, to ...string) (id string, err error) {
	msg, err := message.New(from, subject, html, to...)
	if err != nil {
		return ``, err
	}
	o, err := msg.Send(c)
	if err != nil {
		return ``, err
	}
	return o.ID, nil
}

// New returns a Client for apikey and domain.
// If apikey is zero-length, attempts to use environment
// variable MAILGUN_KEY, failing that returns error.
// If domain is zero-length, attempts to use environment
// variable MAILGUN_DOMAIN, failing that returns error.
func New(apikey, domain string) (*Client, error) {
	if len(apikey) == 0 {
		apikey = os.Getenv("MAILGUN_KEY")
		if len(apikey) == 0 {
			return nil, fmt.Errorf("MAILGUN_KEY not set, apikey not supplied")
		}
	}
	if len(domain) == 0 {
		domain = os.Getenv("MAILGUN_DOMAIN")
		if len(domain) == 0 {
			return nil, fmt.Errorf("MAILGUN_DOMAIN not set, domain not supplied")
		}
	}
	return &Client{Caller: client.New(prefix, apikey, domain)}, nil
}
