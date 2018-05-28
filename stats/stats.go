// Package stats implements the stats/total endpoint.
package stats

import (
	"fmt"
	"time"

	"github.com/j7b/mailgun/client"
	"github.com/j7b/mailgun/stats/types"
)

// Resolution represents reporting resolution.
type Resolution interface {
	name() string
}

type resolution string

func (resolution) name() string {
	return `resolution`
}

// Reporting resolution.
const (
	Hour  = resolution(`hour`)
	Day   = resolution(`day`)
	Month = resolution(`month`)
)

// Event represents a reported event.
type Event interface {
	event()
}

type event string

func (event) event() {}

// Reporting event.
const (
	Accepted     = event(`accepted`)
	Delivered    = event(`delivered`)
	Failed       = event(`failed`)
	Opened       = event(`opened`)
	Clicked      = event(`clicked`)
	Unsubscribed = event(`unsubscribed`)
	Complained   = event(`complained`)
	Stored       = event(`stored`)
)

// AllEvents returns all Event types.
func AllEvents() []Event {
	return []Event{Accepted, Delivered, Failed, Opened, Clicked, Unsubscribed, Complained, Stored}
}

// Events are the count of events that occurred.
type Events struct {
	Accepted     int `json:"accepted"`
	Delivered    int `json:"delivered"`
	Failed       int `json:"failed"`
	Opened       int `json:"opened"`
	Clicked      int `json:"clicked"`
	Unsubscribed int `json:"unsubscribed"`
	Complained   int `json:"complained"`
}

func dedup(events []Event) ([]Event, error) {
	if len(events) == 0 {
		return nil, fmt.Errorf("query: no events specified")
	}
	m := make(map[Event]bool)
	for _, ev := range events {
		m[ev] = true
	}
	events = nil
	for ev := range m {
		events = append(events, ev)
	}
	return events, nil
}

// Stats are the statistics returned by the Query function.
type Stats struct {
	Start      string        `json:"start"`
	End        string        `json:"end"`
	Resolution string        `json:"resolution"`
	Stats      []types.Stats `json:"stats"`
}

// Query calls the stats/total endpoint.
func Query(c client.Caller, start *time.Time, end *time.Time, reso Resolution, events ...Event) (*Stats, error) {
	events, err := dedup(events)
	if err != nil {
		return nil, err
	}
	req := c.Get(`stats/total`)
	for _, ev := range events {
		req.AddQuery("event", fmt.Sprintf(`%s`, ev))
	}
	if start != nil {
		req.SetQuery("start", start.Format(time.RFC1123))
	}
	if end != nil {
		req.SetQuery("end", start.Format(time.RFC1123))
	}
	if reso != nil {
		req.SetQuery(`resolution`, fmt.Sprintf(`%s`, reso))
	}
	var stats *Stats
	return stats, req.Decode(&stats)
}

// BUG(j7b): The "duration" parameter is not implemented,
// it seems redundant.
