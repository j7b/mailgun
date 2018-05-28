// Package event implements event querying.
/*
Callers should first call the Queries function,
which returns a query client. Queries use
filter expression types to represent filter fields.
For example:

	Subject(`Hello Sailor`)

sets the subject filter of the query to
"Hello Sailor" (sans quotes).

Documentation for filter fields is at
https://documentation.mailgun.com/en/latest/api-events.html#event-polling
*/
package event

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/j7b/mailgun/client"
	"github.com/j7b/mailgun/client/pager"
	"github.com/j7b/mailgun/event/types"
)

// Events is a set of returned events.
type Events struct {
	list []types.Interface
	pager.Pager
}

func (e *Events) page(f func(interface{}) error) (*Events, error) {
	var ev *Events
	return ev, f(&ev)
}

// Next page of events.
func (e *Events) Next() (*Events, error) {
	return e.page(e.Pager.Paging.Next)
}

// Previous page of events.
func (e *Events) Previous() (*Events, error) {
	return e.page(e.Pager.Paging.Previous)
}

// Last page of events.
func (e *Events) Last() (*Events, error) {
	return e.page(e.Pager.Paging.Last)
}

// First page of events.
func (e *Events) First() (*Events, error) {
	return e.page(e.Pager.Paging.First)
}

// List returns the list of events. The underlying
// types are not pointers to "types" structs,
// but struct types.
func (e *Events) List() []types.Interface {
	return e.list
}

// FilterField holds a filter expression.
type FilterField interface {
	name() string
}

// Event is a filter expression.
type Event string

func (Event) name() string {
	return `event`
}

// List is a filter expression.
type List string

func (List) name() string {
	return `list`
}

// Attachment is a filter expression.
type Attachment string

func (Attachment) name() string {
	return `attachment`
}

// From is a filter expression.
type From string

func (From) name() string {
	return `from`
}

// MessageID is a filter expression.
type MessageID string

func (MessageID) name() string {
	return `message-id`
}

// Subject is a filter expression.
type Subject string

func (Subject) name() string {
	return `subject`
}

// To is a filter expression.
type To string

func (To) name() string {
	return `to`
}

// Size is a filter expression.
type Size string

func (Size) name() string {
	return `size`
}

// Recipient is a filter expression.
type Recipient string

func (Recipient) name() string {
	return `recipient`
}

// Tags is a filter expression.
type Tags string

func (Tags) name() string {
	return `tags`
}

// Severity is a filter expression.
type Severity string

func (Severity) name() string {
	return `severity`
}

// Client is a client for event queries.
type Client struct {
	c client.Caller
}

// Queries returns an event query client.
func Queries(c client.Caller) *Client {
	return &Client{c: c}
}

func decode(res *http.Response, i interface{}, err error) error {
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return json.NewDecoder(res.Body).Decode(i)
}

func parseresults(c client.Caller, res *http.Response, err error) (*Events, error) {
	var o struct {
		Events []json.RawMessage `json:"items"`
		pager.Pager
	}
	if err = decode(res, &o, err); err != nil {
		return nil, err
	}
	gens := make([]types.Generic, len(o.Events))
	for i := range o.Events {
		data := []byte(o.Events[i])
		gen := types.Generic{}
		if err := json.Unmarshal(data, &gen); err != nil {
			return nil, err
		}
		gens[i] = gen
	}
	events := Events{Pager: o.Pager}
	for i, g := range gens {
		data := []byte(o.Events[i])
		switch g.Event {
		case `accepted`:
			t := types.Accepted{}
			if err = json.Unmarshal(data, &t); err != nil {
				log.Println(string(data))
				return nil, err
			}
			events.list = append(events.list, t)
		case `rejected`: // not documented
			t := types.Failed{}
			if err = json.Unmarshal(data, &t); err != nil {
				log.Println(string(data))
				return nil, err
			}
			events.list = append(events.list, t)
		case `delivered`:
			t := types.Delivered{}
			if err = json.Unmarshal(data, &t); err != nil {
				log.Println(string(data))
				return nil, err
			}
			events.list = append(events.list, t)
		case `failed`:
			t := types.Failed{}
			if err = json.Unmarshal(data, &t); err != nil {
				log.Println(string(data))
				return nil, err
			}
			events.list = append(events.list, t)
		case `opened`:
			t := types.Opened{}
			if err = json.Unmarshal(data, &t); err != nil {
				log.Println(string(data))
				return nil, err
			}
			events.list = append(events.list, t)
		case `clicked`:
			t := types.Clicked{}
			if err = json.Unmarshal(data, &t); err != nil {
				log.Println(string(data))
				return nil, err
			}
			events.list = append(events.list, t)
		case `unsubscribed`:
			t := types.Unsubscribed{}
			if err = json.Unmarshal(data, &t); err != nil {
				log.Println(string(data))
				return nil, err
			}
			events.list = append(events.list, t)
		case `complained`:
			t := types.Complained{}
			if err = json.Unmarshal(data, &t); err != nil {
				log.Println(string(data))
				return nil, err
			}
			events.list = append(events.list, t)
		case `stored`:
			t := types.Stored{}
			if err = json.Unmarshal(data, &t); err != nil {
				log.Println(string(data))
				return nil, err
			}
			events.list = append(events.list, t)
		default:
			return nil, fmt.Errorf("Query: unexpected event name %s", g.Event)
		}
	}
	return &events, nil
}

// Query executes a query using the parameters provided.
func (c *Client) Query(begin *time.Time, end *time.Time, ascending *bool, filters ...FilterField) (*Events, error) {
	req := c.c.Get(`events`).SetQuery("limit", "300")
	if begin != nil {
		req.SetQuery("begin", begin.Format(time.RFC1123))
	}
	if end != nil {
		req.SetQuery("end", end.Format(time.RFC1123))
	}
	if ascending != nil {
		req.SetQuery("ascending", fmt.Sprintf(`%v`, *ascending))
	}
	for _, f := range filters {
		req.AddQuery(f.name(), fmt.Sprintf(`%s`, f))
	}
	res, err := req.Do()
	return parseresults(c.c, res, err)
}
