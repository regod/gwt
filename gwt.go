package gwt

import (
	"net"
	"net/http"
)

type Application struct {
	server *http.Server
	router *Router
	//ctx *Context
}

func New() *Application {
	app := &Application{
		server: new(http.Server),
		router: NewRouter(),
	}
	app.server.Handler = app
	return app
}

func (app *Application) AddRoute(path string, h http.HandlerFunc) {
	app.router.Register(path, h)
}

func (app *Application) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// build req context
	// detect uri to handler
	path := r.URL.Path
	h := app.router.Detect(path)
	// execute handler
	h(w, r)
}

func (app *Application) Run(address string) error {
	app.server.Addr = address
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	return app.server.Serve(listener)
}
