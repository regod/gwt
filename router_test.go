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

func TestRouteSomePathDiffParam(t *testing.T) {
	app := New()
	path1 := "/ab/c/:id/"
	path2 := "/ab/c/:name/"
	app.AddRoute("GET", path1, func(ctx *Context) error {
		ctx.SetStore("path", path1)
		return nil
	}, nil)
	app.AddRoute("POST", path2, func(ctx *Context) error {
		ctx.SetStore("path", path2)
		return nil
	}, nil)

	rpath := "/ab/c/134/"
	r, _ := http.NewRequest("GET", rpath, nil)
	ctx := NewContext(app, nil, r)
	h := app.router.Detect(rpath, ctx)
	h(ctx)
	actual := ctx.GetParam("id")
	if "134" != actual {
		t.Errorf("actual: %s; expect: %s", actual, "134")
	}

	rpath = "/ab/c/134/"
	r, _ = http.NewRequest("POST", rpath, nil)
	ctx = NewContext(app, nil, r)
	h = app.router.Detect(rpath, ctx)
	h(ctx)
	actual = ctx.GetParam("name")
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

func TestRouteTrailingSlash(t *testing.T) {
	app := New()
	path1 := "/ab/c/"
	path2 := "/ab/c"
	app.AddRoute("GET", path1, func(ctx *Context) error {
		ctx.SetStore("trailingslash", true)
		return nil
	}, nil)
	app.AddRoute("GET", path2, func(ctx *Context) error {
		ctx.SetStore("trailingslash", false)
		return nil
	}, nil)

	r, _ := http.NewRequest("GET", path1, nil)
	ctx := NewContext(app, nil, r)
	h := app.router.Detect(path1, ctx)
	h(ctx)
	val1 := ctx.GetStore("trailingslash")

	r, _ = http.NewRequest("GET", path2, nil)
	ctx = NewContext(app, nil, r)
	h = app.router.Detect(path2, ctx)
	h(ctx)
	val2 := ctx.GetStore("trailingslash")

	if val1 == val2 {
		t.Errorf("val1: %s; val2: %s", val1, val2)
	}
}
