package handlers

import (
	"database/sql"
	"encoding/json"
	"leetcodeapp/internal/database"
	"log"
	"net/http"
	"strconv"
	"time"
)

func AddTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("AddTaskHandler: start working")

		if r.Method != "POST" {
			log.Printf("AddTaskHandler: problem with method in http")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("AddTaskHandler: problem with parsing the form")
			http.Error(w, "Can't parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("AddTaskHandler: Form data: %+v", r.Form)

		task := database.Task{
			Number:            parseInt(r.FormValue("number")),
			PlatformDifficult: parseInt(r.FormValue("platform_difficult")),
			MyDifficult:       parseInt(r.FormValue("my_difficult")),
			Description:       r.FormValue("description"),
			SolvedWithHint:    r.FormValue("solved_with_hint") == "on",
			IsMasthaved:       r.FormValue("is_masthaved") == "on",
		}

		log.Printf("AddTaskHandler: Task to insert: %+v", task)

		err := database.AddTask(db, task)
		if err != nil {
			log.Printf("AddTaskHandler: problem with add task to db: %v", err)
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("AddTaskHandler: task inserted successfully")
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func parseInt(s string) int {
	val, _ := strconv.Atoi(s)
	// уж надеюсь в атои не будет еррора..
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

		log.Printf("GetTasksHandler: returning %d tasks", len(tasks))

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(tasks)
	}
}

func DeleteTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("DeleteTaskHandler: start working")

		if r.Method != "POST" {
			log.Printf("DeleteTaskHandler: problem with method in http")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("DeleteTaskHandler: problem with parsing the form")
			http.Error(w, "Can't parse form: "+err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("DeleteTaskHandler: form data: %+v", r.Form)

		id := parseInt(r.FormValue("id"))
		log.Printf("DeleteTaskHandler: deleting task ID: %d", id)

		err := database.DeleteTask(db, id)
		if err != nil {
			log.Printf("DeleteTaskHandler: delete failed: %v", err)
			http.Error(w, "Delete error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("DeleteTaskHandler: task %d deleted successfully", id)
		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}

func GetTaskByNumberHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("GetTaskByIDHandler: start working")
		if r.Method != "GET" {
			log.Printf("GetTaskByIDHandler: problem with method in http")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		idStr := r.URL.Query().Get("id")
		log.Printf("GetTaskByIDHandler: ID from query: '%s'", idStr)

		id, err := strconv.Atoi(idStr)
		if err != nil || id <= 0 {
			log.Printf("GetTaskByIDHandler: invalid ID task")
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		log.Printf("GetTaskByIDHandler: task ID: %d", id)

		task, err := database.FindTaskByNumber(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Task not found", http.StatusNotFound)
			} else {
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(task)
	}
}

func UpdateTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("UpdateTaskHandler: start working")

		if r.Method != "POST" {
			log.Printf("UpdateTaskHandler: problem with method in http")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := r.ParseForm(); err != nil {
			log.Printf("UpdateTaskHandler: parse form error: %v", err)
			http.Error(w, "Can't parse form", http.StatusBadRequest)
			return
		}

		id := parseInt(r.FormValue("id"))
		log.Printf("UpdateTaskHandler: updating task ID: %d", id)
		log.Printf("UpdateTaskHandler: form data: %+v", r.Form)

		task := database.Task{
			ID:                id,
			PlatformDifficult: parseInt(r.FormValue("platform_difficult")),
			MyDifficult:       parseInt(r.FormValue("my_difficult")),
			Description:       r.FormValue("description"),
			SolvedWithHint:    r.FormValue("solved_with_hint") == "on",
			IsMasthaved:       r.FormValue("is_masthaved") == "on",
		}

		// Обработка даты решения
		if solvedAtStr := r.FormValue("solved_at"); solvedAtStr != "" {
			log.Printf("UpdateTaskHandler: solved_at from form: %s", solvedAtStr)
			if solvedAt, err := time.Parse("2006-01-02", solvedAtStr); err == nil {
				task.SolvedAt = &solvedAt
				log.Printf("UpdateTaskHandler: parsed solved_at: %v", solvedAt)
			} else {
				log.Printf("UpdateTaskHandler: solved_at parse error: %v", err)
			}
		}

		err := database.UpdateTask(db, task)
		if err != nil {
			log.Printf("UpdateTaskHandler: database error: %v", err)
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("UpdateTaskHandler: task %d updated successfully", id)
		w.WriteHeader(http.StatusOK)
	}
}

func GetRandomOldTaskHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("GetRandomOldTaskHandler: start working")
		if r.Method != "GET" {
			log.Printf("GetRandomOldTaskHandler: problem with method in http")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		oldTasks, err := database.GetRandomTasks(db)
		if err != nil {
			http.Error(w, "Database error with getting tasks", http.StatusInternalServerError)
			return
		}

		randomTask, err := database.GetRandomTaskFromSlice(oldTasks)
		if err != nil {
			http.Error(w, "Database error with getting task from tasks's slices", http.StatusInternalServerError)
			return
		}

		if randomTask.ID == 0 {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "Нет задач старше 2 недель"})
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(randomTask)
	}
}
