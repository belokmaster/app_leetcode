package database

import (
	"database/sql"
	"fmt"
	"leetcodeapp/internal/config"
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
		is_masthaved
	) VALUES ($1, $2, $3, $4, $5, $6)`

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
