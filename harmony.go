package harmony

import (
	gocontext "context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

const (
	banner = `
	  _   _
	 | | | | __ _ _ __ _ __ ___   ___  _ __  _   _
	 | |_| |/ _` + "`" + ` | '__| '_ ` + "`" + ` _ \ / _ \| '_ \| | | |
	 |  _  | (_| | |  | | | | | | (_) | | | | |_| |
	 |_| |_|\__,_|_|  |_| |_| |_|\___/|_| |_|\__, |
											 |___/
`
	defaultPort = 8080
)

const (
	// HeaderContentType is the header key for Content-Type.
	HeaderContentType = "Content-Type"
	// HeaderVary is the header key for Vary.
	HeaderVary = "Vary"
	// HeaderAcceptEncoding is the header key for Accept-Encoding.
	HeaderAcceptEncoding = "Accept-Encoding"
	// HeaderContentLength is the header key for Content-Length.
	HeaderContentLength = "Content-Length"
	// HeaderContentEncoding is the header key for Content-Encoding.
	HeaderContentEncoding = "Content-Encoding"
)

const (
	charsetUTF8 = "charset=utf-8"
)

const (
	// MIMEApplicationJSON is the MIME type for JSON.
	MIMEApplicationJSON = "application/json"
	// MIMEApplicationJSONCharsetUTF8 is the MIME type for JSON with charset=utf-8.
	MIMEApplicationJSONCharsetUTF8 = MIMEApplicationJSON + "; " + charsetUTF8
	// MIMETextPlain is the MIME type for plain text.
	MIMETextPlain = "text/plain"
	// MIMETextPlainCharsetUTF8 is the MIME type for plain text with charset=utf-8.
	MIMETextPlainCharsetUTF8 = MIMETextPlain + "; " + charsetUTF8
)

type (
	// Harmony is the interface for Harmony.
	Harmony struct {
		// gmux is the underlying router used by Harmony.
		gmux *mux.Router

		// middlewares is the list of middlewares used by Harmony.
		middlewares []mux.MiddlewareFunc

		// ctxPool is a pool of Context.
		ctxPool sync.Pool

		// srv is the underlying http.Server used by Harmony.
		srv *http.Server

		// group is the underlying group of Harmony.
		group map[string]*Harmony

		// binderPool is a pool of Binder.
		binderPool sync.Pool
	}

	// HandlerFunc is the function signature used by all Harmony handlers.
	HandlerFunc func(ctx Context) error

	// MiddlewareFunc is the function signature used by all Harmony middlewares.
	MiddlewareFunc func(next HandlerFunc) HandlerFunc

	// Map is a shortcut for map[string]any.
	Map map[string]any

	// HTTPError is the error returned by Harmony.
	HTTPError struct {
		Code    int
		Message string
	}
)

// New returns a new instance of Harmony.
func New() *Harmony {
	return &Harmony{
		gmux:  mux.NewRouter(),
		group: make(map[string]*Harmony),
	}
}

// ListenAndServe starts the server.
func (h *Harmony) ListenAndServe(port ...int) error {
	p := defaultPort
	if len(port) > 0 {
		p = port[0]
	}

	h.gmux.Use(h.middlewares...)

	h.srv = &http.Server{
		Addr:         "0.0.0.0:" + strconv.Itoa(p),
		WriteTimeout: 60 * time.Second,
		ReadTimeout:  60 * time.Second,
		IdleTimeout:  60 * time.Second,
		Handler:      h.gmux,
	}

	errCh := make(chan error, 1)

	go func() {
		log.Printf("%s\nharmony: server is listening and serving at :%d", banner, port)
		if err := h.srv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	return <-errCh
}

// ServeHTTP implements http.Handler.
func (h *Harmony) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx, ok := h.ctxPool.Get().(Context)
	if !ok {
		bdr, ok := h.binderPool.Get().(Binder)
		if !ok {
			bdr = newBinder()
		}
		ctx = NewContext(w, r, bdr)
	}
	defer func() {
		ctx.reset()
		h.ctxPool.Put(ctx)
	}()

	ctx.setRequest(r)
	ctx.SetResponseWriter(w)

	h.gmux.Use(h.middlewares...)
	h.gmux.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
}

// GracefulShutdown waits for SIGINT and gracefully shutdown the server.
func (h *Harmony) GracefulShutdown() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	ctx, cancel := gocontext.WithTimeout(gocontext.Background(), 10*time.Second)
	defer cancel()

	if err := h.srv.Shutdown(ctx); err != nil {
		return err
	}

	log.Println("harmony: gracefully shutdown")
	os.Exit(1)
	return nil
}

