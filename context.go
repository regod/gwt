package gwt

import (
	"encoding/json"
	"net/http"
	"net/url"
	"sync"
)

// Context contains the context of one request
type Context struct {
	app      *Application
	request  *http.Request
	writer   http.ResponseWriter
	param    sync.Map
	query    url.Values
	postform url.Values
	store    sync.Map
}

// NewContext create a `Context` instance
func NewContext(app *Application, w http.ResponseWriter, r *http.Request) *Context {
	r.ParseForm()
	return &Context{
		app:      app,
		request:  r,
		writer:   w,
		query:    r.URL.Query(),
		postform: r.PostForm,
	}
}

func (ctx *Context) Query() url.Values {
	return ctx.query
}

func (ctx *Context) PostForm() url.Values {
	return ctx.postform
}

// SetParam set path param value
func (ctx *Context) SetParam(key string, value string) {
	ctx.param.Store(key, value)
}

// GetParam get param by `key`
func (ctx *Context) GetParam(key string) string {
	value, _ := ctx.param.Load(key)
	return value.(string)
}

// SetStore set custom variable to Context
func (ctx *Context) SetStore(key string, val interface{}) {
	ctx.store.Store(key, val)
}

// GetStore get from Context
func (ctx *Context) GetStore(key string) (val interface{}) {
	val, _ = ctx.store.Load(key)
	return val
}

// Writer return ResponseWriter
func (ctx *Context) Writer() http.ResponseWriter {
	return ctx.writer
}

// Request return http.Request
func (ctx *Context) Request() *http.Request {
	return ctx.request
}

// SetHeader set response header
func (ctx *Context) SetHeader(key, value string) {
	header := ctx.Writer().Header()
	header.Set(key, value)
}

func (ctx *Context) setContentType(contenttype string) {
	ctx.SetHeader("Content-Type", contenttype)
}

// RespBase response with statusCode, data string, content type
func (ctx *Context) RespBase(statusCode int, data string, contenttype string) error {
	ctx.setContentType(contenttype)
	ctx.writer.WriteHeader(statusCode)
	_, err := ctx.Writer().Write([]byte(data))
	return err
}

// RespText simple text response
func (ctx *Context) RespText(statusCode int, data string) error {
	return ctx.RespBase(statusCode, data, "text/plain")
}

// RespHtml html response
func (ctx *Context) RespHtml(statusCode int, data string) error {
	return ctx.RespBase(statusCode, data, "text/html")
}

// RespJson json response
func (ctx *Context) RespJson(statusCode int, data interface{}) error {
	enc := json.NewEncoder(ctx.writer)
	ctx.setContentType("application/json; charset=UTF-8")
	ctx.writer.WriteHeader(statusCode)
	return enc.Encode(data)
}
