package tags

import (
	"testing"

	"github.com/j7b/mailgun/stats"

	"github.com/j7b/mailgun/client/mock"
)

var m = mock.Client

func TestDelete(t *testing.T) {
	if err := Delete(m(t), `tag`); err != nil {
		t.Fatal(err)
	}
}

func TestUpdate(t *testing.T) {
	if err := Update(m(t), `tag`, `description`); err != nil {
		t.Fatal(err)
	}
}

func TestGet(t *testing.T) {
	res, err := Get(m(t), `tag`)
	if err != nil {
		t.Fatal(err)
	}
	if res.Tag != `tag` {
		t.Fatal(`the universe no longer makes sense`)
	}
}

func TestList(t *testing.T) {
	res, err := List(m(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Tags) != 2 {
		t.Fatal(`want 2 got`, len(res.Tags))
	}
}

func TestDevices(t *testing.T) {
	res, err := Devices(m(t), `exampletag`)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatal(`want 2 got`, len(res))
	}
}

func TestProviders(t *testing.T) {
	res, err := Providers(m(t), `exampletag`)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatal(`want 2 got`, len(res))
	}
}

func TestCountries(t *testing.T) {
	res, err := Countries(m(t), `exampletag`)
	if err != nil {
		t.Fatal(err)
	}
	if len(res) != 2 {
		t.Fatal(`want 2 got`, len(res))
	}
}

func TestStats(t *testing.T) {
	res, err := Stats(m(t), `exampletag`, nil, nil, nil, stats.Delivered)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Stats) != 8 {
		t.Fatal(`want 8 got`, len(res.Stats))
	}
}
