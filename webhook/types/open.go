package types

import "net/textproto"

var _ = textproto.Error{}

type Open struct {
	Event           `json:"event"`         // Event name (“opened”).
	Recipient       string                 `json:"recipient"`     // Recipient who opened.
	Domain          string                 `json:"domain"`        // Domain that sent the original message.
	IP              string                 `json:"ip"`            // IP address the event originated from.
	Country         string                 `json:"country"`       // Two-letter country code (as specified by ISO3166) the event came from or ‘unknown’ if it couldn’t be determined.
	Region          string                 `json:"region"`        // Two-letter or two-digit region code or ‘unknown’ if it couldn’t be determined.
	City            string                 `json:"city"`          // Name of the city the event came from or ‘unknown’ if it couldn’t be determined.
	UserAgent       string                 `json:"user-agent"`    // User agent string of the client triggered the event.
	DeviceType      string                 `json:"device-type"`   // Device type the email was opened on. Can be ‘desktop’, ‘mobile’, ‘tablet’, ‘other’ or ‘unknown’.
	ClientType      string                 `json:"client-type"`   // Type of software the email was opened in, e.g. ‘browser’, ‘mobile browser’, ‘email client’.
	ClientName      string                 `json:"client-name"`   // Name of the client software, e.g. ‘Thunderbird’, ‘Chrome’, ‘Firefox’.
	ClientOS        string                 `json:"client-os"`     // OS family running the client software, e.g. ‘Linux’, ‘Windows’, ‘OSX’.
	CampaignID      string                 `json:"campaign-id"`   // The id of campaign triggering the event.
	CampaignName    string                 `json:"campaign-name"` // The name of campaign triggering the event.
	Tag             string                 `json:"tag"`           // Message tag, if message was tagged. See Tagging
	MailingList     string                 `json:"mailing-list"`  // The address of mailing list the original message was sent to.
	CustomVariables map[string]interface{} `json:"-"`             // Your own custom JSON object included in the header (see Attaching Data to Messages).
	Timestamp       `json:"timestamp"`     // Number of seconds passed since January 1, 1970 (see securing webhooks).
	Token           `json:"token"`         // Randomly generated string with length 50 (see securing webhooks).
	Signature       `json:"signature"`     // String with hexadecimal digits generate by HMAC algorithm (see securing webhooks).
}

func (o Open) Open() Open { return o }

type IOpen interface{ Open() Open }
