// Package client is a low-level implementation
// of a mailgun REST client.
/*
This package is generally not intended for package consumers.
*/
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"path"
	"reflect"
	"strings"
)

// DEBUG enables debugging - requests logged, response copied to stderr.
var DEBUG = false

// Pager is implemented by pager types.
type Pager interface {
	SetCaller(Caller)
}

func apierr(res *http.Response) error {
	if res.StatusCode > 299 {
		defer res.Body.Close()
		buf := new(bytes.Buffer)
		io.CopyN(buf, res.Body, 80)
		return Error{
			Status:     res.Status,
			StatusCode: res.StatusCode,
			Stringer:   buf,
		}
	}
	return nil
}

// Error is an HTTP error from the endpoint.
type Error struct {
	Status     string
	StatusCode int
	fmt.Stringer
}

// Err returns *Error if e is
// Error, nil otherwise.
func Err(e error) *Error {
	switch t := e.(type) {
	case Error:
		return &t
	case *Error:
		return t
	}
	return nil
}

// Error implements error.
func (e Error) Error() string {
	return fmt.Sprintf("(%v) %s", e.StatusCode, e.Status)
}

func postHeader(h http.Header, boundary string) {
	h.Set("Content-Type", fmt.Sprintf(`multipart/form-data; boundary=%s`, boundary))
}

// Requester has methods that return Requests.
type Requester struct {
	Endpoint  string
	APIKey    string
	APIDomain string
	*http.Client
}

var _ = Caller(&Requester{})

// New turns a new Requester.
func New(endpoint, apikey, apidomain string) *Requester {
	return &Requester{Endpoint: endpoint, APIKey: apikey, APIDomain: apidomain}
}

// Domain returns the API domain.
func (r *Requester) Domain() string {
	return r.APIDomain
}

// Key returns the API key.
func (r *Requester) Key() string {
	return r.APIKey
}

// HTTPClient returns the *http.Client associated with this Requester.
func (r *Requester) HTTPClient() *http.Client {
	if r.Client != nil {
		return r.Client
	}
	return http.DefaultClient
}

// Request is a client request.
type Request struct {
	client   Caller
	e        error
	endpoint string
	method   string
	header   http.Header
	query    url.Values
	form     url.Values
	payload  io.Reader
}

func (r *Request) err(e error) error {
	if r.e == nil {
		r.e = e
	}
	return e
}

func (r *Request) done() {
	if r.e == nil {
		r.e = fmt.Errorf("request: done")
	}
}

// Header returns the http.Header for this Request.
func (r *Request) Header() http.Header {
	if r.header == nil {
		r.header = make(http.Header)
	}
	return r.header
}

// SetHeader sets a header, returning Request.
func (r *Request) SetHeader(k, v string) *Request {
	r.Header().Set(k, v)
	return r
}

// AddHeader adds a header, returning Request.
func (r *Request) AddHeader(k, v string) *Request {
	r.Header().Add(k, v)
	return r
}

// Query returns the url.Values for this Request's query string.
func (r *Request) Query() url.Values {
	if r.query == nil {
		r.query = make(url.Values)
	}
	return r.query
}

// SetQuery sets k,v to query string.
func (r *Request) SetQuery(k, v string) *Request {
	r.Query().Set(k, v)
	return r
}

// AddQuery adds k,v to query string.
func (r *Request) AddQuery(k, v string) *Request {
	r.Query().Add(k, v)
	return r
}

// Form returns the url.Values for this Request's form data.
func (r *Request) Form() url.Values {
	if r.form == nil {
		r.form = make(url.Values)
	}
	return r.form
}

// SetForm sets k,v to form values.
func (r *Request) SetForm(k, v string) *Request {
	r.Form().Set(k, v)
	return r
}

// AddForm adds k,v to form values.
func (r *Request) AddForm(k, v string) *Request {
	r.Form().Add(k, v)
	return r
}

// Payload sets the io.Reader for this Request. Precludes
// Form method.
func (r *Request) Payload(reader io.Reader) *Request {
	r.payload = reader
	return r
}

