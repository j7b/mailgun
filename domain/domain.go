// Package domain implements domain-related types and methods.
package domain

import (
	"fmt"

	"github.com/j7b/mailgun/client"
)

type caller client.Caller

// Caller calls domain methods.
type Caller struct {
	c caller
}

// API returns a Caller for c.
func API(c client.Caller) Caller {
	return Caller{c: c}
}

type spamAction string

// Spam action for inbound mail.
const (
	ActionDisabled = spamAction(`disabled`)
	ActionBlock    = spamAction(`block`)
	ActionTag      = spamAction(`tag`)
)

// MX is an MX record.
type MX struct {
	Priority   string `json:"priority"`    // "priority": "10",
	RecordType string `json:"record_type"` // "record_type": "MX",
	Valid      string `json:"valid"`       // "valid": "valid",
	Value      string `json:"value"`       // "value": "mxa.mailgun.org"
}

// RR is a DNS resource record.
type RR struct {
	Name       string `json:"name"`        // "name": "domain.com",
	RecordType string `json:"record_type"` // "record_type": "MX",
	Valid      string `json:"valid"`       // "valid": "valid",
	Value      string `json:"value"`       // "value": "mxa.mailgun.org"
}

// Domain is a mailgun domain configuration.
type Domain struct {
	CreatedAt    string `json:"created_at"`         // "created_at": "Wed, 10 Jul 2013 19:26:52 GMT",
	SMTPLogin    string `json:"smtp_login"`         // "smtp_login": "postmaster@samples.mailgun.org",
	Name         string `json:"name"`               // "name": "samples.mailgun.org",
	SMTPPassword string `json:"smtp_password"`      // "smtp_password": "4rtqo4p6rrx9",
	Wildcard     *bool  `json:"wildcard,omitempty"` // "wildcard": true,
	SpamAction   string `json:"spam_action"`        // "spam_action": "disabled",
	State        string `json:"state"`              // "state": "active"
}

// Info is information about a domain, including
// resource records.
type Info struct {
	Domain Domain `json:"domain"`
	MXList []MX   `json:"receiving_dns_records"`
	RRLIst []RR   `json:"sending_dns_records"`
}

// Delete removes a domain from your account.
func (c Caller) Delete(name string) error {
	return c.c.Delete(`/domains`, name).Err()
}

// New creates a new domain.
func (c Caller) New(name string, password string, action spamAction, wildcard bool) (*Info, error) {
	wc := "false"
	if wildcard {
		wc = "true"
	}
	req := c.c.Post(`/domains`)
	vals := req.Form()
	vals.Set("name", name)
	vals.Set("password", password)
	vals.Set("action", string(action))
	vals.Set("wildcard", wc)
	var info *Info
	return info, req.Decode(&info)
}

// Get returns Info about domain name.
func (c Caller) Get(name string) (*Info, error) {
	var info *Info
	return info, c.c.Get(`/domains`, name).Decode(&info)
}

// List returns a list of Domains. page must be
// greater than 0.
func (c Caller) List(page int) ([]Domain, error) {
	if page < 1 {
		return nil, fmt.Errorf("List: %v is < 1 ", page)
	}
	page--
	var output struct {
		Domains []Domain `json:"items"`
	}
	return output.Domains, c.c.Get(`/domains`).
		SetQuery("skip", fmt.Sprintf(`%v`, page*100)).
		Decode(&output)
}
