package mailgun

import "testing"

func TestNew(t *testing.T) {
	c, err := New(`asdfe`, `asdf`)
	if err != nil {
		t.Fatal(err)
	}
	if c.Key() != `asdfe` {
		t.Fatal(`want "asdfe" got`, c.Key())
	}
}
