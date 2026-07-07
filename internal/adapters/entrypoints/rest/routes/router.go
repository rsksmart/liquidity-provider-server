package routes

import "net/http"

// Router registers HTTP routes and builds the final handler with middleware.
type Router interface {
	Handle(pattern string, handler http.Handler)
	Use(middleware ...func(http.Handler) http.Handler)
	Handler(req *http.Request) (http.Handler, string)
	BuildHandler() http.Handler
}

type routerImpl struct {
	mux         *http.ServeMux
	middlewares []func(http.Handler) http.Handler
}

func NewRouter() Router {
	return &routerImpl{mux: http.NewServeMux()}
}

func (r *routerImpl) Use(middleware ...func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, middleware...)
}

func (r *routerImpl) Handle(pattern string, handler http.Handler) {
	r.mux.Handle(pattern, handler)
}

func (r *routerImpl) Handler(req *http.Request) (http.Handler, string) {
	return r.mux.Handler(req)
}

func (r *routerImpl) BuildHandler() http.Handler {
	handler := http.Handler(r.mux)
	for i := len(r.middlewares) - 1; i >= 0; i-- {
		handler = r.middlewares[i](handler)
	}
	return handler
}
