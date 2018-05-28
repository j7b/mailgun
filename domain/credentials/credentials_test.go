package credentials

import (
	"testing"

	"github.com/j7b/mailgun/client/mock"
)

func TestAll(t *testing.T) {
	errs := func(e error) func(*testing.T) {
		return func(t *testing.T) {
			if e != nil {
				t.Fatal(e)
			}
		}
	}
	c := mock.Client(t)
	m := map[string]error{
		"Delete": Delete(c, `login`),
		"Update": Update(c, `login`, `password`),
		"Create": Create(c, `login`, `password`),
	}
	for k, v := range m {
		t.Run(k, errs(v))
	}
	t.Run(`Credentials`, func(t *testing.T) {
		creds, err := Credentials(c, 1)
		if err != nil {
			t.Fatal(err)
		}
		if len(creds) != 2 {
			t.Fatal("want 2 creds got", len(creds))
		}
	})
}
