package gwt

import (
	"net/http"
	"testing"
)

func TestRouterBase(t *testing.T) {
	app := New()
	path := "/ab/c/"
	app.AddRoute("GET", path, func(ctx *Context) error {
		ctx.SetStore("path", path)
		return nil
	}, nil)

	r, _ := http.NewRequest("GET", path, nil)
	ctx := NewContext(app, nil, r)
	h := app.router.Detect(path, ctx)
	h(ctx)
	actual := ctx.GetStore("path").(string)
	if path != actual {
		t.Errorf("actual: %s; expect: %s", actual, path)
	}
}

func TestRouterParam(t *testing.T) {
	app := New()
	path := "/ab/c/:id/"
	app.AddRoute("POST", path, func(ctx *Context) error {
		ctx.SetStore("path", path)
		return nil
	}, nil)

	rpath := "/ab/c/134/"
	r, _ := http.NewRequest("POST", rpath, nil)
	ctx := NewContext(app, nil, r)
	h := app.router.Detect(rpath, ctx)
	h(ctx)
	actual := ctx.GetParam("id")
	if "134" != actual {
		t.Errorf("actual: %s; expect: %s", actual, "134")
	}
}

func TestRouterAnyMethod(t *testing.T) {
	app := New()
	path := "/ab/c/"
	app.AddRoute("*", path, func(ctx *Context) error {
		ctx.SetStore("path", path)
		return nil
	}, nil)

	for _, method := range []string{"GET", "POST", "PUT", "DELETE", "HEAD"} {
		r, _ := http.NewRequest(method, path, nil)
		ctx := NewContext(app, nil, r)
		h := app.router.Detect(path, ctx)
		h(ctx)
		actual := ctx.GetStore("path").(string)
		if path != actual {
			t.Errorf("actual: %s; expect: %s", actual, path)
		}
	}
}

func TestRouterWrongMethod(t *testing.T) {
	app := New()
	path := "/ab/c/"
	app.AddRoute("GET", path, func(ctx *Context) error {
		ctx.SetStore("path", path)
		return nil
	}, nil)

	r, _ := http.NewRequest("POST", path, nil)
	ctx := NewContext(app, nil, r)
	h := app.router.Detect(path, ctx)
	if h != nil {
		t.Errorf("handler should be nil")
	}
}
