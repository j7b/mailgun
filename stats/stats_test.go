package stats

import (
	"testing"

	"github.com/j7b/mailgun/client/mock"
)

func TestStats(t *testing.T) {
	c := mock.Client(t)
	res, err := Query(c, nil, nil, nil, AllEvents()...)
	if err != nil {
		t.Fatal(err)
	}
	if l := len(res.Stats); l != 1 {
		t.Fatal(`want`, 1, `got`, l)
	}
}
