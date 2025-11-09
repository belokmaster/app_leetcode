package database

import "time"

type Task struct {
	ID                int       `db:"id"`
	Number            int       `db:"number"`
	CreatedAt         time.Time `db:"created_at"`
	SolvedAt          time.Time `db:"solved_at"`
	PlatformDifficult int       `db:"platform_difficult"`
	MyDifficult       int       `db:"my_difficult"`
	SolvedWithHint    bool      `db:"solved_with_hint"`
	Description       string    `db:"description"`
	IsMasthaved       bool      `db:"is_masthaved"`
}
