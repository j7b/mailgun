package domain_test

import (
	"testing"

	"github.com/j7b/mailgun/client/mock"
	"github.com/j7b/mailgun/domain"
)

func TestAPI(t *testing.T) {
	c := mock.Client(t)
	api := domain.API(c)
	t.Run("List", func(t *testing.T) {
		d, err := api.List(1)
		if err != nil {
			t.Fatal(err)
		}
		if d[0].Name != "samples.mailgun.org" {
			t.Fatal(`want "samples.mailgun.org" got`, d[0].Name)
		}
	})
	t.Run("Delete", func(t *testing.T) {
		if err := api.Delete(`domain.mock`); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("Get", func(t *testing.T) {
		d, err := api.Get(`domain.mock`)
		if err != nil {
			t.Fatal(err)
		}
		if d.Domain.Name != `domain.com` {
			t.Fatal(`want domain.com got `, d.Domain.Name)
		}
	})
	t.Run("New", func(t *testing.T) {
		d, err := api.New(`domain.mock`, `dunno`, domain.ActionBlock, false)
		if err != nil {
			t.Fatal(err)
		}
		if d.Domain.Name != `domain.com` {
			t.Fatal(`want domain.com got `, d.Domain.Name)
		}
	})
}
