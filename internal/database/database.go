package database

import (
	"database/sql"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "app.db")
	if err != nil {
		return nil, fmt.Errorf("problem with opening db: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("problem with ping db: %v", err)
	}

	if err := CreateTables(db); err != nil {
		return nil, fmt.Errorf("problem with creating db: %v", err)
	}

	return db, nil
}

func CreateTables(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		number INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		solved_at TIMESTAMP,
		platform_difficult INTEGER,
		my_difficult INTEGER, 
		solved_with_hint BOOLEAN,
		description TEXT,
		is_masthaved BOOLEAN,
		labels TEXT
	);`

	_, err := db.Exec(query)
	return err
}

func encodeLabels(labels []Label) string {
	if len(labels) == 0 {
		return ""
	}

	parts := make([]string, len(labels))
	for i, l := range labels {
		parts[i] = strconv.Itoa(int(l))
	}

	return strings.Join(parts, ",")
}

func decodeLabels(raw string) []Label {
	if strings.TrimSpace(raw) == "" {
		return nil
	}

	parts := strings.Split(raw, ",")
	labels := make([]Label, 0, len(parts))

	for _, p := range parts {
		v, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			continue
		}
		labels = append(labels, Label(v))
	}

	return labels
}

func AddTask(db *sql.DB, task Task) error {
	query := `INSERT INTO tasks (
		number, 
		platform_difficult, 
		my_difficult, 
		description, 
		solved_with_hint, 
		is_masthaved,
		solved_at,
		labels
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	var solvedAt interface{}
	if task.SolvedAt != nil {
		solvedAt = *task.SolvedAt
	} else {
		solvedAt = nil
	}

	_, err := db.Exec(query,
		task.Number,
		task.PlatformDifficult,
		task.MyDifficult,
		task.Description,
		task.SolvedWithHint,
		task.IsMasthaved,
		solvedAt,
		encodeLabels(task.Labels),
	)

	return err
}

func GetAllTasks(db *sql.DB) ([]Task, error) {
	query := `SELECT 
		id, 
		number, 
		created_at, 
		solved_at, 
		platform_difficult, 
		my_difficult, 
		solved_with_hint, 
		description, 
		is_masthaved,
		labels
	FROM tasks 
	ORDER BY created_at DESC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var solvedAt sql.NullTime
		var labelsRaw sql.NullString

		err := rows.Scan(
			&task.ID,
			&task.Number,
			&task.CreatedAt,
			&solvedAt,
			&task.PlatformDifficult,
			&task.MyDifficult,
			&task.SolvedWithHint,
			&task.Description,
			&task.IsMasthaved,
			&labelsRaw,
		)

		if err != nil {
			return nil, err
		}

		if solvedAt.Valid {
			task.SolvedAt = &solvedAt.Time
		}

		if labelsRaw.Valid {
			task.Labels = decodeLabels(labelsRaw.String)
		}

		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func DeleteTask(db *sql.DB, id int) error {
	query := `DELETE FROM tasks WHERE id = ?`
	result, err := db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("error retrieving task ID during delete operation: %v", err)
	}

	// проверка что строка действительно удалилась
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with id %d not found", id)
	}

	return nil
}

func FindTaskByNumber(db *sql.DB, number int) (*Task, error) {
	query := `SELECT 
		id,
		number, 
		created_at, 
		solved_at, 
		platform_difficult, 
		my_difficult, 
		solved_with_hint, 
		description, 
		is_masthaved,
		labels
	FROM tasks 
		WHERE number = ?`

	var task Task
	var solvedAt sql.NullTime
		var labelsRaw sql.NullString

	err := db.QueryRow(query, number).Scan(
		&task.ID,
		&task.Number,
		&task.CreatedAt,
		&solvedAt,
		&task.PlatformDifficult,
		&task.MyDifficult,
		&task.SolvedWithHint,
		&task.Description,
		&task.IsMasthaved,
		&labelsRaw,
	)

	if err != nil {
		return nil, err
	}

	if solvedAt.Valid {
		task.SolvedAt = &solvedAt.Time
	}

	if labelsRaw.Valid {
		task.Labels = decodeLabels(labelsRaw.String)
	}

	return &task, nil
}

func UpdateTask(db *sql.DB, task Task) error {
	query := `UPDATE tasks SET 
		platform_difficult = ?,
		my_difficult = ?,
		description = ?,
		solved_with_hint = ?,
		is_masthaved = ?,
		solved_at = ?,
		labels = ?
	WHERE id = ?`

	var solvedAt interface{}
	if task.SolvedAt != nil {
		solvedAt = *task.SolvedAt
	} else {
		solvedAt = nil
	}

	_, err := db.Exec(query,
		task.PlatformDifficult,
		task.MyDifficult,
		task.Description,
		task.SolvedWithHint,
		task.IsMasthaved,
		solvedAt,
		encodeLabels(task.Labels),
		task.ID,
	)

	return err
}

func GetRandomTasks(db *sql.DB) ([]Task, error) {
	query := `SELECT 
		id,
		number, 
		created_at, 
		solved_at, 
		platform_difficult, 
		my_difficult, 
		solved_with_hint, 
		description, 
		is_masthaved,
		labels 
	FROM tasks 
	WHERE solved_at IS NOT NULL
	  AND solved_at < datetime('now', '-14 days')
	ORDER BY solved_at`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var solvedAt sql.NullTime
		var labelsRaw sql.NullString

		err := rows.Scan(
			&task.ID,
			&task.Number,
			&task.CreatedAt,
			&solvedAt,
			&task.PlatformDifficult,
			&task.MyDifficult,
			&task.SolvedWithHint,
			&task.Description,
			&task.IsMasthaved,
			&labelsRaw,
		)

		if err != nil {
			return nil, err
		}

		if solvedAt.Valid {
			task.SolvedAt = &solvedAt.Time
		}

		if labelsRaw.Valid {
			task.Labels = decodeLabels(labelsRaw.String)
		}

		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func GetRandomTaskFromSlice(tasks []Task) (Task, error) {
	if len(tasks) == 0 {
		return Task{}, fmt.Errorf("no tasks available")
	}

	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(tasks))

	return tasks[randomIndex], nil
}
