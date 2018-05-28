// Package types contains stats types.
package types

// Accepted stats.
type Accepted struct {
	Incoming int `json:"incoming"`
	Outgoing int `json:"outgoing"`
	Total    int `json:"total"`
}

// Clicked stats.
type Clicked struct {
	Total int `json:"total"`
}

// Complained stats.
type Complained struct {
	Total int `json:"total"`
}

// Delivered stats.
type Delivered struct {
	HTTP  int `json:"http"`
	SMTP  int `json:"smtp"`
	Total int `json:"total"`
}

// Failed stats.
type Failed struct {
	Permanent Permanent `json:"permanent"`
	Temporary Temporary `json:"temporary"`
}

// Temporary failures.
type Temporary struct {
	ESPBlock int `json:"espblock"`
	Total    int `json:"total"`
}

// Permanent failures.
type Permanent struct {
	Bounce              int `json:"bounce"`
	DelayedBounce       int `json:"delayed-bounce"`
	SuppressBounce      int `json:"suppress-bounce"`
	SuppressComplaint   int `json:"suppress-complaint"`
	SuppressUnsubscribe int `json:"suppress-unsubscribe"`
	Total               int `json:"http"`
}

// Opened stats.
type Opened struct {
	Total int `json:"total"`
}

// Unsubscribed stats.
type Unsubscribed struct {
	Total int `json:"total"`
}

// Stats contains all stats returned by the endpoint.
type Stats struct {
	Time         string       `json:"time"`
	Accepted     Accepted     `json:"accepted"`
	Clicked      Clicked      `json:"clicked"`
	Complained   Complained   `json:"complained"`
	Delivered    Delivered    `json:"delivered"`
	Failed       Failed       `json:"failed"`
	Opened       Opened       `json:"opened"`
	Unsubscribed Unsubscribed `json:"unsubscribed"`
}
