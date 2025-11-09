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
