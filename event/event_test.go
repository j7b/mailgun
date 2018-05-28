package event_test

import (
	"log"
	"testing"
	"time"

	"github.com/j7b/mailgun"
	"github.com/j7b/mailgun/client/mock"
	"github.com/j7b/mailgun/event"
	"github.com/j7b/mailgun/event/types"
)

func TestQuery(t *testing.T) {
	c := mock.Client(t)
	q := event.Queries(c)
	end := time.Now()
	start := end.Add(time.Hour*24 - 30)
	whynot := true
	res, err := q.Query(&start, &end, &whynot, event.Tags(`summer-sale-2018`), event.Event(`delivered`))
	if err != nil {
		t.Fatal(err)
	}
	if l := len(res.List()); l != 2 {
		t.Fatal(`want 2 events got`, l)
	}
}

func dosomethingwith(interface{}) {}

func Example_event_Query() {
	c, err := mailgun.New(``, ``) // use env variables
	if err != nil {
		log.Fatal(err)
	}
	q := event.Queries(c)
	res, err := q.Query(nil, nil, nil, event.Tags(`summer-sale-2018`), event.Event(`delivered`))
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range res.List() {
		switch t := i.(type) {
		case types.Delivered:
			dosomethingwith(t)
		}
	}
}
