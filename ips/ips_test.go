package ips

import (
	"testing"

	"github.com/j7b/mailgun/client/mock"
)

func api(t *testing.T) Caller {
	c := mock.Client(t)
	return API(c)
}

func TestAssigned(t *testing.T) {
	list, err := api(t).Assigned()
	if err != nil {
		t.Fatal(err)
	}
	if l := len(list); l != 2 {
		t.Fatal(`want 2 got`, l)
	}
}

func TestAssign(t *testing.T) {
	if err := api(t).Assign(`1.2.3.4`); err != nil {
		t.Fatal(err)
	}
}

func TestUnassign(t *testing.T) {
	if err := api(t).Unassign(`1.2.3.4`); err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	if info, err := api(t).Get(`1.2.3.4`); err != nil {
		t.Fatal(err)
	} else {
		if info.RDNS != "luna.mailgun.net" {
			t.Fatal(`want`, "luna.mailgun.net", `got`, info.RDNS)
		}
	}
}

func TestList(t *testing.T) {
	if list, err := api(t).List(false); err != nil {
		t.Fatal(err)
	} else {
		if l := len(list); l != 2 {
			t.Fatal(`want 2 got`, l)
		}
	}
}
