package bounce

import (
	"testing"
	"time"

	"github.com/j7b/mailgun/client/mock"
)

func api(t *testing.T) *API {
	return Bounces(mock.Client(t))
}

func TestList(t *testing.T) {
	lst, err := api(t).List()
	if err != nil {
		t.Fatal(err)
	}
	w := 1
	if l := len(lst.Bounces); l != w {
		t.Fatal(`want`, w, `got`, l)
	}
}

func TestGet(t *testing.T) {
	a := `foo@bar.com`
	b, err := api(t).Get(a)
	if err != nil {
		t.Fatal(err)
	}
	if b.Address != a {
		t.Fatal(`want`, a, `got`, b.Address)
	}
}

func TestAdd(t *testing.T) {
	c, e, now := 550, `bounce`, time.Now()
	if err := api(t).Add(`a`, &c, &e, &now); err != nil {
		t.Fatal(err)
	}
}

func TestAddList(t *testing.T) {
	if err := api(t).AddList([]Bounce{}); err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	if err := api(t).Delete(`foo@bar.com`); err != nil {
		t.Fatal(err)
	}
}

func TestClear(t *testing.T) {
	if err := api(t).Clear(); err != nil {
		t.Fatal(err)
	}
}
