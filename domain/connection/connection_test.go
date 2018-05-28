package connection

import (
	"testing"

	"github.com/j7b/mailgun/client/mock"
)

func TestAll(t *testing.T) {
	c := mock.Client(t)
	t.Run(`Settings`, func(t *testing.T) {
		con, err := Settings(c)
		if err != nil {
			t.Fatal(err)
		}
		if !con.RequireTLS {
			t.Fatal("con not requiretls")
		}
	})
	t.Run(`Update`, func(t *testing.T) {
		if err := Update(c, &Connection{}); err != nil {
			t.Fatal(err)
		}
	})
}
