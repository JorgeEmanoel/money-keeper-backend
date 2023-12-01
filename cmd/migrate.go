package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

func Migrate(db *sql.DB) {
	defer db.Close()
	driver, err := mysql.WithInstance(db, &mysql.Config{})

	if err != nil {
		log.Fatalf("Failed to initialize MySQL instance for migration: %v", err)
	}

	currentPath, err := os.Getwd()

	if err != nil {
		log.Fatalf("Failed to get cwd: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s/migrations", currentPath),
		"migrations", driver)

	if err != nil {
		log.Fatalf("Failed to create migration instance: %v", err)
	}

	log.Println("Running migrations")
	m.Up()
	log.Println("Finished running migrations")
}
