package main

import (
	"leetcodeapp/internal/database"
	"leetcodeapp/internal/handlers"
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

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/index.html")
	})

	http.HandleFunc("/tasks", handlers.AddTaskHandler(db))
	http.HandleFunc("/api/tasks", handlers.GetTasksHandler(db))
	http.HandleFunc("/tasks/delete", handlers.DeleteTaskHandler(db))
	http.HandleFunc("/api/task", handlers.GetTaskByNumberHandler(db))
	http.HandleFunc("/api/tasks/update", handlers.UpdateTaskHandler(db))

	log.Println("Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
