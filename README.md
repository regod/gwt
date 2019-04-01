gwt - Golang Web Tool

## example 
```go
package main

import (
	"github.com/regod/gwt"
        "github.com/regod/gwt/middleware"
)

func main() {
	app := gwt.New()

        app.AddRoute("/hello/", hello, nil)
        app.AddRoute("/ping/:id/:name", ping, []gwt.MiddlewareFunc{middleware.BasicAuth})
        app.Run(":9001")
}

func ping(ctx *gwt.Context) error {
        return ctx.RespJson(200, map[string]string{"ping": "pong"})
}

func hello(ctx *gwt.Context) error {
    return ctx.RespText(200, "hello world")
}
```
