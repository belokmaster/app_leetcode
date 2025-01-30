package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/xuri/excelize/v2"
)

type Task struct {
	RowNumber   int
	Date        string
	TaskNum     string
	IsSolved    string
	Difficulty  string
	countSolved string
}

func getNeededTasks(f *excelize.File, sheetName string, now time.Time) []Task {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatal(err)
	}

	var neededTasks []Task
	for i, row := range rows {
		if i == 0 || len(row) < 4 {
			continue
		}

		date, err := time.Parse("02-01-06", row[0])
		if err != nil {
			continue
		}

		if now.Sub(date).Hours() > 14*24 && row[2] != "0" {
			neededTasks = append(neededTasks, Task{
				RowNumber:   i + 1,
				Date:        row[0],
				TaskNum:     row[1],
				IsSolved:    row[2],
				Difficulty:  row[3],
				countSolved: row[4],
			})
		}
	}

	return neededTasks
}

func pickRandomTask(tasks []Task) Task {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(tasks))
	return tasks[randomIndex]
}
