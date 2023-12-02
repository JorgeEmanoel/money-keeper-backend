package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/JorgeEmanoel/money-keeper-backend/sec"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		jwt := req.Header.Get("Authorization")

		if len(jwt) < 1 {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		payload, err := sec.JWTToPayload(jwt)
		if err != nil {
			log.Printf("Failed to decode token with err: %v", err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		id, _ := strconv.Atoi(fmt.Sprintf("%s", payload["id"]))

		authReq := req.Clone(context.WithValue(req.Context(), "user.id", id))

		next.ServeHTTP(w, authReq)
	})
}
