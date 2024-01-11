package harmony

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"sync"
)

type (
	Context interface {
		// Request returns the *http.Request object.
		Request() *http.Request

		// ResponseWriter returns the http.ResponseWriter object.
		ResponseWriter() http.ResponseWriter

		// Bind binds the request body into dest.
		Bind(dest any) error

		// JSON writes the response in JSON format.
		JSON(code int, body any) error

		// PathParams returns the path parameters of the request in map[string]string.
		PathParams() map[string]string

		// PathParam returns the path parameter of the request by key in string.
		PathParam(key string) string

		// PathParamInt returns the path parameter of the request by key in int.
		// If the value is not an integer, it will return an error.
		// Handle the error like strconv.Atoi.
		PathParamInt(key string) (int, error)

		// SetPathParam sets the path parameter of the request by key and value.
		SetPathParam(key, value string)

		// QueryString returns the query string of the request in string.
		QueryString(key string, defaultValue ...string) string

		// QueryInt returns the query parameter of the request by key in int.
		QueryInt(key string, defaultValue ...int) int

		// QueryFloat64 returns the query parameter of the request by key in float64.
		QueryFloat64(key string, defaultValue ...float64) float64

		// QueryBool returns the query parameter of the request by key in bool.
		QueryBool(key string, defaultValue ...bool) bool

		// SendStatus writes the response status code.
		SendStatus(code int) error

		// String writes the response in string format.
		String(code int, body string) error

		// Get returns the value in the context by key.
		Get(key string) any

		// Set sets the value in the context by key and value.
		Set(key string, value any)

		// reset resets the context.
		reset()

		// SetResponseWriter sets the http.ResponseWriter.
		SetResponseWriter(w http.ResponseWriter)

		// setRequest sets the *http.Request.
		setRequest(r *http.Request)
	}

	context struct {
		w     http.ResponseWriter
		r     *http.Request
		store Map
		lock  sync.RWMutex
		bdr   Binder
	}
)

func NewContext(w http.ResponseWriter, r *http.Request, binder Binder) Context {
	return &context{
		w:     w,
		r:     r,
		store: make(Map),
		bdr:   binder,
	}
}

func (c *context) Request() *http.Request {
	return c.r
}

func (c *context) ResponseWriter() http.ResponseWriter {
	return c.w
}

func (c *context) Bind(dest any) error {
	return c.bdr.Bind(c, dest)
}

func (c *context) JSON(code int, body any) error {
	c.w.Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	c.w.WriteHeader(code)
	return json.NewEncoder(c.w).Encode(body)
}

func (c *context) PathParams() map[string]string {
	return mux.Vars(c.r)
}

func (c *context) PathParam(key string) string {
	return mux.Vars(c.r)[key]
}

func (c *context) PathParamInt(key string) (int, error) {
	return strconv.Atoi(mux.Vars(c.r)[key])
}

func (c *context) SetPathParam(key, value string) {
	vars := make(map[string]string)
	if len(mux.Vars(c.r)) > 0 {
		vars = mux.Vars(c.r)
	}
	if _, ok := vars[key]; !ok {
		vars[key] = value
	}
	c.r = mux.SetURLVars(c.r, vars)
}

func (c *context) QueryString(key string, defaultValue ...string) string {
	qs := c.r.URL.Query().Get(key)
	if qs == "" && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return qs
}

func (c *context) QueryInt(key string, defaultValue ...int) int {
	i, err := strconv.Atoi(c.r.URL.Query().Get(key))
	if err != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return i
}

func (c *context) QueryFloat64(key string, defaultValue ...float64) float64 {
	f64, err := strconv.ParseFloat(c.r.URL.Query().Get(key), 64)
	if err != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return f64
}

func (c *context) QueryBool(key string, defaultValue ...bool) bool {
	b, err := strconv.ParseBool(c.r.URL.Query().Get(key))
	if err != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	return b
}

func (c *context) SendStatus(code int) error {
	c.w.WriteHeader(code)
	return nil
}

func (c *context) String(code int, s string) error {
	c.w.Header().Set(HeaderContentType, MIMETextPlainCharsetUTF8)
	c.w.WriteHeader(code)
	_, err := c.w.Write([]byte(s))
	return err
}

func (c *context) Get(key string) any {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.store[key]
}

func (c *context) Set(key string, value any) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.store == nil {
		c.store = make(Map)
	}
	c.store[key] = value
}

func (c *context) reset() {
	c.w = nil
	c.r = nil
	c.store = make(Map)
}

func (c *context) SetResponseWriter(w http.ResponseWriter) {
	c.w = w
}

func (c *context) setRequest(r *http.Request) {
	c.r = r
}
