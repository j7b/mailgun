package batch

import (
	"testing"

	"github.com/j7b/mailgun/client/mock"
	"github.com/j7b/mailgun/message"
)

func TestSend(t *testing.T) {
	c := mock.Client(t)
	msg, err := message.New(`Mailgun <postmaster@sandbox.mailgun.org>`,
		`Hey bud!`,
		`You you are <b>truly</b> awesome!`)
	if err != nil {
		t.Fatal(err)
	}
	if res, err := Send(c, msg, map[string]map[string]interface{}{`rick@roll.org`: nil}); err != nil {
		t.Fatal(err)
	} else {
		if res.ID != `<20111114174239.25659.5817@samples.mailgun.org>` {
			t.Fatal(`want`, `<20111114174239.25659.5817@samples.mailgun.org>`, `got`, res.ID)
		}
	}
}
