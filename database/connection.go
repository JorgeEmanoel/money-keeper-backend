package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	dbHost, ok := os.LookupEnv("DB_HOST")

	if !ok {
		log.Fatal("Missing DB_HOST env")
	}

	dbPort, ok := os.LookupEnv("DB_PORT")

	if !ok {
		log.Fatal("Missing DB_PORT env")
	}

	dbUser, ok := os.LookupEnv("DB_USER")

	if !ok {
		log.Fatal("Missing DB_USER env")
	}

	dbPassword, ok := os.LookupEnv("DB_PASSWORD")

	if !ok {
		log.Fatal("Missing DB_PASSWORD env")
	}

	dbName, ok := os.LookupEnv("DB_NAME")

	if !ok {
		log.Fatal("Missing DB_NAME env")
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@(%s:%s)/%s?multiStatements=true", dbUser, dbPassword, dbHost, dbPort, dbName))

	if err != nil {
		log.Fatalf("Failed to connect to database host: %s. Error: %v", dbHost, err)
	}

	return db
}
