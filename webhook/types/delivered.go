package types

import "net/textproto"

var _ = textproto.Error{}

type Delivered struct {
	Event           `json:"event"`         // Event name (“delivered”).
	Recipient       string                 `json:"recipient"`       // Intended recipient.
	Domain          string                 `json:"domain"`          // Domain that sent the original message.
	MessageHeaders  string                 `json:"message-headers"` // String list of all MIME headers dumped to a JSON string (order of headers preserved).
	MessageID       string                 `json:"Message-Id"`      // String id of the original message delivered to the recipient.
	CustomVariables map[string]interface{} `json:"-"`               // Your own custom JSON object included in the header of the original message (see Attaching Data to Messages).
	Timestamp       `json:"timestamp"`     // Number of seconds passed since January 1, 1970 (see securing webhooks).
	Token           `json:"token"`         // Randomly generated string with length 50 (see securing webhooks).
	Signature       `json:"signature"`     // String with hexadecimal digits generate by HMAC algorithm (see securing webhooks).
}

func (d Delivered) Delivered() Delivered { return d }

type IDelivered interface{ Delivered() Delivered }
