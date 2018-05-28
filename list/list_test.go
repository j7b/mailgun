package list

import (
	"testing"

	"github.com/j7b/mailgun/list/member"

	"github.com/j7b/mailgun/client/mock"
)

func TestLists(t *testing.T) {
	c := mock.Client(t)
	page, err := Lists(c)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(page.Lists); l != 2 {
		t.Fatal(`want 2 got`, l)
	}
}

func TestAddress(t *testing.T) {
	c := mock.Client(t)
	list, err := Address(c, `list@domain.mock`)
	if err != nil {
		t.Fatal(err)
	}
	if list.Members != 3 {
		t.Fatal(`want 3 got`, list.Members)
	}
}

func TestNew(t *testing.T) {
	c := mock.Client(t)
	err := New(c, `herp@derp.com`, `Herp`, `The Derp`, AccessReadOnly)
	if err != nil {
		t.Fatal(err)
	}
}

func TestUpdate(t *testing.T) {
	c := mock.Client(t)
	s := `S`
	err := Update(c, `list@domain.mock`, &s, &s, &s, AccessMembers)
	if err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	c := mock.Client(t)
	err := Delete(c, `list@domain.mock`)
	if err != nil {
		t.Fatal(err)
	}
}

func mgr(t *testing.T) *Membership {
	return Manager(mock.Client(t), `list@domain.mock`)
}

func TestMembers(t *testing.T) {
	c := mgr(t)
	page, err := c.Members()
	if err != nil {
		t.Fatal(err)
	}
	if l := len(page.Members); l != 1 {
		t.Fatal(`want 1 got`, l)
	}
}

func TestMember(t *testing.T) {
	c := mgr(t)
	addr := `bar@example.com`
	mbr, err := c.Member(addr)
	if err != nil {
		t.Fatal(err)
	}
	if mbr.Address != addr {
		t.Fatal(`want`, addr, `got`, mbr.Address)
	}
}

func TestMemberUpdate(t *testing.T) {
	c := mgr(t)
	s := `s`
	if err := c.Update(`bar@example.com`, &s, &s, nil); err != nil {
		t.Fatal(err)
	}
}

func TestSubscribe(t *testing.T) {
	c := mgr(t)
	if err := c.Subscribe(`bar@example.com`, `name`, map[string]interface{}{"a": "b"}); err != nil {
		t.Fatal(err)
	}
}

func TestMemberDelete(t *testing.T) {
	c := mgr(t)
	if err := c.Delete(`bar@example.com`); err != nil {
		t.Fatal(err)
	}
}

func TestAdd(t *testing.T) {
	c := mgr(t)
	l := member.NewList()
	l.Add(`a`, `b`, nil)
	if err := c.Add(l); err != nil {
		t.Fatal(err)
	}
}