// Err performs this request, returning any error.
func (r *Request) Err() error {
	res, err := r.Do()
	if err != nil {
		return r.err(err)
	}
	return res.Body.Close()
}

func (r *Request) setpager(i interface{}) {
	v := reflect.ValueOf(i)
	switch v.Kind() {
	case reflect.Ptr:
		v = v.Elem()
	default:
		return
	}
	if sc, ok := v.Interface().(Pager); ok {
		sc.SetCaller(r.client)
	}
}

// Decode decodes the result of performing r to i.
func (r *Request) Decode(i interface{}) error {
	res, err := r.Do()
	if err != nil {
		return r.err(err)

	}
	defer res.Body.Close()
	if err = json.NewDecoder(res.Body).Decode(i); err != nil {
		return r.err(err)
	}
	r.setpager(i)
	return nil
}

// Do performs this Request.
func (r *Request) Do() (*http.Response, error) {
	defer r.done()
	if r.e != nil {
		return nil, r.e
	}
	client := r.client.HTTPClient()
	if len(r.query) > 0 {
		r.endpoint += fmt.Sprintf(`?%s`, r.query.Encode())
	}
	if len(r.form) > 0 {
		buf := new(bytes.Buffer)
		w := multipart.NewWriter(buf)
		for k, v := range r.form {
			for _, v := range v {
				if err := w.WriteField(k, v); err != nil {
					return nil, r.err(err)
				}
			}
		}
		if err := w.Close(); err != nil {
			return nil, r.err(err)
		}
		if r.header == nil {
			r.header = make(http.Header)
		}
		postHeader(r.header, w.Boundary())
		r.payload = buf
	}
	if DEBUG {
		log.Println("method", r.method)
		log.Println("endpoint", r.endpoint)
		log.Printf("payload: %T", r.payload)
	}
	req, err := http.NewRequest(r.method, r.endpoint, r.payload)
	if err != nil {
		return nil, r.err(err)
	}
	for k, v := range r.header {
		req.Header[k] = v
	}
	req.SetBasicAuth("api", r.client.Key())
	res, err := client.Do(req)
	if err != nil {
		return nil, r.err(err)
	}
	if err = apierr(res); err != nil {
		return nil, r.err(err)
	}
	if DEBUG {
		defer res.Body.Close()
		buf := new(bytes.Buffer)
		if _, err = io.Copy(buf, res.Body); err != nil {
			panic(err)
		}
		log.Println(buf.String())
		res.Body = ioutil.NopCloser(buf)
	}
	return res, nil
}

func join(uri []string) string {
	pth := path.Join(uri...)
	if strings.HasPrefix(pth, `https:/`) {
		pth = strings.Replace(pth, `https:/`, `https://`, 1)
	}
	return pth
}

func (r *Requester) request(method string, uri ...string) *Request {
	req := new(Request)
	if len(uri) == 0 {
		req.e = fmt.Errorf("request: uri not supplied")
		return req
	}
	req.method = method
	req.client = r
	pth := join(uri)
	switch {
	case strings.HasPrefix(pth, `https:/`):
		req.endpoint = pth
	case pth[0] == '/':
		req.endpoint = r.Endpoint + pth[1:]
	default:
		req.endpoint = r.Endpoint + path.Join(r.APIDomain, pth)
	}
	return req
}

// Get returns a GET Request.
func (r *Requester) Get(uri ...string) *Request {
	return r.request(http.MethodGet, uri...)
}

// Post returns a POST Request.
func (r *Requester) Post(uri ...string) *Request {
	return r.request(http.MethodPost, uri...)
}

// Put returns a PUT Request.
func (r *Requester) Put(uri ...string) *Request {
	return r.request(http.MethodPut, uri...)
}

// Delete returns a DELETE Request.
func (r *Requester) Delete(uri ...string) *Request {
	return r.request(http.MethodDelete, uri...)
}

// Caller is the interface to a Requester.
type Caller interface {
	HTTPClient() *http.Client
	Get(uri ...string) *Request
	Post(uri ...string) *Request
	Put(uri ...string) *Request
	Delete(uri ...string) *Request
	Domain() string
	Key() string
}
