package api

import "net/http"

type Router struct{}

func MakeRouter() *Router {
	return &Router{}
}

func (r *Router) HandleRoot(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
}

func (r *Router) HandleHealth(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok"))
	w.WriteHeader(200)
}
