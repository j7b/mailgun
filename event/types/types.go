// Package types contains event types.
package types

// Interface is a marker shared by
// all event types.
type Interface interface {
	event()
}

// Generic is shared by all events.
type Generic struct {
	Event     string  `json:"event"`
	ID        string  `json:"id"`
	Timestamp float64 `json:"timestamp"`
}

func (Generic) event() {
	_ = 1
}

// Geolocation contains location info.
type Geolocation struct {
	Country string `json:"country"` // "country": "US",
	Region  string `json:"region"`  // "region": "TX",
	City    string `json:"city"`    // "city": "San Antonio"
}

// ClientInfo has browser data.
type ClientInfo struct {
	Type       string `json:"client-type"` // "client-type": "browser",
	OS         string `json:"client-os"`   // "client-os": "OS X",
	DeviceType string `json:"device-type"` // "device-type": "desktop",
	Name       string `json:"client-name"` // "client-name": "Chrome",
	UserAgent  string `json:"user-agent"`  // "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1500.95 Safari/537.36"
}

// Route is an inbound mail route.
type Route struct {
	Priority    int      `json:"priority"`    // "priority": 1,
	Expression  string   `json:"expression"`  //       "expression": "match_recipient(\".*@samples.mailgun.org\")",
	Description string   `json:"description"` //       "description": "Sample route",
	Actions     []string `json:"actions"`     //       "actions": [
	//         "stop()",
	//         "forward(\"http://host.com/messages\")"
	//       ]
}

// Message identifies an email message.
type Message struct {
	Headers map[string]string `json:"headers"` //     "headers": {
	//       "to": "",
	//       "message-id": "77AF5C3CA1416D93FC47AF8AD42A60AD@example.com",
	//       "from": "John Doe <sender@example.com>",
	//       "subject": "Test Subject"
	//     },
	Attachments []interface{} `json:"attachments"` //     "attachments": [],
	Recipients  []string      `json:"recipients"`  //     "recipients": [
	//       "recipient@example.com"
	//     ],
	Size int `json:"size"` //     "size": 6021
}

// DeliveryStatus is reported by receiving MTA.
type DeliveryStatus struct {
	Message     string `json:"message"`     //     "message": "",
	Code        int    `json:"code"`        //     "code": 0,
	Description string `json:"description"` //     "description": null
}

// Storage is storage information.
type Storage struct {
	URL string `json:"url"` //       "url":"https://api.mailgun.net/v3/domains/ninomail.com/messages/WyI3MDhjODgwZTZlIiwgIjF6",
	Key string `json:"key"` //       "key":"WyI3MDhjODgwZTZlIiwgIjF6"
}

// Accepted event.
type Accepted struct {
	// {
	Generic //   "event": "accepted",
	//   "id": "ncV2XwymRUKbPek_MIM-Gw",
	//   "timestamp": 1377211256.096436,
	Tags     []string               `json:"tags"`     //   "tags": [],
	Envelope map[string]interface{} `json:"envelope"` //   "envelope": {
	//     "sender": "sender@example.com"
	//   },
	Campaigns     interface{} `json:"campaigns"`      //   "campaigns": [],
	UserVariables interface{} `json:"user-variables"` //   "user-variables": {},
	Flags         interface{} `json:"flags"`          //   "flags": {
	//     "is-authenticated": false,
	//     "is-test-mode": false
	//   },
	Routes []Route `json:"routes"` //   "routes": [
	//     {
	//       "priority": 1,
	//       "expression": "match_recipient(\".*@samples.mailgun.org\")",
	//       "description": "Sample route",
	//       "actions": [
	//         "stop()",
	//         "forward(\"http://host.com/messages\")"
	//       ]
	//     }
	//   ],
	Message Message `json:"message"` //   "message": {
	//     "headers": {
	//       "to": "",
	//       "message-id": "77AF5C3CA1416D93FC47AF8AD42A60AD@example.com",
	//       "from": "John Doe <sender@example.com>",
	//       "subject": "Test Subject"
	//     },
	//     "attachments": [],
	//     "recipients": [
	//       "recipient@example.com"
	//     ],
	//     "size": 6021
	//   },
	Recipient string `json:"recipient"` //   "recipient": "recipient@example.com",
	Method    string `json:"method"`    //   "method": "smtp"
	// }
	//
}

