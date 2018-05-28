// Package connection implements domain connection parameters.
package connection

import (
	"fmt"

	"github.com/j7b/mailgun/client"
)

// Connection is domain connection parameters.
type Connection struct {
	RequireTLS       bool `json:"require_tls"`       // "require_tls": false,
	SkipVerification bool `json:"skip_verification"` // "skip_verification": false
}

// Settings retrieve connection settings for API domain.
func Settings(c client.Caller) (*Connection, error) {
	var con struct {
		Connection *Connection `json:"connection,omitempty"`
	}
	return con.Connection, c.Get(`/domains`, c.Domain(), `connection`).Decode(&con)
}

// Update updates connection parameters with con.
func Update(c client.Caller, con *Connection) error {
	return c.Put(`/domains`, c.Domain(), `connection`).
		SetForm("require_tls", fmt.Sprintf(`%v`, con.RequireTLS)).
		SetForm("skip_verification", fmt.Sprintf(`%v`, con.SkipVerification)).
		Err()
}
