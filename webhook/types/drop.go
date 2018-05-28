package types

import "net/textproto"

var _ = textproto.Error{}

type Drop struct {
	Event           `json:"event"`                   // Event name (“dropped”).
	Recipient       string                           `json:"recipient"`       // Intended recipient.
	Domain          string                           `json:"domain"`          // Domain that sent the original message.
	MessageHeaders  string                           `json:"message-headers"` // String list of all MIME headers of the original message dumped to a JSON string (order of headers preserved).
	Reason          string                           `json:"reason"`          // Reason for failure. Can be one either “hardfail” or “old”. See below.
	Code            string                           `json:"code"`            // ESP response code, e.g. if the message was blocked as a spam (optional).
	Description     string                           `json:"description"`     // Detailed explanation of why the messages was dropped
	CustomVariables map[string]interface{}           `json:"-"`               // Your own custom JSON object included in the header (see Attaching Data to Messages).
	Timestamp       `json:"timestamp"`               // Number of seconds passed since January 1, 1970 (see securing webhooks).
	Token           `json:"token"`                   // Randomly generated string with length 50 (see securing webhooks).
	Signature       `json:"signature"`               // String with hexadecimal digits generate by HMAC algorithm (see securing webhooks).
	Attachments     map[*textproto.MIMEHeader][]byte `json:"-"` // attached file (‘x’ stands for number of the attachment). Attachments are included if the recipient ESP includes them in the bounce message. They are handled as file uploads, encoded as multipart/form-data.
}

func (d Drop) Drop() Drop { return d }

type IDrop interface{ Drop() Drop }
