// Package list implements the mailing list endpoint.
/*
The functions in this package are for operations on mailing lists.
The Membership type, created with the Manager function,
performs operations on members of mailing lists.

Documentation for this endpoint is at
https://documentation.mailgun.com/en/latest/api-mailinglists.html
*/
package list

import (
	"encoding/json"
	"fmt"

	"github.com/j7b/mailgun/client"
	"github.com/j7b/mailgun/client/pager"
	"github.com/j7b/mailgun/list/member"
)

/*
Docs say the maximum size of a mailing list
is 2.5 million. It doesn't say what happens
if you add a mailing list to a mailing list.
Probably nothing exciting.
*/

// AccessLevel is a list access level.
type AccessLevel interface {
	string() string
}

type accesslevel string

func (a accesslevel) string() string {
	return string(a)
}

// Access levels.
const (
	AccessReadOnly = accesslevel(`readonly`)
	AccessMembers  = accesslevel(`members`)  // members can post to list
	AccessEveryone = accesslevel(`everyone`) // everyone can post to list (don't)
)

// Page of mailing lists.
type Page struct {
	Lists []List `json:"items"`
	pager.Pager
}

func (p *Page) page(f func(interface{}) error) (*Page, error) {
	var page *Page
	return page, f(&page)
}

// Next page of mailing lists.
func (p *Page) Next() (*Page, error) {
	return p.page(p.Pager.Paging.Next)
}

// Previous page of mailing lists.
func (p *Page) Previous() (*Page, error) {
	return p.page(p.Pager.Paging.Previous)
}

// Last page of mailing lists.
func (p *Page) Last() (*Page, error) {
	return p.page(p.Pager.Paging.Last)
}

// First page of mailing lists.
func (p *Page) First() (*Page, error) {
	return p.page(p.Pager.Paging.First)
}

// List is a mailing list.
type List struct {
	AccessLevel string `json:"access_level"`  // : "everyone",
	Address     string `json:"address"`       // : "dev@samples.mailgun.org",
	CreatedAt   string `json:"created_at"`    // : "Tue, 06 Mar 2012 05:44:45 GMT",
	Description string `json:"description"`   // : "Mailgun developers list",
	Members     int    `json:"members_count"` // : 1,
	Name        string `json:"name"`          // : ""
}

// Lists returns a page of mailing lists.
func Lists(c client.Caller) (*Page, error) {
	var p *Page
	return p, c.Get(`/lists/pages`).Decode(&p)
}

// Address returns a mailing list by address.
func Address(c client.Caller, listaddress string) (*List, error) {
	var l struct {
		List *List `json:"list"`
	}
	return l.List, c.Get(`/lists`, listaddress).Decode(&l)
}

// New creates a new mailing list. Name and description may be
// zero-length, and access may be nil (default to readonly).
func New(c client.Caller, listaddress string, name string, description string, access AccessLevel) error {
	req := c.Post(`/lists`).SetForm(`address`, listaddress)
	if len(name) > 0 {
		req.SetForm(`name`, name)
	}
	if len(description) > 0 {
		req.SetForm(`description`, description)
	}
	if access != nil {
		req.SetForm(`access_level`, access.string())
	}
	return req.Err()
}

// Update updates an existing mailing list.
func Update(c client.Caller, listaddress string, newaddress *string, name *string, description *string, access AccessLevel) error {
	counter := 0
	req := c.Put(`/lists`, listaddress)
	if newaddress != nil {
		counter++
		req.SetForm(`address`, *newaddress)
	}
	if name != nil {
		counter++
		req.SetForm(`name`, *name)
	}
	if description != nil {
		counter++
		req.SetForm(`description`, *description)
	}
	if access != nil {
		counter++
		req.SetForm(`access_level`, access.string())
	}
	if counter == 0 {
		return fmt.Errorf("New: all parameters were nil")
	}
	return req.Err()
}

// Delete deletes a mailing list.
func Delete(c client.Caller, listaddress string) error {
	return c.Delete(`/lists`, listaddress).Err()
}

type caller client.Caller

// Membership allows operations on members of a list.
type Membership struct {
	caller
	address string
}

// Members returns a page of list members.
func (m *Membership) Members() (*member.Page, error) {
	var p *member.Page
	return p, m.Get(`/lists`, m.address, `members/pages`).Decode(&p)
}

// Member retrieves one list member by address.
func (m *Membership) Member(address string) (*member.Member, error) {
	var r struct {
		Member *member.Member `json:"member"`
	}
	return r.Member, m.Get(`/lists`, m.address, `members`, address).Decode(&r)
}

// Update performs an upsert operation on address. Nil vars
// will not update member vars. To delete vars, use
// an empty non-nil map (sets vars to `{}`).
func (m *Membership) Update(address string, newaddress *string, newname *string, vars map[string]interface{}) error {
	b, err := json.Marshal(vars)
	if err != nil {
		return err
	}
	req := m.Put(`/lists`, m.address, `members`)
	if newaddress != nil {
		req.SetForm(`address`, *newaddress)
	}
	if newaddress != nil {
		req.SetForm(`name`, *newname)
	}
	if s := string(b); s != `null` {
		req.SetForm(`vars`, s)
	}
	return req.Err()
}

// Subscribe subscribes address to list.
func (m *Membership) Subscribe(address string, name string, vars map[string]interface{}) error {
	b, err := json.Marshal(vars)
	if err != nil {
		return err
	}
	req := m.Post(`/lists`, m.address, `members`).SetForm(`address`, address).SetForm(`name`, name)
	if s := string(b); s != `null` {
		req.SetForm(`vars`, s)
	}
	return req.Err()
}

// Delete deletes address from list.
func (m *Membership) Delete(address string) error {
	return m.caller.Delete(`/lists`, m.address, `members`, address).Err()
}

type list interface {
	Map() *map[string]member.Member
}

func (m *Membership) add(mems []member.Member) error {
	if len(mems) == 0 {
		return nil
	}
	b, err := json.Marshal(mems)
	if err != nil {
		return err
	}
	return m.Post(`/lists`, m.address, `members.json`).SetForm(`members`, string(b)).SetForm(`upsert`, `yes`).Err()
}

// Add adds the members in l to the mailing list
// in batches of 1000. Performed as upsert but
// does not change subscription status.
func (m *Membership) Add(l member.List) error {
	p := l.(list).Map()
	list := *p
	mems := make([]member.Member, 0, 1000)
	for _, v := range list {
		mems = append(mems, v)
		if len(mems) == 1000 {
			if err := m.add(mems); err != nil {
				return err
			}
			for _, m := range mems {
				delete(list, m.Address)
			}
			mems = mems[:0]
		}
	}
	if err := m.add(mems); err != nil {
		return err
	}
	list = nil
	*p = list
	return nil
}

// Manager returns a member manager for listaddress.
func Manager(c client.Caller, listaddress string) *Membership {
	return &Membership{caller: c, address: listaddress}
}
