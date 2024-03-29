package api

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/JorgeEmanoel/money-keeper-backend/plan"
	"github.com/gorilla/handlers"
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

	planRepository := plan.MakePlanRepository(h.Db)

	usrController := MakeUserController(h.Db, baseRouter, planRepository)

	r.HandleFunc("/register", usrController.HandleRegistration).Methods(http.MethodPost)
	r.HandleFunc("/login", usrController.HandleLogin).Methods(http.MethodPost)

	authRouter := r.PathPrefix("/").Subrouter()
	authRouter.Use(AuthMiddleware)

	authRouter.HandleFunc("/me", usrController.HandleMe).Methods(http.MethodGet)

	planRouter := r.PathPrefix("/plans").Subrouter()
	planRouter.Use(AuthMiddleware)

	planController := MakePlanController(
		planRepository,
		plan.MakeTransactionRepository(h.Db),
		plan.MakeSkeletonRepository(h.Db),
		baseRouter,
	)

	planRouter.HandleFunc("", planController.HandleList).Methods(http.MethodGet)
	planRouter.HandleFunc("/summary/{period}", planController.HandleSummary).Methods(http.MethodGet)
	planRouter.HandleFunc("/init/{period}", planController.HandleInit).Methods(http.MethodGet)
	planRouter.HandleFunc("/{id}", planController.HandleDelete).Methods(http.MethodDelete)
	planRouter.HandleFunc("", planController.HandleCreate).Methods(http.MethodPost)

	skeletonRouter := r.PathPrefix("/plans/{planId}/skeletons").Subrouter()
	skeletonRouter.Use(AuthMiddleware)

	skeletonController := MakeSkeletonController(
		plan.MakeSkeletonRepository(h.Db),
		baseRouter,
	)

	skeletonRouter.HandleFunc("", skeletonController.HandleList).Methods(http.MethodGet)
	skeletonRouter.HandleFunc("/incoming", skeletonController.HandleIncomingList).Methods(http.MethodGet)
	skeletonRouter.HandleFunc("/outcoming", skeletonController.HandleOutocomingList).Methods(http.MethodGet)
	skeletonRouter.HandleFunc("/{id}", skeletonController.HandleDelete).Methods(http.MethodDelete)
	skeletonRouter.HandleFunc("", skeletonController.HandleCreate).Methods(http.MethodPost)

	transactionRouter := r.PathPrefix("/transactions").Subrouter()
	transactionRouter.Use(AuthMiddleware)

	transactionsController := MakeTransactionController(
		plan.MakeTransactionRepository(h.Db),
		baseRouter,
	)

	transactionRouter.HandleFunc("/", transactionsController.HandleList).Methods(http.MethodGet)
	transactionRouter.HandleFunc("", transactionsController.HandleCreate).Methods(http.MethodPost)
	transactionRouter.HandleFunc("/outcoming/{period}", transactionsController.HandleOutcomingList).Methods(http.MethodGet)
	transactionRouter.HandleFunc("/incoming/{period}", transactionsController.HandleIncomingList).Methods(http.MethodGet)
	transactionRouter.HandleFunc("/{id}/status/{status}", transactionsController.HandleChangeStatus).Methods(http.MethodPut)
	transactionRouter.HandleFunc("/{id}", transactionsController.HandleDelete).Methods(http.MethodDelete)

	allowedMethods := handlers.AllowedMethods([]string{"OPTIONS", "GET", "PATCH", "PUT", "POST", "DELETE"})
	allowedOrigins := handlers.AllowedOrigins([]string{"https://saturn.jorgeemanoel.com", "https://srn.jorgeemanoel.com"})
	allowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Accept", "Authorization"})

	http.ListenAndServe(
		fmt.Sprintf("%s:%d", h.Addr, h.Port),
		handlers.CORS(allowedOrigins, allowedHeaders, allowedMethods)(r),
	)
}
