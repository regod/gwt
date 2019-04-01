package gwt

import (
	"net"
	"net/http"
)

type (
	// Application main instance
	Application struct {
		server      *http.Server
		router      *Router
		ctx         *Context
		middlewares []MiddlewareFunc
	}

	// HandlerFunc function to serve request
	HandlerFunc func(*Context) error

	// MiddlewareFunc function for middleware
	MiddlewareFunc func(HandlerFunc) HandlerFunc
)

func New() *Application {
	app := &Application{
		server: new(http.Server),
		router: NewRouter(),
	}
	app.server.Handler = app
	return app
}

//SetMiddlewares set global middlewares
func (app *Application) SetMiddlewares(middlewares []MiddlewareFunc) {
	app.middlewares = middlewares
}

// AddRoute add route rule, specify middlewares of the rule by `middlewares`
// app.AddRoute("/user/list/", UserList, nil)
// app.AddRoute("/user/:id/", UserGet, nil)
func (app *Application) AddRoute(path string, h HandlerFunc, middlewares []MiddlewareFunc) {
	// chain middleware functions
	contains := append(app.middlewares, middlewares...)
	for i := len(contains) - 1; i >= 0; i-- {
		h = contains[i](h)
	}

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
