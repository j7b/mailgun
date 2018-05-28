// Package pager implements the mailgun paging convention.
/*
Not intended for direct use by user code.
*/
package pager

import (
	"fmt"
	"io"

	"github.com/j7b/mailgun/client"
)

// Interface allows caller to be set.
type Interface interface {
	SetCaller(client.Caller)
}

var _ = Interface(Pager{})

// Pager is a container for Paging.
type Pager struct {
	Paging *Paging `json:"paging"`
}

// SetCaller implements Interface.
func (p Pager) SetCaller(c client.Caller) {
	if p.Paging != nil {
		p.Paging.c = c
	}
}

// Paging contains paging info.
type Paging struct {
	c client.Caller
	N string `json:"next"` // "next":
	// "https://url_to_next_page",
	P string `json:"previous"` // "previous":
	// "https://url_to_previous_page",
	F string `json:"first"` // "first":
	// "https://url_to_first_page",
	L string `json:"last"` // "last":
	// "https://url_to_last_page"
}

func (p *Paging) get(uri string, i interface{}) error {
	if p.c == nil {
		return fmt.Errorf("paging: caller not set")
	}
	return p.c.Get(uri).Decode(i)
}

// Next decodes the next url to i. Returns
// io.EOF if p is nil or N is empty.
func (p *Paging) Next(i interface{}) error {
	if p == nil || len(p.N) == 0 {
		return io.EOF
	}
	return p.get(p.N, i)
}

// Previous decodes the previous url to i. Returns
// io.EOF if p is nil or P is empty.
func (p *Paging) Previous(i interface{}) error {
	if p == nil || len(p.P) == 0 {
		return io.EOF
	}
	return p.get(p.P, i)
}

// First decodes the first url to i. Returns
// io.EOF if p is nil or F is empty.
func (p *Paging) First(i interface{}) error {
	if p == nil || len(p.F) == 0 {
		return io.EOF
	}
	return p.get(p.F, i)
}

// Last decodes the last url to i. Returns
// io.EOF if p is nil or L is empty.
func (p *Paging) Last(i interface{}) error {
	if p == nil || len(p.L) == 0 {
		return io.EOF
	}
	return p.get(p.L, i)
}

// BUG(j7b): Paging might be inconsistently expressed by the API, and
// it's awful hard to test without 2 pages worth of data.
