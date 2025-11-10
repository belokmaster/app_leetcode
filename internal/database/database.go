package database

import (
	"database/sql"
	"fmt"
	"leetcodeapp/internal/config"
	"math/rand"
	"time"
)

func InitDB(path string) (*sql.DB, error) {
	config, err := config.ReadConfig(path)
	if err != nil {
		return nil, fmt.Errorf("problem with gettig config: %v", err)
	}

	connStr := config.ConnectionString()

	db, err := sql.Open("postgres", connStr)
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
		id SERIAL PRIMARY KEY,
		number INTEGER,
		created_at TIMESTAMP DEFAULT NOW(),
		solved_at TIMESTAMP,
		platform_difficult INTEGER,
		my_difficult INTEGER, 
		solved_with_hint BOOLEAN,
		description TEXT,
		is_masthaved BOOLEAN
	);`

	_, err := db.Exec(query)
	return err
}

func AddTask(db *sql.DB, task Task) error {
	query := `INSERT INTO tasks (
		number, 
		platform_difficult, 
		my_difficult, 
		description, 
		solved_with_hint, 
		is_masthaved,
		solved_at 
	) VALUES ($1, $2, $3, $4, $5, $6, NOW())`

	_, err := db.Exec(query,
		task.Number,
		task.PlatformDifficult,
		task.MyDifficult,
		task.Description,
		task.SolvedWithHint,
		task.IsMasthaved,
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
		is_masthaved 
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
		)

		if err != nil {
			return nil, err
		}

		if solvedAt.Valid {
			task.SolvedAt = &solvedAt.Time
		}

		tasks = append(tasks, task)
	}

	return tasks, rows.Err()
}

func DeleteTask(db *sql.DB, id int) error {
	query := `DELETE FROM tasks WHERE id = $1`
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
		is_masthaved 
	FROM tasks 
	WHERE number = $1`

	var task Task
	var solvedAt sql.NullTime

	err := db.QueryRow(query, number).Scan(
		&task.ID, &task.Number, &task.CreatedAt, &solvedAt,
		&task.PlatformDifficult, &task.MyDifficult, &task.SolvedWithHint,
		&task.Description, &task.IsMasthaved,
	)

	if err != nil {
		return nil, err
	}

	if solvedAt.Valid {
		task.SolvedAt = &solvedAt.Time
	}

	return &task, nil
}

func UpdateTask(db *sql.DB, task Task) error {
	query := `UPDATE tasks SET 
		platform_difficult = $1,
		my_difficult = $2,
		description = $3,
		solved_with_hint = $4,
		is_masthaved = $5,
		solved_at = $6
	WHERE id = $7`

	_, err := db.Exec(query,
		task.PlatformDifficult,
		task.MyDifficult,
		task.Description,
		task.SolvedWithHint,
		task.IsMasthaved,
		task.SolvedAt,
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
		is_masthaved 
	FROM tasks 
	WHERE solved_at < CURRENT_DATE - INTERVAL '2 weeks'
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
		)

		if err != nil {
			return nil, err
		}

		if solvedAt.Valid {
			task.SolvedAt = &solvedAt.Time
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
