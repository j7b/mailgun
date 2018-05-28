// Package types defines webhook types.
package types

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/textproto"
	"net/url"
	"strings"
)

func decodeform(reader *bufio.Reader, boundary string) (Type, error) {
	if len(boundary) == 0 {
		b, _ := reader.Peek(72)
		i := bytes.IndexByte(b, '\n')
		if i == -1 {
			return nil, fmt.Errorf("Decode: couldn't find form boundary")
		}
		boundary = strings.TrimSpace(string(b[2:i]))
	}
	m := make(url.Values)
	files := make(map[*textproto.MIMEHeader][]byte)
	r := multipart.NewReader(reader, boundary)
	part, err := r.NextPart()
	for err == nil {
		formname := part.FormName()
		filename := part.FileName()
		var b []byte
		b, err = ioutil.ReadAll(part)
		if err != nil {
			return nil, err
		}
		switch {
		case len(filename) > 0:
			files[&part.Header] = b
		case len(formname) > 0:
			m.Add(formname, string(b))
		default:
			return nil, fmt.Errorf("Decode: form part with no file/field")
		}
		part.Close()
		part, err = r.NextPart()
	}
	if err != io.EOF {
		return nil, err
	}
	return decode(m, files)
}

func decode(vals url.Values, files map[*textproto.MIMEHeader][]byte) (Type, error) {
	event := vals.Get("event")
	dv := func(i interface{}) (Type, error) {
		return decodevals(vals, files, i)
	}
	switch event {
	case `clicked`:
		return dv(&Click{})
	case `delivered`:
		return dv(&Delivered{})
	case `dropped`:
		return dv(&Drop{})
	case `bounced`:
		return dv(&Bounce{})
	case `opened`:
		return dv(&Open{})
	case `complained`:
		return dv(&Complaint{})
	case `unsubscribed`:
		return dv(&Unsubscribe{})
	}
	return nil, fmt.Errorf("Decode: unknown event %s", event)
}

func decodevars(r io.Reader) (Type, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	vals, err := url.ParseQuery(string(b))
	if err != nil {
		return nil, err
	}
	return decode(vals, nil)
}

// Decode decodes webhook payload from reader. If boundary
// is empty, guesses if x-form-urlencoded or multipart.
func Decode(reader io.Reader, boundary string) (Type, error) {
	r := bufio.NewReader(reader)
	if len(boundary) > 0 {
		return decodeform(r, boundary)
	}
	peek, err := r.Peek(2)
	if err != nil {
		return nil, err
	}
	switch string(peek) {
	case `--`:
		return decodeform(r, boundary)
	}
	return decodevars(r)
}
