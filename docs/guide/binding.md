# Binding

## Bind
You only use [Bind](/guide/context#bind) when you want to bind the request body, path parameters, or query string parameters to a struct by using the struct tags to define the binding rules.
### JSON
``` go
// POST /users with JSON body: {"name": "John Doe"}
type User struct {
    Name string `json:"name"`
}
```

### Path Parameters
``` go
// GET /users/:username
type User struct {
    Username string `path:"username"`
}
```

### Query String Parameters
``` go
// GET /users?limit=10
type User struct {
    Limit int `query:"limit"`
}
```