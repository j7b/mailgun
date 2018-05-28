package main

import (
	"bufio"
	"bytes"
	"go/format"
	"log"
	"net/textproto"
	"os"
	"path/filepath"
	"strings"
)

func togo(name string) string {
	name = strings.Replace(name, `-`, ` `, -1)
	name = strings.Title(name)
	return strings.Replace(name, ` `, ``, -1)
}

func parse(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	structname := strings.Replace(filename, filepath.Ext(filename), ``, 1)
	buf := new(bytes.Buffer)
	bw := bufio.NewWriter(buf)
	w := textproto.NewWriter(bw)
	pl := w.PrintfLine
	pl("package types\n")
	pl(`import "net/textproto"`)
	pl(`var _ = textproto.Error{}`)
	pl(`type %s struct {`, strings.Title(structname))
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		txt := scanner.Text()
		parts := strings.Split(txt, "\t")
		if len(parts) != 2 {
			continue
		}
		js := parts[0]
		typ := `string`
		name := togo(js)
		if strings.Contains(js, `“`) {
			js = `-`
			name = `CustomVariables`
			typ = `map[string]interface{}`
		}
		switch js {
		case `event`, `timestamp`, `token`, `signature`:
			typ = ``
		}
		if strings.Contains(js, `attachment-x`) {
			js = `-`
			name = `Attachments`
			typ = `map[*textproto.MIMEHeader][]byte`
		}
		if strings.HasSuffix(name, `Url`) {
			l := len(name)
			name = name[:l-3] + `URL`
		}
		if strings.HasSuffix(name, `Id`) {
			l := len(name)
			name = name[:l-2] + `ID`
		}
		if strings.HasSuffix(name, `Ip`) {
			l := len(name)
			name = name[:l-2] + `IP`
		}
		if strings.HasSuffix(name, `Os`) {
			l := len(name)
			name = name[:l-2] + `OS`
		}
		comment := parts[1]
		pl("%s %s `json:\"%s\"` // %s", name, typ, js, comment)
	}
	pl("}\n")
	pl("func (%s %s) %s() %s {return %s}\n", structname[:1], strings.Title(structname), strings.Title(structname), strings.Title(structname), structname[:1])
	pl("type I%s interface {%s() %s}\n", strings.Title(structname), strings.Title(structname), strings.Title(structname))
	bw.Flush()
	src, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}
	f, err = os.Create(filepath.Join(`..`, structname+`.go`))
	if err != nil {
		return err
	}
	if _, err = f.Write(src); err != nil {
		return err
	}
	return f.Close()
}

func main() {
	for _, name := range []string{
		"bounce.txt",
		"click.txt",
		"complaint.txt",
		"delivered.txt",
		"drop.txt",
		"open.txt",
		"unsubscribe.txt",
	} {
		if err := parse(name); err != nil {
			log.Fatal(err)
		}
	}
}
