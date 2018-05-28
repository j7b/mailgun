// Package unsubscribe implements unsubscribe endpoints.
package unsubscribe

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/j7b/mailgun/client"
	"github.com/j7b/mailgun/client/pager"
)

// Unsubscribe is information associated with an unsubscribe.
type Unsubscribe struct {
	Address   string `json:"address"`              // "address": "alice@example.com",
	Tag       string `json:"tag,omitempty"`        // "tag": "*",
	CreatedAt string `json:"created_at,omitemtpy"` // "created_at": "Fri, 21 Oct 2011 11:02:55 GMT"
}

// List is a list of Unsubscribe with paging.
type List struct {
	Unsubscribes []Unsubscribe `json:"items"`
	pager.Pager
}

func (u *List) pager(f func(interface{}) error) (*List, error) {
	var ul *List
	return ul, f(&ul)
}

// Next returns the next page.
func (u *List) Next() (*List, error) {
	return u.pager(u.Pager.Paging.Next)
}

// Previous returns the previous page.
func (u *List) Previous() (*List, error) {
	return u.pager(u.Pager.Paging.Previous)
}

// First returns the first page.
func (u *List) First() (*List, error) {
	return u.pager(u.Pager.Paging.First)
}

// Last returns the last page.
func (u *List) Last() (*List, error) {
	return u.pager(u.Pager.Paging.Last)
}

// API implements the unsubscribe API.
type API struct {
	c client.Caller
}

// Unsubscribes returns the unsubscribe API.
func Unsubscribes(c client.Caller) *API {
	return &API{c: c}
}

// List returns a list of unsubscribes.
func (a *API) List() (*List, error) {
	var ul *List
	return ul, a.c.Get(`unsubscribes`).Decode(&ul)
}

// Get a single unsubscribe.
func (a *API) Get(address string) (*Unsubscribe, error) {
	var u *Unsubscribe
	return u, a.c.Get(`unsubscribes`, address).Decode(&u)
}

// Unsubscribe address.
func (a *API) Unsubscribe(address string, tag *string, timestamp *time.Time) error {
	req := a.c.Post(`unsubscribes`).SetForm(`address`, address)
	if tag != nil {
		req.SetForm(`tag`, *tag)
	}
	if timestamp != nil {
		req.SetForm(`created_at`, timestamp.Format(time.RFC1123))
	}
	return req.Err()
}

// UnsubscribeList unsubscribes list.
func (a *API) UnsubscribeList(list []Unsubscribe) error {
	b, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return a.c.Post(`unsubscribes`).SetHeader(`Content-Type`, `application/json`).Payload(bytes.NewReader(b)).Err()
}

// Delete an address.
func (a *API) Delete(address string) error {
	return a.c.Delete(`unsubscribes`, address).Err()
}
