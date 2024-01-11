# Quick Start

## Installation
Download and install:
``` bash
go get -u github.com/SyntaxCrew/harmony
```

## Ping-Pong Server
1. Create `main.go`
``` go
package main

import (
    "log"
    "net/http"

    "github.com/SyntaxCrew/harmony"
)

func main() {
    app := harmony.New()
    app.Get("/ping", func(ctx harmony.Context) error {
        return ctx.String(http.StatusOK, "pong")
    })
    if err := app.ListenAndServe(); err != http.ErrServerClosed {
        log.Fatal(err)
    }
}
```
2. You can run the server with
``` bash
go run main.go
```
3. Open your browser and visit to `0.0.0.0:8080/ping`