package main

import (
	"leetcodeapp/internal/database"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	path := "config.json"
	db, err := database.InitDB(path)
	if err != nil {
		log.Fatalf("problem with init db: %v", err)
	}
	defer db.Close()

	setupRoutes(db)

	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