// Delivered event.
type Delivered struct {
	// {
	Generic //   "event": "delivered",
	//   "id": "W3X4JOhFT-OZidZGKKr9iA",
	//   "timestamp": 1377208314.173742,
	Tags     []string               `json:"tags"`     //   "tags": [],
	Envelope map[string]interface{} `json:"envelope"` //   "envelope": {
	//     "transport": "smtp",
	//     "sender": "postmaster@samples.mailgun.org",
	//     "sending-ip": "184.173.153.199"
	//   },
	Status DeliveryStatus `json:"delivery-status"` //   "delivery-status": {
	//     "message": "",
	//     "code": 0,
	//     "description": null
	//   },
	Campaigns     interface{} `json:"campaigns"`      //   "campaigns": [],
	UserVariables interface{} `json:"user-variables"` //   "user-variables": {},
	Flags         interface{} `json:"flags"`          //   "flags": {}
	Message       Message     `json:"message"`        //   "message": {
	//     "headers": {
	//       "to": "recipient@example.com",
	//       "message-id": "20130822215151.29325.59996@samples.mailgun.org",
	//       "from": "sender@example.com",
	//       "subject": "Sample Message"
	//     },
	//     "attachments": [],
	//     "recipients": [
	//       "recipient@example.com"
	//     ],
	//     "size": 31143
	//   },
	Recipient string `json:"recipient"` //   "recipient": "recipient@example.com",
	// }
	//
}

// Failed event.
type Failed struct {
	// {
	Generic //   "event": "failed",
	//   "id": "pVqXGJWhTzysS9GpwF2hlQ",
	//   "timestamp": 1377198389.769129,
	Severity string                 `json:"severity"` //   "severity": "permanent",
	Tags     []string               `json:"tags"`     //   "tags": [],
	Envelope map[string]interface{} `json:"envelope"` //   "envelope": {
	//     "transport": "smtp",
	//     "sender": "postmaster@samples.mailgun.org",
	//     "sending-ip": "184.173.153.199"
	//   },
	Status DeliveryStatus `json:"delivery-status"` //   "delivery-status": {
	//     "message": "Relay Not Permitted",
	//     "code": 550,
	//     "description": null
	//   },
	Campaigns     interface{} `json:"campaigns"`      //   "campaigns": [],
	Reason        string      `json:"reason"`         //   "reason": "bounce",
	UserVariables interface{} `json:"user-variables"` //   "user-variables": {},
	Flags         interface{} `json:"flags"`          //   "flags": {
	//     "is-authenticated": true,
	//     "is-test-mode": false
	//   },
	Message Message `json:"message"` //   "message": {
	//     "headers": {
	//       "to": "recipient@example.com",
	//       "message-id": "20130822185902.31528.73196@samples.mailgun.org",
	//       "from": "John Doe <sender@example.com>",
	//       "subject": "Test Subject"
	//     },
	//     "attachments": [],
	//     "recipients": [
	//       "recipient@example.com"
	//     ],
	//     "size": 557
	//   },
	Recipient string `json:"recipient"` //   "recipient": "recipient@example.com",
	// }
	//
}

// Opened event.
type Opened struct {
	// {
	Generic //   "event": "opened",
	//   "id": "-laxIqj9QWubsjY_3pTq_g",
	//   "timestamp": 1377047343.042277,
	Recipient   string       `json:"recipient"`             //   "recipient": "recipient@example.com",
	Geolocation *Geolocation `json:"geolocation,omitempty"` //   "geolocation": {
	//     "country": "US",
	//     "region": "Texas",
	//     "city": "Austin"
	//   },
	Tags          []string    `json:"tags"`                  //   "tags": [],
	Campaigns     interface{} `json:"campaigns"`             //   "campaigns": [],
	UserVariables interface{} `json:"user-variables"`        //   "user-variables": {},
	IP            string      `json:"ip"`                    //   "ip": "111.111.111.111",
	ClientInfo    *ClientInfo `json:"client-info,omitempty"` //   "client-info": {
	//     "client-type": "mobile browser",
	//     "client-os": "iOS",
	//     "device-type": "mobile",
	//     "client-name": "Mobile Safari",
	//     "user-agent": "Mozilla/5.0 (iPhone; CPU iPhone OS 6_1 like Mac OS X) AppleWebKit/536.26 (KHTML, like Gecko) Mobile/10B143"
	//   },
	Message Message `json:"message"` //   "message": {
	//     "headers": {
	//       "message-id": "20130821005614.19826.35976@samples.mailgun.org"
	//     }
	//   },
	// }
}

