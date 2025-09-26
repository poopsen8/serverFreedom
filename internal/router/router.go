package router

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

type Router struct {
	muxRouter *mux.Router
}

func NewRouter() *Router {
	return &Router{
		muxRouter: mux.NewRouter(),
	}
}

func (r *Router) Add(path string, handler http.HandlerFunc) {
	// Проверяем, есть ли в пути параметры {id}
	if strings.Contains(path, "{") {
		r.muxRouter.HandleFunc(path, handler)
	} else {
		r.muxRouter.HandleFunc(path, handler).Methods("POST", "GET", "PUT", "DELETE")
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.muxRouter.ServeHTTP(w, req)
}
