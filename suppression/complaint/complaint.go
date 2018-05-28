// Package complaint implements complaint endpoints.
package complaint

import (
	"bytes"
	"encoding/json"
	"time"

	"github.com/j7b/mailgun/client"
	"github.com/j7b/mailgun/client/pager"
)

// Complaint is information associated with a complaint.
type Complaint struct {
	Address   string `json:"address"`              // "address": "alice@example.com",
	CreatedAt string `json:"created_at,omitempty"` // "created_at": "Fri, 21 Oct 2011 11:02:55 GMT"
}

// List is a list of Complaints.
type List struct {
	Complaints []Complaint `json:"items"`
	pager.Pager
}

func (*List) pager(f func(interface{}) error) (*List, error) {
	var c *List
	return c, f(&c)
}

// Next returns the next page.
func (c *List) Next() (*List, error) {
	return c.pager(c.Pager.Paging.Next)
}

// Previous returns the previous page.
func (c *List) Previous() (*List, error) {
	return c.pager(c.Pager.Paging.Previous)
}

// First returns the first page.
func (c *List) First() (*List, error) {
	return c.pager(c.Pager.Paging.First)
}

// Last returns the last page.
func (c *List) Last() (*List, error) {
	return c.pager(c.Pager.Paging.Last)
}

// API to complaints.
type API struct {
	c client.Caller
}

// Complaints returns the endpoint interface.
func Complaints(c client.Caller) *API {
	return &API{c: c}
}

// List complaints.
func (a *API) List() (*List, error) {
	var l *List
	return l, a.c.Get(`complaints`).Decode(&l)
}

// Get single complaint.
func (a *API) Get(address string) (*Complaint, error) {
	var c *Complaint
	return c, a.c.Get(`complaints`, address).Decode(&c)
}

// Add a complaint.
func (a *API) Add(address string, timestamp *time.Time) error {
	req := a.c.Post(`complaints`).SetForm(`address`, address)
	if timestamp != nil {
		req.SetForm(`created_at`, timestamp.Format(time.RFC1123))
	}
	return req.Err()
}

// AddList adds list of complaints.
func (a *API) AddList(list []Complaint) error {
	b, err := json.Marshal(list)
	if err != nil {
		return err
	}
	return a.c.Post(`complaints`).SetHeader(`Content-Type`, `application/json`).Payload(bytes.NewReader(b)).Err()
}

// Delete a single complaint.
func (a *API) Delete(address string) error {
	return a.c.Delete(`complaints`, address).Err()
}
