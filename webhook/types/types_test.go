package types

import (
	"bytes"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/j7b/mailgun"
)

func TestValidate(t *testing.T) {
	c, err := mailgun.New(`key-fake`, `domain.fake`)
	if err != nil {
		t.Fatal(err)
	}
	tok := make([]byte, 32)
	rand.Read(tok)
	ts := fmt.Sprintf(`%v`, time.Now().Unix())
	sig := hmac.New(sha256.New, []byte(`key-fake`))
	sig.Write([]byte(ts))
	sig.Write([]byte(fmt.Sprintf(`%x`, tok)))
	e := Bounce{
		Signature: Signature(fmt.Sprintf(`%x`, sig.Sum(nil))),
		Timestamp: Timestamp(ts),
		Token:     Token(fmt.Sprintf(`%x`, tok)),
	}
	if err = Validate(e, c); err != nil {
		t.Fatal(err)
	}
}

func TestTypes(t *testing.T) {
	var table = map[string]string{
		`_payload/click/click.txt`:           `clicked`,
		`_payload/delivered/delivered.txt`:   `delivered`,
		`_payload/dropped/dropped.txt`:       `dropped`,
		`_payload/hardbounce/hardbounce.txt`: `bounced`,
		`_payload/open/open.txt`:             `opened`,
		`_payload/scomp/scomp.txt`:           `complained`,
		`_payload/unsub/unsub.txt`:           `unsubscribed`,
	}
	for k, v := range table {
		t.Run(v, func(t *testing.T) {
			f, err := os.Open(k)
			if err != nil {
				t.Fatal(err)
			}
			defer f.Close()
			typ, err := Decode(f, ``)
			if err != nil {
				t.Fatal(err)
			}
			if typ.Type() != v {
				t.Fatalf(`want %s got %s`, v, typ.Type())
			}
		})
	}
	_, err := Decode(bytes.NewBufferString(`event=bogus`), ``)
	if err == nil {
		t.Fatal(`decode bogus event should error`)
	}
	if _, err = Decode(bytes.NewBufferString(`---- - ------------herpderp`), ``); err == nil {
		t.Fatal(`decode broken multipart should error`)
	}
}
