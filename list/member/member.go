// Package member implements membership of mailing lists.
/*
The List type is intended for adding subscribers
to a mailing list. The Membership.Add method will add
list members in batches of 1000 and delete the
members from the list when they are added.
The intent is to have a mechanism to retry failed
calls or to queue subscriptions and execute in batches.

Note the mailgun documentation says something about
inaccurate counts after 100k members are added to a list,
and asserts the maximum  size of a list is 2.5 million.
*/
package member

import (
	"fmt"

	"github.com/j7b/mailgun/client/pager"
)

// Page of members. With paging.
type Page struct {
	Members []Member `json:"items"`
	pager.Pager
}

// Member of a list.
type Member struct {
	Address    string                 `json:"address"`
	Name       string                 `json:"name"`
	Subscribed *bool                  `json:"subscribed,omitempty"`
	Vars       map[string]interface{} `json:"vars"`
}

type list struct {
	m map[string]Member
}

func (l *list) Map() *map[string]Member {
	return &l.m
}

func (l *list) Add(address, name string, vars map[string]interface{}) error {
	if _, ok := l.m[address]; ok {
		return fmt.Errorf("Add: %s already in batch", address)
	}
	m := Member{Address: address, Name: name, Vars: vars}
	l.m[address] = m
	return nil
}

func (l *list) Members() []Member {
	if len(l.m) == 0 {
		return nil
	}
	s := make([]Member, 0, len(l.m))
	for _, v := range l.m {
		s = append(s, v)
	}
	return s
}

// List is the list interface.
type List interface {
	// Add adds subscriber info to the list.
	// If address already exists, error is returned.
	Add(address, name string, vars map[string]interface{}) error
	// Members is the set of members contained by this list.
	Members() []Member // Mostly to determine what wasn't done if error occurs.
}

// NewList returns a new List.
func NewList() List {
	return &list{m: make(map[string]Member)}
}
