package complaint

import (
	"testing"
	"time"

	"github.com/j7b/mailgun/client/mock"
)

func api(t *testing.T) *API {
	return Complaints(mock.Client(t))
}

func TestList(t *testing.T) {
	lst, err := api(t).List()
	if err != nil {
		t.Fatal(err)
	}
	if len(lst.Complaints) != 1 {
		t.Fatal(`want 1 got`, len(lst.Complaints))
	}
}

func TestGet(t *testing.T) {
	a := `baz@example.com`
	res, err := api(t).Get(a)
	if err != nil {
		t.Fatal(err)
	}
	if a != res.Address {
		t.Fatal(`want`, a, `got`, res.Address)
	}
}

func TestAdd(t *testing.T) {
	now := time.Now()
	if err := api(t).Add(`a`, &now); err != nil {
		t.Fatal(err)
	}
}

func TestAddList(t *testing.T) {
	if err := api(t).AddList([]Complaint{}); err != nil {
		t.Fatal(err)
	}
}

func TestDelete(t *testing.T) {
	if err := api(t).Delete(`a`); err != nil {
		t.Fatal(err)
	}
}
