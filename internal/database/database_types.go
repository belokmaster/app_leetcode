package database

import "time"

type Task struct {
	ID                int        `db:"id" json:"id"`
	Number            int        `db:"number" json:"number"`
	CreatedAt         time.Time  `db:"created_at" json:"created_at"`
	SolvedAt          *time.Time `db:"solved_at" json:"solved_at"`
	PlatformDifficult int        `db:"platform_difficult" json:"platform_difficult"`
	MyDifficult       int        `db:"my_difficult" json:"my_difficult"`
	SolvedWithHint    bool       `db:"solved_with_hint" json:"solved_with_hint"`
	Description       string     `db:"description" json:"description"`
	IsMasthaved       bool       `db:"is_masthaved" json:"is_masthaved"`
}
