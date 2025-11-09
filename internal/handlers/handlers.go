package handlers

import (
	"database/sql"
	"encoding/json"
	"leetcodeapp/internal/database"
	"log"
	"net/http"
	"strconv"
)

func AddTaskHandler(db *sql.DB) http.HandlerFunc {
	log.Printf("start working AddTaskHandler")
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			http.Error(w, "Can't parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Form data: %+v", r.Form)

		task := database.Task{
			Number:            parseInt(r.FormValue("number")),
			PlatformDifficult: parseInt(r.FormValue("platform_difficult")),
			MyDifficult:       parseInt(r.FormValue("my_difficult")),
			Description:       r.FormValue("description"),
			SolvedWithHint:    r.FormValue("solved_with_hint") == "on",
			IsMasthaved:       r.FormValue("is_masthaved") == "on",
		}

		log.Printf("Task to insert: %+v", task)

		err := database.AddTask(db, task)
		if err != nil {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Task inserted successfully")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func parseInt(s string) int {
	val, _ := strconv.Atoi(s)
	return val
}

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tasks, err := database.GetAllTasks(db)
		if err != nil {
			http.Error(w, "Database error with getting tasks", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}
}
