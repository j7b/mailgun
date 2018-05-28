// Package tags contains types and implementations
// related to message tags.
package tags

import (
	"fmt"
	"time"

	"github.com/j7b/mailgun/client"
	"github.com/j7b/mailgun/client/pager"
	"github.com/j7b/mailgun/stats"
)

// Tag is a domain tag.
type Tag struct {
	Tag         string `json:"tag"`         // "tag": "red",
	Description string `json:"description"` // "description": "red signup button",
}

// Tags contains Tags and paging info returned by queries.
type Tags struct {
	Tags []Tag `json:"items"`
	pager.Pager
}

// Next returns the next page.
func (t *Tags) Next() (*Tags, error) {
	var tags *Tags
	return tags, t.Paging.Next(&tags)
}

// Previous returns the previous page.
func (t *Tags) Previous() (*Tags, error) {
	var tags *Tags
	return tags, t.Paging.Previous(&tags)
}

// First returns the first page.
func (t *Tags) First() (*Tags, error) {
	var tags *Tags
	return tags, t.Paging.First(&tags)
}

// Last returns the last page.
func (t *Tags) Last() (*Tags, error) {
	var tags *Tags
	return tags, t.Paging.Last(&tags)
}

// Delete deletes tag by name.
func Delete(c client.Caller, name string) error {
	return c.Delete(`tags`, name).Err()
}

// Update updates a tag description.
func Update(c client.Caller, name, description string) error {
	return c.Put(`tags`, name).SetForm("description", description).
		Err()
}

// Get returns a single tag.
func Get(c client.Caller, name string) (*Tag, error) {
	var tag *Tag
	return tag, c.Get(`tags`, name).Decode(&tag)
}

// List returns Tags for API domain.
func List(c client.Caller) (*Tags, error) {
	var tags *Tags
	return tags, c.Get(`tags`).Decode(&tags)
}

// Devices retrieves aggregate stats for devices.
func Devices(c client.Caller, name string) (map[string]stats.Events, error) {
	var o struct {
		Devices map[string]stats.Events `json:"device,omitempty"` // docs are wrong
	}
	return o.Devices, c.Get(`tags`, name, `stats/aggregates/devices`).Decode(&o)
}

// Providers retrieves aggregate stats for providers.
func Providers(c client.Caller, name string) (map[string]stats.Events, error) {
	var o struct {
		Providers map[string]stats.Events `json:"provider,omitempty"` // docs are wrong
	}
	return o.Providers, c.Get(`tags`, name, `stats/aggregates/providers`).Decode(&o)
}

// Countries retrieves aggregate stats for countries.
func Countries(c client.Caller, name string) (map[string]stats.Events, error) {
	var o struct {
		Countries map[string]stats.Events `json:"country,omitempty"` // docs are wrong
	}
	return o.Countries, c.Get(`tags`, name, `stats/aggregates/countries`).Decode(&o)
}

// Stats retrieves stats for name. The ev parameter must not be nil, reso may be nil.
func Stats(c client.Caller, name string, start *time.Time, end *time.Time, reso stats.Resolution, ev ...stats.Event) (*stats.Stats, error) {
	if len(ev) == 0 {
		return nil, fmt.Errorf("Stats: ev is required")
	}
	req := c.Get(`tags`, name, `stats`)
	vals := req.Query()
	for _, e := range ev {
		vals.Add("event", fmt.Sprintf("%s", e))
	}
	if start != nil {
		vals.Set("start", start.Format(time.RFC1123))
	}
	if end != nil {
		vals.Set("end", end.Format(time.RFC1123))
	}
	if reso != nil {
		vals.Set("resolution", fmt.Sprintf(`%s`, reso))
	}
	var s *stats.Stats
	return s, req.Decode(&s)
}

// BUG(j7b): Stats duration parameter not implemented, seems redundant.
