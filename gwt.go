package gwt

import (
	"net"
	"net/http"
)

type (
        // Application main instance
	Application struct {
		server *http.Server
		router *Router
		ctx    *Context
	}

        // HandlerFunc function to serve request
	HandlerFunc func(*Context) error
)

func New() *Application {
	app := &Application{
		server: new(http.Server),
		router: NewRouter(),
	}
	app.server.Handler = app
	return app
}

func (app *Application) AddRoute(path string, h HandlerFunc) {
	app.router.Register(path, h)
}

func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// build req context
	ctx := NewContext(app, w, r)
	// detect uri to handler
	path := r.URL.Path
	h := app.router.Detect(path, ctx)
	// execute handler
	h(ctx)
}

func (app *Application) Run(address string) error {
	app.server.Addr = address
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	return app.server.Serve(listener)
}
