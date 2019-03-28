package gwt

import "net/http"

type (
	Router struct {
		maps map[string]http.HandlerFunc
	}
)

func NewRouter() *Router {
	r := &Router{
		maps: make(map[string]http.HandlerFunc),
	}
	return r
}

// register path&method pair to router
func (r *Router) Register(path string, h http.HandlerFunc) {
	r.maps[path] = h
}

func (r *Router) Detect(path string) http.HandlerFunc {
	h, _ := r.maps[path]
	return h
}
