# Context
`harmony.Context` is a wrapper around `http.ResponseWriter`, and `*http.Request` that provides a simple interface for writing HTTP handlers.

## Request
Returns the `*http.Request` object
### Function Signature
``` go
func (ctx *context) Request() *http.Request
```

## ResponseWriter
Returns the `http.ResponseWriter` object
### Function Signature
``` go
func (ctx *context) Response() http.ResponseWriter
```

## Bind
Binds the request body into dest.
### Function Signature
``` go
func (ctx *context) Bind(dest any) error
```
### Example
``` go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

app.Post("/user", func(ctx *harmony.Context) error {
    var user User
    err := ctx.Bind(&user)
    if err != nil {
        return ctx.JSON(http.StatusBadRequest, harmony.Map{"error": err.Error()})
    }

    // ...
})
```

## JSON
Writes the response in JSON format
### Function Signature
``` go
func (ctx *context) JSON(code int, body any) error
```
### Example
``` go
type User struct {
    Name string `json:"name"`
    Age  int    `json:"age"`
}

app.Get("/user", func(ctx *harmony.Context) error {
    user := User{
        Name: "John Doe",
        Age:  30,
    }
    return ctx.JSON(http.StatusOK, user)
})
```

## PathParams
Returns the path parameters
### Function Signature
``` go
func (ctx *context) PathParams() map[string]string
```
### Example
``` go
app.Get("/user/:username", func(ctx *harmony.Context) error {
    pathParams := ctx.PathParams()
    username := pathParams["username"]

    // ...
})
```

## PathParam
Returns the path parameters as a string
### Function Signature
``` go
func (ctx *context) PathParam(key string) string
```
### Example
``` go
app.Get("/user/:username", func(ctx *harmony.Context) error {
    username := ctx.PathParam("username")

    // ...
})
```

## PathParamInt
Returns the path parameters as an integer
### Function Signature
``` go
func (ctx *context) PathParamInt(key string) (int, error)
```
### Example
``` go
app.Get("/user/:id", func(ctx *harmony.Context) error {
    id, err := ctx.PathParamInt("id")

    // ...
})
```

## QueryString
Returns the query parameters as a string
### Function Signature
``` go
func (ctx *context) QueryString(key string, defaultValue ...string) string
```
### Example
``` go
// GET /users?status=active
app.Get("/user", func(ctx *harmony.Context) error {
    // get query string params
    status := ctx.QueryString("status")
    // get query string params with default value
    sort := ctx.QueryString("sort", "desc")

    // ...
})
```


## QueryInt
Returns the query parameters as an integer
### Function Signature
``` go
func (ctx *context) QueryInt(key string, defaultValue ...int) int
```
### Example
``` go
// GET /users?organization_id=1
app.Get("/user", func(ctx *harmony.Context) error {
// get query int params
organizationID := ctx.QueryInt("organization_id")
// get query int params with default value
limit := ctx.Query("limit", 10)

// ...
})
```

## QueryFloat64
Returns the query parameters as a float64
### Function Signature
``` go
func (ctx *context) QueryFloat64(key string, defaultValue ...float64) float64
```
### Example
``` go
// GET /users?latitude=1.2345
app.Get("/user", func(ctx *harmony.Context) error {
    // get query float64 params
    latitude := ctx.QueryFloat64("latitude")
    // get query float64 params with default value
    longitude := ctx.QueryFloat64("longitude", 1.2345)

    // ...
})
```

## QueryBool
Returns the query parameters as a boolean
### Function Signature
``` go
func (ctx *context) QueryBool(key string, defaultValue ...bool) bool
```
### Example
``` go
// GET /users?is_active=true
app.Get("/user", func(ctx *harmony.Context) error {
    // get query bool params
    isActive := ctx.QueryBool("is_active")
    // get query bool params with default value
    isKYCVerified := ctx.QueryBool("is_kyc_verified", true)

    // ...
})
```

## SendStatus
Sends an HTTP response with only the given status code
### Function Signature
``` go
func (ctx *context) SendStatus(code int) error
```
### Example
``` go
app.Get("/livez", func(ctx *harmony.Context) error {
    return ctx.SendStatus(http.StatusOK)
})
```

## String
Writes the given string to the response
### Function Signature
``` go
func (ctx *context) String(code int, body string) error
```
### Example
``` go
app.Get("/hello", func(ctx *harmony.Context) error {
    return ctx.String(http.StatusOK, "Welcome to Harmony!")
})
```

## Get
Returns the value of the given key in the context
### Function Signature
``` go
func (ctx *context) Get(key string) string
```
### Example
``` go
app.Get("/hello", func(ctx *harmony.Context) error {
    traceID := ctx.Get("traceID")

    // ...
})
```

## Set
Sets the value of the given key in the context
### Function Signature
``` go
func (ctx *context) Set(key string, value any)
```
### Example
``` go
app.Get("/hello", func(ctx *harmony.Context) error {
    ctx.Set("userID", 1)

    // ...
})
```