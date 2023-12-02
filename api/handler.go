package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/JorgeEmanoel/money-keeper-backend/plan"
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
	baseRouter := MakeRouter()

	usrController := MakeUserController(h.Db, baseRouter)

	r.HandleFunc("/register", usrController.HandleRegistration).Methods(http.MethodPost)
	r.HandleFunc("/login", usrController.HandleLogin).Methods(http.MethodPost)

	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.Use(AuthMiddleware)

	authRouter.HandleFunc("/me", usrController.HandleMe).Methods(http.MethodGet)

	planRouter := r.PathPrefix("/plans").Subrouter()
	planRouter.Use(AuthMiddleware)

	planController := MakePlanController(
		plan.MakePlanRepository(h.Db),
		baseRouter,
	)

	planRouter.HandleFunc("", planController.HandleList).Methods(http.MethodGet)
	planRouter.HandleFunc("/{id}", planController.HandleDelete).Methods(http.MethodDelete)
	planRouter.HandleFunc("", planController.HandleCreate).Methods(http.MethodPost)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", h.Addr, h.Port),
		Handler: r,
	}

	srv.ListenAndServe()
}
