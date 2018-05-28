// Package message implements the mailgun message convention.
package message

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"os"
	"time"

	"github.com/j7b/mailgun/client"
)

// TrackOption is the interface shared by tracking options.
type TrackOption interface {
	o()
}

type trackoption string

func (trackoption) o() {}

// TrackOptions
const (
	Track    = trackoption(`yes`)
	NoTrack  = trackoption(`no`)
	HTMLOnly = trackoption(`htmlonly`)
)

type recipvars map[string]map[string]interface{}

func (r *recipvars) SetRecipVars(m map[string]map[string]interface{}) {
	*r = recipvars(m)
}

// Message is an email message. Methods with
// slice parameters do not retain references
// to those slices.
type Message struct {
	from          string
	to            []string
	cc            []string
	bcc           []string
	Subject       string
	Text          string
	HTML          string
	Attachments   map[string]io.Reader
	Inlines       map[string]io.Reader
	Tag           string
	DKIM          *bool // this doesn't seem to have an associated domain API endpoint.
	DeliveryTime  *time.Time
	TestMode      *bool
	Tracking      *bool
	OpenTracking  *bool
	ClickTracking TrackOption
	RequireTLS    *bool
	SkipVerify    *bool
	Headers       textproto.MIMEHeader // keys will be canonicalized.
	Vars          map[string]string
	recipvars
}

func (m *Message) sum(l []string, s []string) error {
	sum := len(m.to) + len(m.cc) + len(m.bcc) + len(l) - len(s)
	if sum > 1000 {
		return fmt.Errorf("message would have %v receipients", sum)
	}
	return nil
}

// To sets the "to" recipients of Message, returning the previous list or
// error if total > 1000.
func (m *Message) To(to ...string) ([]string, error) {
	if err := m.sum(to, m.to); err != nil {
		return nil, err
	}
	was := m.to
	m.to = make([]string, len(to))
	copy(m.to, to)
	return was, nil
}

// CC sets the "cc" recipients of Message, returning the previous list or
// error if total > 1000.
func (m *Message) CC(cc ...string) ([]string, error) {
	if err := m.sum(cc, m.cc); err != nil {
		return nil, err
	}
	was := m.cc
	m.cc = make([]string, len(cc))
	copy(m.cc, cc)
	return was, nil
}

// BCC sets the "bcc" recipients of Message, returning the previous list or
// error if total > 1000.
func (m *Message) BCC(bcc ...string) ([]string, error) {
	if err := m.sum(bcc, m.bcc); err != nil {
		return nil, err
	}
	was := m.bcc
	m.cc = make([]string, len(bcc))
	copy(m.bcc, bcc)
	return was, nil
}

// Response is a response to a Send request.
type Response struct {
	ID string `json:"id"`
}

func bs(b bool) string {
	if b == true {
		return `yes`
	}
	return `no`
}

func (m *Message) textfields(f func(string, string) error) error {
	var err error
	wf := func(k, v string) error {
		if err != nil {
			return err
		}
		err = f(k, v)
		return err
	}
	bp := func(name string, b *bool) {
		if b == nil {
			return
		}
		wf(name, bs(*b))
	}
	wf("from", m.from)
	wf("subject", m.Subject)
	wf("html", m.HTML)
	bp("o:dkim", m.DKIM)
	bp("o:testmode", m.TestMode)
	bp("o:tracking", m.Tracking)
	bp("o:tracking-opens", m.OpenTracking)
	bp("o:require-tls", m.RequireTLS)
	bp("o:skip-verification", m.SkipVerify)
	if m.ClickTracking != nil {
		wf("o:tracking-clicks", fmt.Sprintf(`%s`, m.ClickTracking))
	}
	if len(m.Text) > 0 {
		wf("text", m.Text)
	}
	if len(m.Tag) > 0 {
		wf("o:tag", m.Tag)
	}
	if m.DeliveryTime != nil {
		wf("o:deliverytime", m.DeliveryTime.Format(time.RFC1123))
	}
	for k, v := range m.Vars {
		key := fmt.Sprintf(`v:%s`, k)
		wf(key, v)
	}
	if len(m.recipvars) > 0 {
		b, err := json.Marshal(m.recipvars)
		if err != nil {
			return err
		}
		wf(`recipient-variables`, string(b))
	}
	return err
}

