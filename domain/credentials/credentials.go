// Package credentials implements domain credentials.
package credentials

import (
	"fmt"

	"github.com/j7b/mailgun/client"
)

// Credential is domain credentials.
type Credential struct {
	SizeBytes int    `json:"size_bytes"` // "size_bytes": 0,
	CreatedAt string `json:"created_at"` // "created_at": "Tue, 27 Sep 2011 20:24:22 GMT",
	Maibox    string `json:"mailbox"`    // "mailbox": "user@samples.mailgun.org"
	Login     string `json:"login"`      // "login": "user@samples.mailgun.org"
}

// Delete deletes credentials for login.
func Delete(c client.Caller, login string) error {
	return c.Delete(`/domains`, c.Domain(), `credentials`, login).Err()
}

// Update changes the password for login.
func Update(c client.Caller, login, password string) error {
	return c.Put(`/domains`, c.Domain(), `credentials`, login).
		SetForm("password", password).
		Err()
}

// Create creates a set of credentials.
func Create(c client.Caller, login, password string) error {
	return c.Post(`/domains`, c.Domain(), `credentials`).
		SetForm("login", login).
		SetForm("password", password).
		Err()
}

// Credentials returns a page of Credential for the API domain.
func Credentials(c client.Caller, page int) ([]Credential, error) {
	if page < 1 {
		return nil, fmt.Errorf("List: %v is < 1 ", page)
	}
	page--
	var creds struct {
		Credentials []Credential `json:"items"`
	}
	return creds.Credentials, c.Get(`/domains`, c.Domain(), `credentials`).
		SetQuery("skip", fmt.Sprintf(`%v`, page*100)).
		Decode(&creds)
}
