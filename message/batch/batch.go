// Package batch implements the batch mailing convention.
/*
The Send method sends a Message to the recipients in
recipmap. The recipmap key should be an email address
(bare or RFC5322 compliant). The value (map[string]interface{})
will be marshalled to the recipient-variables
field.
*/
package batch

import (
	"github.com/j7b/mailgun/client"
	"github.com/j7b/mailgun/message"
)

type recipvars interface {
	SetRecipVars(m map[string]map[string]interface{})
}

func rv(i interface{}, recipmap map[string]map[string]interface{}) {
	rv := i.(recipvars)
	rv.SetRecipVars(recipmap)
}

// Send sends m to recipmap.
func Send(c client.Caller, m *message.Message, recipmap map[string]map[string]interface{}) (*message.Response, error) {
	rv(m, recipmap)
	recips := make([]string, 0, len(recipmap))
	for k := range recipmap {
		recips = append(recips, k)
	}
	if _, err := m.To(recips...); err != nil {
		return nil, err
	}
	return m.Send(c)
}

// BUG(j7b): The results of using non-builtin types in the
// map[string]interface{} can be astonishing and the
// templating supported by the endpoint is very limited.
