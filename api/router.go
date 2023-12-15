package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type Router struct{}

func MakeRouter() *Router {
	return &Router{}
}

func (r *Router) json(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(code)
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Failed to marshal response: %v", err)
	}

	w.Write(jsonData)
}
