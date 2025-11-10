package main

import (
	"database/sql"
	"leetcodeapp/internal/handlers"
	"net/http"
)

func setupRoutes(db *sql.DB) {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/templates/index.html")
	})

	http.HandleFunc("/tasks", handlers.AddTaskHandler(db))
	http.HandleFunc("/tasks/delete", handlers.DeleteTaskHandler(db))

	http.HandleFunc("/api/tasks", handlers.GetTasksHandler(db))
	http.HandleFunc("/api/task", handlers.GetTaskByNumberHandler(db))
	http.HandleFunc("/api/tasks/update", handlers.UpdateTaskHandler(db))
	http.HandleFunc("/api/tasks/random-old", handlers.GetRandomOldTaskHandler(db))
}
