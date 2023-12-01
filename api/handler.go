package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	Addr string
	Port int
	Db   *sql.DB
}

func CreateHandler(addr string, port int, db *sql.DB) *Handler {
	return &Handler{
		Addr: addr,
		Port: port,
		Db:   db,
	}
}

func (h *Handler) Start() {
	log.Printf("HTTP handler running at %s:%d\n", h.Addr, h.Port)
	r := mux.NewRouter()

	router := MakeRouter(h.Db)

	r.HandleFunc("/", router.HandleRoot).Methods("GET")
	r.HandleFunc("/health", router.HandleHealth).Methods("GET")

	r.HandleFunc("/register", router.HandleRegistration).Methods("POST")

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", h.Addr, h.Port),
		Handler: r,
	}

	srv.ListenAndServe()
}
