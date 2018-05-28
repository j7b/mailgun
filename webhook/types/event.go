package types

import (
	"crypto/hmac"
	"crypto/sha256"
	"fmt"
	"net/textproto"
	"net/url"
	"reflect"
	"strings"
)

// Event is shared by all webhook types.
type Event string

// Type implements Type.
func (e Event) Type() string {
	return string(e)
}

// Timestamp is shared by all webhook types.
type Timestamp string

// TS implements Type.
func (t Timestamp) TS() string {
	return string(t)
}

// Token is shared by all webhook types.
type Token string

// Tok implements Type.
func (t Token) Tok() string {
	return string(t)
}

// Signature is shared by all webhook types.
type Signature string

// Sig implements Type.
func (s Signature) Sig() string {
	return string(s)
}

// Type is the event interface.
type Type interface {
	Type() string
	Sig() string
	TS() string
	Tok() string
}

// APIKey defines Key() method.
type APIKey interface {
	Key() string
}

// Validate signature on t with a.
func Validate(t Type, a APIKey) error {
	if t == nil {
		return fmt.Errorf(`Validate: Type is nil`)
	}
	if a == nil {
		return fmt.Errorf(`Validate: APIKey is nil`)
	}
	hash := hmac.New(sha256.New, []byte(a.Key()))
	hash.Write([]byte(t.TS()))
	hash.Write([]byte(t.Tok()))
	want := fmt.Sprintf(`%x`, hash.Sum(nil))
	if want == t.Sig() {
		return nil
	}
	return fmt.Errorf(`Validate: want %s got %s`, t.Sig(), want)
}

func decodevals(vals url.Values, files map[*textproto.MIMEHeader][]byte, i interface{}) (event Type, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
		}
	}()
	val := reflect.ValueOf(i)
	val = val.Elem()
	typ := val.Type()
	var cv reflect.Value
	var fv reflect.Value
	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if f.Name == `CustomVariables` {
			cv = val.Field(i)
			continue
		}
		if f.Name == `Attachments` {
			fv = val.Field(i)
		}
		fieldname := strings.TrimSpace(strings.Split(f.Tag.Get(`json`), `,`)[0])
		if len(fieldname) == 0 || fieldname == `-` {
			continue
		}
		v := vals.Get(fieldname)
		if len(v) > 0 {
			switch fieldname {
			case `event`:
				val.Field(i).Set(reflect.ValueOf(Event(v)))
			case `timestamp`:
				val.Field(i).Set(reflect.ValueOf(Timestamp(v)))
			case `token`:
				val.Field(i).Set(reflect.ValueOf(Token(v)))
			case `signature`:
				val.Field(i).Set(reflect.ValueOf(Signature(v)))
			default:
				val.Field(i).Set(reflect.ValueOf(v))
			}
		}
		delete(vals, fieldname)
	}
	if cv.IsValid() {
		m := make(map[string]interface{})
		for k := range vals {
			if v := vals.Get(k); len(v) > 0 {
				m[k] = v
			}
		}
		if len(m) > 0 {
			cv.Set(reflect.ValueOf(m))
		}
	}
	if fv.IsValid() && len(files) > 0 {
		fv.Set(reflect.ValueOf(files))
	}
	return val.Interface().(Type), nil
}
