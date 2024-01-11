package harmony

import "net/http"

type (
	// Group is the interface for Harmony's Group.
	Group struct {
		harmony     *Harmony
		prefix      string
		middlewares []MiddlewareFunc
	}
)

// Use adds middlewares to the group.
func (g *Group) Use(middlewares ...MiddlewareFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// Group creates a new Harmony subgroup in the current group
func (g *Group) Group(path string, middlewares ...MiddlewareFunc) *Group {
	return g.harmony.Group(g.prefix+path, middlewares...)
}

// Get adds a GET route to Harmony's Group.
func (g *Group) Get(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	g.add(http.MethodGet, path, handlerFunc, middlewares...)
}

// Post adds a POST route to Harmony.
func (g *Group) Post(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	g.add(http.MethodPost, path, handlerFunc, middlewares...)
}

// Put adds a PUT route to Harmony's Group.
func (g *Group) Put(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	g.add(http.MethodPut, path, handlerFunc, middlewares...)
}

// Patch adds a PATCH route to Harmony's Group.
func (g *Group) Patch(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	g.add(http.MethodPatch, path, handlerFunc, middlewares...)
}

// Delete adds a DELETE route to Harmony's Group.
func (g *Group) Delete(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	g.add(http.MethodDelete, path, handlerFunc, middlewares...)
}

// Connect adds a CONNECT route to Harmony's Group.
func (g *Group) Connect(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	g.add(http.MethodConnect, path, handlerFunc, middlewares...)
}

// Options adds an OPTIONS route to Harmony's Group.
func (g *Group) Options(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	g.add(http.MethodOptions, path, handlerFunc, middlewares...)
}

// Head adds a HEAD route to Harmony's Group.
func (g *Group) Head(path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	g.add(http.MethodHead, path, handlerFunc, middlewares...)
}

func (g *Group) add(method, path string, handlerFunc HandlerFunc, middlewares ...MiddlewareFunc) {
	g.harmony.add(method, g.prefix+path, handlerFunc, append(g.middlewares, middlewares...)...)
}

func newGroup(prefix string, harmony *Harmony, middlewares ...MiddlewareFunc) *Group {
	g := &Group{
		harmony:     harmony,
		prefix:      prefix,
		middlewares: middlewares,
	}
	g.Use(middlewares...)
	return g
}
