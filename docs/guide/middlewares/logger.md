# Logger

## Usage
``` go
app.Use(middleware.Logger())
```

## Custom Config
``` go
type 	LoggerConfig struct {
    // Skipper defines a function to skip middleware.
    Skipper Skipper

    // Format defines the logging format with defined variables.
    // Optional. Default value {remote_ip} - {host} "{method} {path} {protocol}" {status} {latency}
    // possible variables:
    // - {remote_ip}
    // - {host}
    // - {method}
    // - {path}
    // - {protocol}
    // - {status}
    // - {latency}
    Format string
}
```
### Example
``` go
app.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
    Skipper: middleware.DefaultSkipper,
    // GET /users HTTP/1.1 200 10µs
    Format: "{method} {path} {protocol} {status} {latency}",
}))
```