// Use adds a middleware to Harmony.
func (h *Harmony) Use(middlewares ...MiddlewareFunc) {
	for _, m := range middlewares {
		h.middlewares = append(h.middlewares, h.applyMiddleware(m))
	}
}

// Group creates a new Harmony group.
func (h *Harmony) Group(path string, middlewares ...MiddlewareFunc) *Group {
	return newGroup(path, h, middlewares...)
}

// Get adds a GET route to Harmony.
func (h *Harmony) Get(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	h.add(http.MethodGet, path, handlerFunc, middlewares...)
}

// Post adds a POST route to Harmony.
func (h *Harmony) Post(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	h.add(http.MethodPost, path, handlerFunc, middlewares...)
}

// Put adds a PUT route to Harmony.
func (h *Harmony) Put(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	h.add(http.MethodPut, path, handlerFunc, middlewares...)
}

// Patch adds a PATCH route to Harmony.
func (h *Harmony) Patch(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	h.add(http.MethodPatch, path, handlerFunc, middlewares...)
}

// Delete adds a DELETE route to Harmony.
func (h *Harmony) Delete(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	h.add(http.MethodDelete, path, handlerFunc, middlewares...)
}

// Connect adds a CONNECT route to Harmony.
func (h *Harmony) Connect(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	h.add(http.MethodConnect, path, handlerFunc, middlewares...)
}

// Options adds an OPTIONS route to Harmony.
func (h *Harmony) Options(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	h.add(http.MethodOptions, path, handlerFunc, middlewares...)
}

// Head adds a HEAD route to Harmony.
func (h *Harmony) Head(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	h.add(http.MethodHead, path, handlerFunc, middlewares...)
}

// Trace adds a TRACE route to Harmony.
func (h *Harmony) Trace(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	h.add(http.MethodTrace, path, handlerFunc, middlewares...)
}

// NewHTTPError returns a new HTTP error.
func NewHTTPError(code int, message string) *HTTPError {
	return &HTTPError{Code: code, Message: message}
}

func (h *Harmony) add(method, path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	h.gmux.
		HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			ctx, ok := h.ctxPool.Get().(Context)
			if !ok {
				bdr, ok := h.binderPool.Get().(Binder)
				if !ok {
					bdr = newBinder()
				}
				ctx = NewContext(w, r, bdr)
			}
			defer func() {
				ctx.reset()
				h.ctxPool.Put(ctx)
			}()

			ctx.SetResponseWriter(w)
			ctx.setRequest(r)

			for _, middleware := range middlewares {
				handlerFunc = middleware(handlerFunc)
			}

			_ = handlerFunc(ctx)
		}).
		Methods(method)
}

func (h *Harmony) applyMiddleware(middleware MiddlewareFunc) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, ok := h.ctxPool.Get().(Context)
			if !ok {
				bdr, ok := h.binderPool.Get().(Binder)
				if !ok {
					bdr = newBinder()
				}
				ctx = NewContext(w, r, bdr)
			}
			defer func() {
				ctx.reset()
				h.ctxPool.Put(ctx)
			}()

			ctx.SetResponseWriter(w)
			ctx.setRequest(r)

			_ = middleware(func(ctx Context) error {
				next.ServeHTTP(ctx.ResponseWriter(), ctx.Request())
				return nil
			})(ctx)
		})
	}
}

func (e *HTTPError) Error() string {
	return fmt.Sprintf("harmony: code=%d, message=%s", e.Code, e.Message)
}
