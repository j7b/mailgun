// Package handler implements types for webhook consumption.
/*
The Dispatcher type is a simple http.Handler for receiving webhooks.
Care should be taken to not block in handler events (lengthy operations
should use goroutines or channels, see examples).
*/
package handler

import (
	"fmt"
	"io"
	"log"
	"mime"
	"net/http"

	"github.com/j7b/mailgun/webhook/types"
)

func ctype(h http.Header) string {
	ctype := h.Get(`Content-Type`)
	if len(ctype) == 0 {
		return ctype
	}
	d, params, err := mime.ParseMediaType(ctype)
	if err != nil || d != `multipart/form-data` {
		return ``
	}
	return params["boundary"]
}

// APIKey is an interface to
// API key.
type APIKey interface {
	Key() string
}

// Key is an APIKey.
type Key string

// Key implements APIKey.
func (k Key) Key() string {
	return string(k)
}

// Dispatcher is a utility type to dispatch
// webhook events to handler functions. If Logger is not nil
// errors will be logged to it when
// a Dispatcher is used as an http.Handler.
// If Key is not nil, the signature
// of webhooks will be validated
// and invalid payloads will not be
// dispatched.
type Dispatcher struct {
	Key         APIKey
	Logger      *log.Logger
	Bounce      func(*types.Bounce) error
	Click       func(*types.Click) error
	Complaint   func(*types.Complaint) error
	Delivered   func(*types.Delivered) error
	Drop        func(*types.Drop) error
	Open        func(*types.Open) error
	Unsubscribe func(*types.Unsubscribe) error
}

func (d *Dispatcher) bounce(b *types.Bounce) error {
	if d.Bounce != nil {
		return d.Bounce(b)
	}
	return nil
}
func (d *Dispatcher) click(c *types.Click) error {
	if d.Click != nil {
		return d.Click(c)
	}
	return nil
}
func (d *Dispatcher) complaint(c *types.Complaint) error {
	if d.Complaint != nil {
		return d.Complaint(c)
	}
	return nil
}
func (d *Dispatcher) delivered(v *types.Delivered) error {
	if d.Delivered != nil {
		return d.Delivered(v)
	}
	return nil
}
func (d *Dispatcher) drop(p *types.Drop) error {
	if d.Drop != nil {
		return d.Drop(p)
	}
	return nil
}
func (d *Dispatcher) open(o *types.Open) error {
	if d.Open != nil {
		return d.Open(o)
	}
	return nil
}
func (d *Dispatcher) unsubscribe(u *types.Unsubscribe) error {
	if d.Unsubscribe != nil {
		return d.Unsubscribe(u)
	}
	return nil
}

// ServeHTTP implements http.Handler. If a handler function
// returns an error, the error will be returned to
// the caller with status code 429.
func (d *Dispatcher) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := d.HTTPRequest(r)
	if err != nil {
		if d.Logger != nil {
			d.Logger.Println(err)
		}
		http.Error(w, err.Error(), 429)
	}
}

// HTTPRequest dispatches an http.Request.
func (d *Dispatcher) HTTPRequest(r *http.Request) error {
	ev, err := ParseRequest(r)
	if err != nil {
		return err
	}
	return d.Dispatch(ev)
}

// Reader dispatches the content of r.
func (d *Dispatcher) Reader(r io.Reader) error {
	ev, err := types.Decode(r, ``)
	if err != nil {
		return err
	}
	return d.Dispatch(ev)
}

// Dispatch dispatches ev.
func (d *Dispatcher) Dispatch(ev types.Type) error {
	if d.Key != nil {
		if err := types.Validate(ev, d.Key); err != nil {
			return fmt.Errorf(`Dispatch: validation failed`)
		}
	}
	switch t := ev.(type) {
	case *types.Bounce:
		return d.bounce(t)
	case *types.Click:
		return d.click(t)
	case *types.Complaint:
		return d.complaint(t)
	case *types.Delivered:
		return d.delivered(t)
	case *types.Drop:
		return d.drop(t)
	case *types.Open:
		return d.open(t)
	case *types.Unsubscribe:
		return d.unsubscribe(t)
	}
	return fmt.Errorf("unexpected event type %T", ev)
}

// ParseRequest decodes a webhook Type
// from r.
func ParseRequest(r *http.Request) (types.Type, error) {
	return types.Decode(r.Body, ctype(r.Header))
}
