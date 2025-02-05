package main

import (
	"backend/api"
	"backend/pkg/db/sqlite"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	mux := mux.NewRouter()
	dbPath := "./pkg/db/database.db"
	migrationsPath := "pkg/db/migrations/sqlite"

	// Connect to the database and apply migrations
	db, err := sqlite.ConnectAndMigrate(dbPath, migrationsPath)
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
	}
	defer db.Close()

	api.Router(mux, db)

	fmt.Println("Server is running on http://localhost:8080")
	http.ListenAndServe(":8080", nil)
}
