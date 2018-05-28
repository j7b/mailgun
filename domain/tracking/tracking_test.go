package tracking

import (
	"testing"

	"github.com/j7b/mailgun/client/mock"
)

func TestAll(t *testing.T) {
	c := mock.Client(t)
	t.Run(`Settings`, func(t *testing.T) {
		tracking, err := Settings(c)
		if err != nil {
			t.Fatal(err)
		}
		if !tracking.Open.Active {
			t.Fatal("expect open tracking active")
		}
		if len(tracking.Unsubscribe.HTMLFooter) != 57 {
			t.Fatal("want html footer len 57 got", len(tracking.Unsubscribe.HTMLFooter))
		}
	})
	errs := func(e error) func(*testing.T) {
		return func(t *testing.T) {
			if e != nil {
				t.Fatal(e)
			}
		}
	}
	m := map[string]error{
		"Clicks": Clicks(c, true),
		"Opens":  Opens(c, true),
		"Unsubs": Unsubs(c, Unsub{Type: Type{Active: true}, HTMLFooter: `derp`, TextFooter: `derp`}),
	}
	for k, v := range m {
		t.Run(k, errs(v))
	}
}
