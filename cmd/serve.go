package cmd

import (
	"database/sql"
	"log"
	"os"
	"strconv"

	"github.com/JorgeEmanoel/money-keeper-backend/api"
)

func Serve(db *sql.DB) {
	host, ok := os.LookupEnv("HTTP_HOST")

	if !ok {
		log.Println("HTTP_HOST variable not present. Using fallback value: 0.0.0.0")
		host = "0.0.0.0"
	}

	port, ok := os.LookupEnv("HTTP_PORT")

	if !ok {
		log.Println("HTTP_PORT variable not present. Using fallback value: 8080")
		port = "8080"
	}

	portNumber, _ := strconv.Atoi(port)

	h := api.CreateHandler(host, portNumber, db)
	h.Start()
}
