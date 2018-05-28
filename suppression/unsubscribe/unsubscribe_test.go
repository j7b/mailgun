package unsubscribe

import (
	"testing"
	"time"

	"github.com/j7b/mailgun/client/mock"
)

func api(t *testing.T) *API {
	return &API{c: mock.Client(t)}
}

func TestUnsubscribes(t *testing.T) {
	res, err := api(t).List()
	if err != nil {
		t.Fatal(err)
	}
	want := 1
	if l := len(res.Unsubscribes); l != want {
		t.Fatal(`want`, want, `got`, l)
	}
}

func TestGet(t *testing.T) {
	a := `alice@example.com`
	res, err := api(t).Get(a)
	if err != nil {
		t.Fatal(err)
	}
	if res.Address != a {
		t.Fatal(`want`, a, `got`, res.Address)
	}
}

func TestUnsubscribe(t *testing.T) {
	s, now := `tst`, time.Now()
	if err := api(t).Unsubscribe(`a`, &s, &now); err != nil {
		t.Fatal(err)
	}
}

func TestUnsubscribeList(t *testing.T) {
	if err := api(t).UnsubscribeList([]Unsubscribe{}); err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	if err := api(t).Delete(`alice@example.com`); err != nil {
		t.Fatal(err)
	}
}
