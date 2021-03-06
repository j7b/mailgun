package types

import "net/textproto"

var _ = textproto.Error{}

type Bounce struct {
	Event           `json:"event"`                   // Event name (“bounced”).
	Recipient       string                           `json:"recipient"`       // Recipient who could not be reached.
	Domain          string                           `json:"domain"`          // Domain that sent the original message.
	MessageHeaders  string                           `json:"message-headers"` // String list of all MIME headers of the original message dumped to a JSON string (order of headers preserved).
	Code            string                           `json:"code"`            // SMTP bounce error code in form (X.X.X).
	Error           string                           `json:"error"`           // SMTP bounce error string.
	Notification    string                           `json:"notification"`    // Detailed reason for bouncing (optional).
	CampaignID      string                           `json:"campaign-id"`     // The id of campaign triggering the event.
	CampaignName    string                           `json:"campaign-name"`   // The name of campaign triggering the event.
	Tag             string                           `json:"tag"`             // Message tag, if it was tagged. See Tagging.
	MailingList     string                           `json:"mailing-list"`    // The address of mailing list the original message was sent to.
	CustomVariables map[string]interface{}           `json:"-"`               // Your own custom JSON object included in the header (see Attaching Data to Messages).
	Timestamp       `json:"timestamp"`               // Number of seconds passed since January 1, 1970 (see securing webhooks).
	Token           `json:"token"`                   // Randomly generated string with length 50 (see securing webhooks).
	Signature       `json:"signature"`               // String with hexadecimal digits generate by HMAC algorithm (see securing webhooks).
	Attachments     map[*textproto.MIMEHeader][]byte `json:"-"` // attached file (‘x’ stands for number of the attachment). Attachments are included if the recipient ESP includes them in the bounce message. They are handled as file uploads, encoded as multipart/form-data.
}

func (b Bounce) Bounce() Bounce { return b }

type IBounce interface{ Bounce() Bounce }