// Send sends m, buffering in memory.
func (m *Message) Send(c client.Caller) (*Response, error) {
	buf := new(bytes.Buffer)
	return m.send(buf, c)
}

// TmpSend sends m, buffering to disk.
func (m *Message) TmpSend(c client.Caller) (*Response, error) {
	tf, err := ioutil.TempFile("", "mailgun-")
	if err != nil {
		return nil, err
	}
	os.Remove(tf.Name())
	defer tf.Close()
	return m.send(tf, c)
}

func (m *Message) send(rw io.ReadWriter, c client.Caller) (*Response, error) {
	var err error
	w := multipart.NewWriter(rw)
	wf := w.WriteField
	if err = m.textfields(wf); err != nil {
		return nil, err
	}
	for _, t := range m.to {
		if err = wf(`to`, t); err != nil {
			return nil, err
		}
	}
	for _, c := range m.cc {
		if err = wf(`cc`, c); err != nil {
			return nil, err
		}
	}
	for _, b := range m.bcc {
		if err = wf(`bcc`, b); err != nil {
			return nil, err
		}
	}
	for k, v := range m.Headers {
		key := fmt.Sprintf(`h:%s`, textproto.CanonicalMIMEHeaderKey(k))
		for _, v := range v {
			if err = wf(key, v); err != nil {
				return nil, err
			}
		}
	}
	for k, r := range m.Attachments {
		w, err := w.CreateFormFile("attachment", k)
		if err != nil {
			return nil, err
		}
		if _, err = io.Copy(w, r); err != nil {
			return nil, err
		}
	}
	for k, r := range m.Inlines {
		w, err := w.CreateFormFile("inline", k)
		if err != nil {
			return nil, err
		}
		if _, err = io.Copy(w, r); err != nil {
			return nil, err
		}
	}
	if err = w.Close(); err != nil {
		return nil, err
	}
	if s, ok := rw.(io.Seeker); ok {
		if _, err = s.Seek(0, 0); err != nil {
			return nil, err
		}
	}
	formdata := fmt.Sprintf(`multipart/form-data; boundary=%s`, w.Boundary())
	req := c.Post(`messages`)
	req.Header().Set("Content-Type", formdata)
	var re *Response
	return re, req.Payload(rw).Decode(&re)
}

// New returns a *Message with from, subject, and html content.
// If len(to) > 1000 an error is returned. Vars, Headers, Attachments
// and Inlines maps are initialized.
func New(from, subject, html string, to ...string) (*Message, error) {
	switch {
	case len(to) > 1000:
		return nil, fmt.Errorf("new: >1000 receipients")
	}
	return &Message{from: from, Subject: subject, HTML: html, to: to, Attachments: make(map[string]io.Reader), Inlines: make(map[string]io.Reader), Headers: make(textproto.MIMEHeader), Vars: make(map[string]string)}, nil
}

// BUG(j7b): Message.Vars is a little problematic. The relevant
// documentation at https://documentation.mailgun.com/en/latest/user_manual.html#attaching-data-to-messages
// is at best unclear. Empirically JSON-formatted strings
// with an ASCII character set appear to work well,
// callers may want to test other Message.Vars payload for potential astonishments.

// BUG(j7b): In general the endpoint does an OK job inferring
// MIME headers from Inline and Attachment payload. A map was chosen because
// although the endpoint may support attachments with duplicate filenames (and
// IDs for attachments) difficulties may present themselves with duplicates
// with the same content disposition.
