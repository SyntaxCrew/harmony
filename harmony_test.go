package harmony

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type (
	user struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
)

var (
	testUser = user{ID: 1, Name: "John Doe"}
)

const (
	userJSON = `{"id":1,"name":"John Doe"}`
)

func TestHarmony_Handler(t *testing.T) {
	app := New()
	app.Get("/", writeStringOKHandler())

	recCode, recBody := newRequest(http.MethodGet, "/", app)

	assert.Equal(t, http.StatusOK, recCode)
	assert.Equal(t, "OK", recBody)
}

func TestHarmony_Get(t *testing.T) {
	app := New()
	testMethod(t, http.MethodGet, "/", app)
}

func TestHarmony_Post(t *testing.T) {
	app := New()
	testMethod(t, http.MethodPost, "/", app)
}

func TestHarmony_Put(t *testing.T) {
	app := New()
	testMethod(t, http.MethodPut, "/", app)
}

func TestHarmony_Patch(t *testing.T) {
	app := New()
	testMethod(t, http.MethodPatch, "/", app)
}

func TestHarmony_Delete(t *testing.T) {
	app := New()
	testMethod(t, http.MethodDelete, "/", app)
}

func TestHarmony_Connect(t *testing.T) {
	app := New()
	testMethod(t, http.MethodConnect, "/", app)
}

func TestHarmony_Options(t *testing.T) {
	app := New()
	testMethod(t, http.MethodOptions, "/", app)
}

func TestHarmony_Trace(t *testing.T) {
	app := New()
	testMethod(t, http.MethodTrace, "/", app)
}

func TestHarmony_Head(t *testing.T) {
	app := New()
	testMethod(t, http.MethodHead, "/", app)
}

func TestHarmony_Middleware(t *testing.T) {
	app := New()
	buf := bytes.NewBuffer([]byte{})
	app.Use(func(next HandlerFunc) HandlerFunc {
		return func(ctx Context) error {
			buf.WriteString("1")
			return next(ctx)
		}
	})
	app.Use(func(next HandlerFunc) HandlerFunc {
		return func(ctx Context) error {
			buf.WriteString("1")
			return next(ctx)
		}
	})
	app.Use(func(next HandlerFunc) HandlerFunc {
		return func(ctx Context) error {
			buf.WriteString("2")
			return next(ctx)
		}
	})

	app.Get("/", writeStringOKHandler())

	recCode, recBody := newRequest(http.MethodGet, "/", app)
	assert.Equal(t, "112", buf.String())
	assert.Equal(t, http.StatusOK, recCode)
	assert.Equal(t, "OK", recBody)
}

func TestHarmony_Group(t *testing.T) {
	app := New()
	v1 := app.Group("/v1")
	v1.Get("/users", writeStringOKHandler())
	testMethod(t, http.MethodGet, "/v1/users", app)
}

func TestHarmony_MultiGroups(t *testing.T) {
	app := New()
	v1 := app.Group("/v1")
	v1.Get("/users", writeStringOKHandler())
	v2 := app.Group("/v2")
	v2.Get("/users", writeStringOKHandler())
	testMethod(t, http.MethodGet, "/v1/users", app)
	testMethod(t, http.MethodGet, "/v2/users", app)
}

func newRequest(method, path string, h *Harmony, body ...string) (int, string) {
	var b string
	if len(body) > 0 {
		b = body[0]
	}
	req := httptest.NewRequest(method, path, strings.NewReader(b))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, req)
	return rec.Code, rec.Body.String()
}

func testMethod(t *testing.T, method, path string, h *Harmony) {
	h.add(method, path, writeStringOKHandler())
	recCode, recBody := newRequest(method, path, h)
	assert.Equal(t, http.StatusOK, recCode)
	assert.Equal(t, "OK", recBody)
}

func writeStringOKHandler() HandlerFunc {
	return func(ctx Context) error {
		return ctx.String(http.StatusOK, "OK")
	}
}
