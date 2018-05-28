// Package webhook implements webhook API endpoints.
/*
Note docs say "Mailgun imposes a rate limit for the Webhook API endpoint.
Users may issue no more than 300 requests per minute, per account."

Callers with grave concerns about this should see bugs section.
*/
package webhook

import (
	"fmt"

	"github.com/j7b/mailgun/client"
)

// Hook is an interface for webhook names.
type Hook interface {
	h() string
}

type hook string

func (h hook) h() string {
	return string(h)
}

func (h hook) String() string {
	return string(h)
}

// Webhook names.
const (
	Bounce      = hook(`bounce`)
	Deliver     = hook(`deliver`)
	Drop        = hook(`drop`)
	Spam        = hook(`spam`)
	Unsubscribe = hook(`unsubscribe`)
	Click       = hook(`click`)
	Open        = hook(`open`)
)

// Webhook is a webhook URL.
type Webhook struct {
	URL string `json:"url"`
}

// Webhooks are the set
// of webhooks supported
// by a domain.
type Webhooks struct {
	Bounce      Webhook `json:"bounce"`
	Deliver     Webhook `json:"deliver"`
	Drop        Webhook `json:"drop"`
	Spam        Webhook `json:"spam"`
	Unsubscribe Webhook `json:"unsubscribe"`
	Click       Webhook `json:"click"`
	Open        Webhook `json:"open"`
}

// Get retrieves Webhooks for API domain.
func Get(c client.Caller) (*Webhooks, error) {
	var wh struct {
		Webhooks *Webhooks `json:"webhooks,omitempty"`
	}
	return wh.Webhooks, c.Get(`/domains`, c.Domain(), `webhooks`).Decode(&wh)
}

// URL returns the URL for webhook name.
func URL(c client.Caller, name Hook) (fmt.Stringer, error) {
	if name == nil {
		return nil, fmt.Errorf("URL: Hook must not be nil")
	}
	var wh struct {
		Webhook Webhook `json:"webhook"`
	}
	if err := c.Get(`/domains`, c.Domain(), `webhooks`, name.h()).Decode(&wh); err != nil {
		return nil, err
	}
	return hook(wh.Webhook.URL), nil
}

// New creates a new webhook for name.
func New(c client.Caller, name Hook, url string) error {
	if name == nil {
		return fmt.Errorf("New: Hook must not be nil")
	}
	return c.Post(`/domains`, c.Domain(), `webhooks`).SetForm("id", name.h()).SetForm(`url`, url).Err()
}

// Update webhook url.
func Update(c client.Caller, name Hook, url string) error {
	if name == nil {
		return fmt.Errorf("Update: Hook must not be nil")
	}
	return c.Put(`/domains`, c.Domain(), `webhooks`, name.h()).SetForm(`url`, url).Err()
}

// Delete webhook.
func Delete(c client.Caller, name Hook) error {
	if name == nil {
		return fmt.Errorf("Delete: Hook must not be nil")
	}
	return c.Delete(`/domains`, c.Domain(), `webhooks`, name.h()).Err()
}

// BUG(j7b): https://documentation.mailgun.com/en/latest/api-webhooks.html#webhooks
// is painfully incorrect.

// BUG(j7b): Callers should feel free to inspect retry-seconds themselves.