// Clicked event.
type Clicked struct {
	// {
	Generic //   "event": "clicked",
	//   "id": "G5zMz2ysS6OxZ2C8xb2Tqg",
	//   "timestamp": 1377075564.094891,
	Recipient   string       `json:"recipient"`             //   "recipient": "recipient@example.com",
	Geolocation *Geolocation `json:"geolocation,omitempty"` //   "geolocation": {
	//     "country": "US",
	//     "region": "TX",
	//     "city": "Austin"
	//   },
	Tags          []string    `json:"tags"`                  //   "tags": [],
	URL           string      `json:"url"`                   //   "url": "http://google.com",
	IP            string      `json:"ip"`                    //   "ip": "127.0.0.1",
	Campaigns     interface{} `json:"campaigns"`             //   "campaigns": [],
	UserVariables interface{} `json:"user-variables"`        //   "user-variables": {},
	ClientInfo    *ClientInfo `json:"client-info,omitempty"` //   "client-info": {
	//     "client-type": "browser",
	//     "client-os": "Linux",
	//     "device-type": "desktop",
	//     "client-name": "Chromium",
	//     "user-agent": "Mozilla/5.0 (X11; Linux i686) AppleWebKit/537.36 (KHTML, like Gecko) Ubuntu Chromium/28.0.1500.71 Chrome/28.0.1500.71 Safari/537.36"
	//   },
	Message Message `json:"message"` //   "message": {
	//     "headers": {
	//       "message-id": "20130821085807.30688.67706@samples.mailgun.org"
	//     }
	//   },
	// }
}

// Unsubscribed event.
type Unsubscribed struct {
	// {
	Generic //   "event": "unsubscribed",
	//   "id": "W3X4JOhFT-OZidZGKKr9iA",
	//   "timestamp": 1377213791.421473,
	Recipient   string       `json:"recipient"`             //   "recipient": "recipient@example.com",
	Geolocation *Geolocation `json:"geolocation,omitempty"` //   "geolocation": {
	//     "country": "US",
	//     "region": "TX",
	//     "city": "San Antonio"
	//   },
	Campaigns     interface{} `json:"campaigns"`             //   "campaigns": [],
	Tags          []string    `json:"tags"`                  //   "tags": [],
	UserVariables interface{} `json:"user-variables"`        //   "user-variables": {},
	IP            string      `json:"ip"`                    //   "ip": "50.51.14.451",
	ClientInfo    *ClientInfo `json:"client-info,omitempty"` //   "client-info": {
	//     "client-type": "browser",
	//     "client-os": "OS X",
	//     "device-type": "desktop",
	//     "client-name": "Chrome",
	//     "user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_8_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/28.0.1500.95 Safari/537.36"
	//   },
	Message Message `json:"message"` //   "message": {
	//     "headers": {
	//       "message-id": "20130822232216.13966.79700@samples.mailgun.org"
	//     }
	//   },
	// }
}

// Complained event.
type Complained struct {
	// {
	Generic //   "event": "complained",
	//   "id": "ncV2XwymRUKbPek_MIM-Gw",
	//   "timestamp": 1377214260.049634,
	Recipient     string      `json:"recipient"`      //   "recipient": "foo@example.com",
	Tags          []string    `json:"tags"`           //   "tags": [],
	Campaigns     interface{} `json:"campaigns"`      //   "campaigns": [],
	UserVariables interface{} `json:"user-variables"` //   "user-variables": {},
	Flags         interface{} `json:"flags"`          //   "flags": {
	//     "is-test-mode": false
	//   },
	Message Message `json:"message"` //   "message": {
	//     "headers": {
	//       "to": "foo@example.com",
	//       "message-id": "20130718032413.263EE2E0926@example.com",
	//       "from": "John Doe <sender@example.com>",
	//       "subject": "This is the subject."
	//     },
	//     "attachments": [],
	//     "size": 18937
	//   },
	// }
}

// Stored event.
type Stored struct {
	// {
	Generic //    "event":"stored",
	//    "id": "czsjqFATSlC3QtAK-C80nw",
	//    "timestamp":1378335036.859382,
	Storage Storage `json:"storage"` //    "storage":{
	//       "url":"https://api.mailgun.net/v3/domains/ninomail.com/messages/WyI3MDhjODgwZTZlIiwgIjF6",
	//       "key":"WyI3MDhjODgwZTZlIiwgIjF6"
	//    },
	Campaigns     interface{} `json:"campaigns"`      //    "campaigns":[],
	UserVariables interface{} `json:"user-variables"` //    "user-variables":{},
	Flags         interface{} `json:"flags"`          //    "flags":{
	//       "is-test-mode":false
	//    },
	Tags    []string `json:"tags"`    //    "tags":[],
	Message Message  `json:"message"` //    "message":{
	//       "headers":{
	//          "to":"satshabad <satshabad@mailgun.com>",
	//          "message-id":"CAC8xyJxAO7Y0sr=3r-rJ4C6ULZs3cSVPPqYEXLHtarKOKaOCKw@mail.gmail.com",
	//          "from":"Someone <someone@example.com>",
	//          "subject":"Re: A TEST"
	//       },
	//       "attachments":[],
	//       "recipients":[
	//          "satshabad@mailgun.com"
	//       ],
	//       "size":2566
	//    },
	// }
}

// BUG(j7b): Most of the applicable documentation for these types
// is non-normative and by example, as a result the corresponding
// Go types might not be optimal. As better types can be ascertained,
// methods or functions may be added to this package in the interest
// of increasing utility while preserving backwards compatibility.
