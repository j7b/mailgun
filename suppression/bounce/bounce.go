// Package bounce implements bounce endpoints.
package bounce

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"

	"github.com/j7b/mailgun/client"
	"github.com/j7b/mailgun/client/pager"
)

// Bounce is the information assciated with a bounce.
type Bounce struct {
	Address   string `json:"address,omitempty"`    // "address": "alice@example.com",
	Code      string `json:"code,omitempty"`       // "code": "550",
	Error     string `json:"error,omitempty"`      // "error": "No such mailbox",
	CreatedAt string `json:"created_at,omitempty"` // "created_at": "Fri, 21 Oct 2011 11:02:55 GMT"
}

// List is a list of Bounce with paging.
type List struct {
	Bounces []Bounce `json:"items"`
	pager.Pager
}

func (b *List) pager(f func(interface{}) error) (*List, error) {
	var bl *List
	return bl, f(&bl)
}

// Next returns the next page.
func (b *List) Next() (*List, error) {
	return b.pager(b.Pager.Paging.Next)
}

// Previous returns the previous page.
func (b *List) Previous() (*List, error) {
	return b.pager(b.Pager.Paging.Previous)
}

// First returns the first page.
func (b *List) First() (*List, error) {
	return b.pager(b.Pager.Paging.First)
}

// Last returns the last page.
func (b *List) Last() (*List, error) {
	return b.pager(b.Pager.Paging.Last)
}

// API implements the bounce API.
type API struct {
	c client.Caller
}

// List returns a list of bounces.
func (a *API) List() (*List, error) {
	var bl *List
	return bl, a.c.Get(`bounces`).Decode(&bl)
}

// Bounces returns a list of bounces.
func Bounces(c client.Caller) *API {
	return &API{c: c}
}

// Get retrieves a bounce by address.
func (a *API) Get(address string) (*Bounce, error) {
	var b *Bounce
	return b, a.c.Get(`bounces`, address).Decode(&b)
}

// Add a bounce.
func (a *API) Add(address string, code *int, err *string, timestamp *time.Time) error {
	req := a.c.Post(`bounces`).SetForm(`address`, address)
	if code != nil {
		req.SetForm(`code`, fmt.Sprintf(`%v`, *code))
	}
	if err != nil {
		req.SetForm(`error`, *err)
	}
	if timestamp != nil {
		req.SetForm(`created_at`, timestamp.Format(time.RFC1123))
	}
	return req.Err()
}

// AddList adds a list of bounces.
func (a *API) AddList(bounces []Bounce) error {
	payload, err := json.Marshal(bounces)
	if err != nil {
		return err
	}
	return a.c.Post(`bounces`).SetHeader(`Content-Type`, `application/json`).Payload(bytes.NewReader(payload)).Err()
}

// Delete a single bounce
func (a *API) Delete(address string) error {
	return a.c.Delete(`bounces`, address).Err()
}

// Clear clears the bounce list.
func (a *API) Clear() error {
	return a.c.Delete(`bounces`).Err()
}
