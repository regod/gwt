package middleware

import (
	"github.com/regod/gwt"
	"net/http"
)

func BasicAuth(handler gwt.HandlerFunc) gwt.HandlerFunc {
	return func(ctx *gwt.Context) error {
		r := ctx.Request()
		username, password, ok := r.BasicAuth()
		if ok && validate(username, password) {
			return handler(ctx)
		}
		return ctx.RespText(http.StatusUnauthorized, "")
	}
}

func validate(username, password string) bool {
	if username == "test" && password == "123456" {
		return true
	}
	return false
}
