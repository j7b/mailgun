// Package tracking implements domain tracking settings.
package tracking

import (
	"fmt"

	"github.com/j7b/mailgun/client"
)

// Type holds a boolean.
type Type struct {
	Active bool `json:"active"`
}

// Unsub includes footers.
type Unsub struct {
	Type
	HTMLFooter string `json:"html_footer"`
	TextFooter string `json:"text_footer"`
}

// Tracking is domain tracking parameters.
type Tracking struct {
	Click Type `json:"click"` //  "click": {
	//	"active": false
	//  },
	Open Type `json:"open"` //  "open": {
	//	"active": false
	//  },
	Unsubscribe Unsub `json:"unsubscribe"` //  "unsubscribe": {
	//	"active": false,
	//	"html_footer": "\n<br>\n<p><a href=\"%unsubscribe_url%\">unsubscribe</a></p>\n",
	//	"text_footer": "\n\nTo unsubscribe click: <%unsubscribe_url%>\n\n"
	//  }
}

// Settings retrieves tracking settings for API domain.
func Settings(c client.Caller) (*Tracking, error) {
	var t struct {
		Tracking *Tracking `json:"tracking,omitempty"`
	}
	return t.Tracking, c.Get(`/domains`, c.Domain(), `tracking`).
		Decode(&t)
}

func puts(c client.Caller, name string, active bool) error {
	return c.Put(`/domains`, c.Domain(), `tracking`, name).
		SetForm("active", fmt.Sprintf(`%v`, active)).
		Err()
}

// Clicks updates click tracking.
func Clicks(c client.Caller, active bool) error {
	return puts(c, "click", active)
}

// Opens updates open tracking.
func Opens(c client.Caller, active bool) error {
	return puts(c, "open", active)
}

// Unsubs updates unsubscribe tracking and footers.
func Unsubs(c client.Caller, u Unsub) error {
	req := c.Put(`/domains`, c.Domain(), `tracking`, `unsubscribe`).
		SetForm("active", fmt.Sprintf(`%v`, u.Active))
	if len(u.HTMLFooter) > 0 {
		req.SetForm("html_footer", u.HTMLFooter)
	}
	if len(u.TextFooter) > 0 {
		req.SetForm("text_footer", u.TextFooter)
	}
	return req.Err()
}
