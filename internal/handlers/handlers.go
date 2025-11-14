package handlers

import (
	"database/sql"
	"encoding/json"
	"leetcodeapp/internal/database"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func parseLabels(labelsStr string) []database.Label {
	if labelsStr == "" {
		return nil
	}

	labelStrs := strings.Split(labelsStr, ",")
	labels := make([]database.Label, 0, len(labelStrs))

	for _, labelStr := range labelStrs {
		if labelStr == "" {
			continue
		}

		labelVal := parseInt(strings.TrimSpace(labelStr))
		labels = append(labels, database.Label(labelVal))
	}

	return labels
}

func parseInt(s string) int {
	val, _ := strconv.Atoi(s)
	// уж надеюсь в атои не будет еррора..
	return val
}

func parseDifficulty(s string) database.Difficulty {
	val, _ := strconv.Atoi(s)
	return database.Difficulty(val)
}

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
			MyDifficult:       parseDifficulty(r.FormValue("my_difficult")),
			Description:       r.FormValue("description"),
			SolvedWithHint:    r.FormValue("solved_with_hint") == "on",
			IsMasthaved:       r.FormValue("is_masthaved") == "on",
		}

		if labelsStr := r.FormValue("labels"); labelsStr != "" {
			task.Labels = parseLabels(labelsStr)
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

func GetTasksHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			log.Printf("GetTasksHandler: problem with method")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		tasks, err := database.GetAllTasks(db)
		if err != nil {
			log.Printf("GetTasksHandler: database error: %v", err)
			http.Error(w, "Database error with getting tasks", http.StatusInternalServerError)
			return
		}

		log.Printf("GetTasksHandler: returning %d tasks", len(tasks))

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(tasks); err != nil {
			log.Printf("GetTasksHandler: error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
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
		if id <= 0 {
			log.Printf("DeleteTaskHandler: invalid ID: %d", id)
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

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
		log.Printf("GetTaskByNumberHandler: start working")
		if r.Method != "GET" {
			log.Printf("GetTaskByNumberHandler: problem with method in http")
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		numberStr := r.URL.Query().Get("number")
		log.Printf("GetTaskByNumberHandler: number from query: '%s'", numberStr)

		num, err := strconv.Atoi(numberStr)
		if err != nil || num <= 0 {
			log.Printf("GetTaskByNumberHandler: invalid number: %s", numberStr)
			http.Error(w, "Invalid task number", http.StatusBadRequest)
			return
		}

		log.Printf("GetTaskByNumberHandler: searching for task number: %d", num)

		task, err := database.FindTaskByNumber(db, num)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("GetTaskByNumberHandler: task with number %d not found", num)
				http.Error(w, "Task not found", http.StatusNotFound)
			} else {
				log.Printf("GetTaskByNumberHandler: database error: %v", err)
				http.Error(w, "Database error", http.StatusInternalServerError)
			}
			return
		}

		log.Printf("GetTaskByNumberHandler: found task ID: %d, Number: %d", task.ID, task.Number)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(task); err != nil {
			log.Printf("GetTaskByNumberHandler: error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
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
		if id <= 0 {
			log.Printf("UpdateTaskHandler: invalid ID: %d", id)
			http.Error(w, "Invalid task ID", http.StatusBadRequest)
			return
		}

		log.Printf("UpdateTaskHandler: updating task ID: %d", id)
		log.Printf("UpdateTaskHandler: form data: %+v", r.Form)

		task := database.Task{
			ID:                id,
			PlatformDifficult: parseInt(r.FormValue("platform_difficult")),
			MyDifficult:       parseDifficulty(r.FormValue("my_difficult")),
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

		if labelsStr := r.FormValue("labels"); labelsStr != "" {
			task.Labels = parseLabels(labelsStr)
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
			log.Printf("GetRandomOldTaskHandler: database error: %v", err)
			http.Error(w, "Database error with getting tasks", http.StatusInternalServerError)
			return
		}

		randomTask, err := database.GetRandomTaskFromSlice(oldTasks)
		if err != nil {
			log.Printf("GetRandomOldTaskHandler: error getting random task: %v", err)
			http.Error(w, "Database error with getting task from tasks's slices", http.StatusInternalServerError)
			return
		}

		if randomTask.ID == 0 {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"error": "Нет задач старше 2 недель"})
			return
		}

		log.Printf("GetRandomOldTaskHandler: returning task ID: %d", randomTask.ID)
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(randomTask); err != nil {
			log.Printf("GetRandomOldTaskHandler: error encoding response: %v", err)
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}
