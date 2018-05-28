package message

import (
	"testing"

	"github.com/j7b/mailgun/client/mock"
)

func TestMessage(t *testing.T) {
	msg, err := New(`Mailgun <postmaster@sandbox.mailgun.org>`,
		`Hey bud!`,
		`You you are <b>truly</b> awesome!`,
		`rick@roll.net`)
	if err != nil {
		t.Fatal(err)
	}
	old, err := msg.To(`a@b.din`, `b@a.din`)
	if err != nil {
		t.Fatal(err)
	}
	if len(old) != 1 {
		t.Fatal(`old want 1 got`, len(old))
	}
	if len(msg.to) != 2 {
		t.Fatal(`new want 2 got`, len(msg.to))
	}
}

func TestSend(t *testing.T) {
	c := mock.Client(t)
	msg, err := New(`Mailgun <postmaster@sandbox.mailgun.org>`,
		`Hey bud!`,
		`You you are <b>truly</b> awesome!`,
		`rick@roll.net`)
	if err != nil {
		t.Fatal(err)
	}
	if res, err := msg.Send(c); err != nil {
		t.Fatal(err)
	} else {
		if res.ID != `<20111114174239.25659.5817@samples.mailgun.org>` {
			t.Fatal(`want`, `<20111114174239.25659.5817@samples.mailgun.org>`, `got`, res.ID)
		}
	}
}
