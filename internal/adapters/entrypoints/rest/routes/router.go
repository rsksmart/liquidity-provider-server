package routes

import "net/http"

type Router struct {
	mux         *http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

func NewRouter() *Router {
	return &Router{mux: http.NewServeMux()}
}

func (r *Router) Use(middleware ...func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, middleware...)
}

func (r *Router) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, handler)
}

func (r *Router) Handler(req *http.Request) (http.Handler, string) {
	return r.mux.Handler(req)
}

func (r *Router) BuildHandler() http.Handler {
	handler := http.Handler(r.mux)
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}
	return handler
}
