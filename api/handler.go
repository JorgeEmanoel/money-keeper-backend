package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Handler struct {
	Addr string
	Port int
}

func CreateHandler(addr string, port int) *Handler {
	return &Handler{
		Addr: addr,
		Port: port,
	}
}

func (h *Handler) Start() {
	log.Printf("HTTP handler running at %s:%d\n", h.Addr, h.Port)
	r := mux.NewRouter()

	router := MakeRouter()

	r.HandleFunc("/", router.HandleRoot)
	r.HandleFunc("/health", router.HandleHealth)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", h.Addr, h.Port),
		Handler: r,
	}

	srv.ListenAndServe()
}
