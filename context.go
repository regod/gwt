package gwt

import (
	"net/http"
	"net/url"
)

// Context contains the context of one request
type Context struct {
	app      *Application
	request  *http.Request
	writer   http.ResponseWriter
	param    map[string]string
	query    url.Values
	postform url.Values
}

// NewContext create a `Context` instance
func NewContext(app *Application, w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		app:     app,
		request: r,
		writer:  w,
		param:   make(map[string]string),
	}
}

// SetParam set path param value
func (ctx *Context) SetParam(key string, value string) {
	ctx.param[key] = value
}

// GetParam get param by `key`
func (ctx *Context) GetParam(key string) string {
	return ctx.param[key]
}

// Writer return ResponseWriter
func (ctx *Context) Writer() http.ResponseWriter {
	return ctx.writer
}
