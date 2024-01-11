# Routing

## Function Signatures
``` go
func (h *Harmony) Get(path string, handler HandlerFunc, middlewares ...MiddlewareFunc)
func (h *Harmony) Post(path string, handler HandlerFunc, middlewares ...MiddlewareFunc)
func (h *Harmony) Put(path string, handler HandlerFunc, middlewares ...MiddlewareFunc)
func (h *Harmony) Delete(path string, handler HandlerFunc, middlewares ...MiddlewareFunc)
func (h *Harmony) Patch(path string, handler HandlerFunc, middlewares ...MiddlewareFunc)
func (h *Harmony) Options(path string, handler HandlerFunc, middlewares ...MiddlewareFunc)
func (h *Harmony) Head(path string, handler HandlerFunc, middlewares ...MiddlewareFunc)
func (h *Harmony) Trace(path string, handler HandlerFunc, middlewares ...MiddlewareFunc)
```

## Simple Examples
### Get
``` go
app.Get("/", func(ctx harmony.Context) error {
    return ctx.String(http.StatusOK, "I'm a GET request")
})
```
### Post
``` go
app.Post("/", func(ctx harmony.Context) error {
    return ctx.String(http.StatusOK, "I'm a POST request")
})
```

## Grouping
### Function Signatures
``` go
func (h *Harmony) Group(prefix string, middlewares ...MiddlewareFunc) *Group
```
### Example
``` go

api := app.Group("/api", authMiddleware)
api.Get("/users", func(ctx harmony.Context) error {
    return ctx.String(http.StatusOK, "I'm a GET request to /api/users")
})
```

## Apply Middlewares
### Function Signature
``` go
func (h *Harmony) Use(middlewares ...MiddlewareFunc)
```
### Example
``` go
// Apply one middleware to all routes
app.Use(harmony.Gzip())

// Apply multiple middlewares to all routes
app.Use(
    harmony.Gzip(),
    harmony.Logger(),
    myCustomMiddleware(),
)
```