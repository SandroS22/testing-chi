package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	return r
}

func StartServer(r *chi.Mux, addr string) error {
	return http.ListenAndServe(addr, r)
}

func AddRoute(r *chi.Mux, method, path string, handlerFunc http.HandlerFunc) {
	r.MethodFunc(method, path, handlerFunc)
}
