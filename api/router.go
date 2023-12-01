package api

import (
	"database/sql"
	"net/http"
)

type Router struct {
	Db *sql.DB
}

func MakeRouter(db *sql.DB) *Router {
	return &Router{
		Db: db,
	}
}

func (r *Router) HandleRoot(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
}

func (r *Router) HandleHealth(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok"))
	w.WriteHeader(200)
}
