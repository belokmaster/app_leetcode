package main

import (
	"fmt"
	"log"
	"math/rand"
	"strconv"
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

func findTaskByNumber(f *excelize.File, sheetName string, taskNum string) (Task, error) {
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return Task{}, fmt.Errorf("ошибка при получении строк: %v", err)
	}

	for rowIndex, row := range rows {
		if len(row) == 5 {
			if row[1] == taskNum {
				return Task{
					RowNumber:   rowIndex + 1,
					Date:        row[0],
					TaskNum:     row[1],
					IsSolved:    row[2],
					Difficulty:  row[3],
					countSolved: row[4],
				}, nil
			}
		}
	}

	return Task{}, fmt.Errorf("задача с номером %s не найдена", taskNum)
}

func getNeededTasks(f *excelize.File, sheetName string, now time.Time, taskInd int) []Task {
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

		if taskInd == 0 {
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
		} else {
			if now.Sub(date).Hours() > 14*24 && row[2] == "0" {
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
	}

	return neededTasks
}

func addNewRow(f *excelize.File, sheetName string, newTask Task) {
	updateExcelCell(f, sheetName, fmt.Sprintf("A%d", newTask.RowNumber), newTask.Date)
	updateExcelCell(f, sheetName, fmt.Sprintf("B%d", newTask.RowNumber), newTask.TaskNum)
	updateExcelCell(f, sheetName, fmt.Sprintf("C%d", newTask.RowNumber), newTask.IsSolved)
	updateExcelCell(f, sheetName, fmt.Sprintf("D%d", newTask.RowNumber), newTask.Difficulty)
	updateExcelCell(f, sheetName, fmt.Sprintf("E%d", newTask.RowNumber), newTask.countSolved)
	saveExcelFile(f, "example.xlsx")
}

func pickRandomTask(tasks []Task) Task {
	rand.Seed(time.Now().UnixNano())
	randomIndex := rand.Intn(len(tasks))
	return tasks[randomIndex]
}

func changeTaskStatus(f *excelize.File, sheetName, numTask string) {
	task, err := findTaskByNumber(f, sheetName, numTask)
	if err != nil {
		log.Fatalf("Ошибка при обработке задачи. Перепроверьте ввод.")
	}

	today := time.Now().Format("02-01-06")
	updateExcelCell(f, sheetName, fmt.Sprintf("A%d", task.RowNumber), today)     // обновляем дату на сегодняшную
	updateExcelCellCountSolved(f, sheetName, fmt.Sprintf("E%d", task.RowNumber)) //+= 1 решений
	fmt.Printf("Задача %s была обновлена.\n", numTask)

	newCountSolved, err := strconv.Atoi(task.countSolved)
	if err != nil {
		fmt.Println("Ошибка при преобразовании строки в число:", err)
		return
	}

	fmt.Printf("Общее количество решений данной задачи: %d.\n", newCountSolved+1)
	saveExcelFile(f, "example.xlsx")
}

func removeTask(tasks []Task, task Task) []Task {
	for i, t := range tasks {
		if t.RowNumber == task.RowNumber {
			return append(tasks[:i], tasks[i+1:]...)
		}
	}
	return tasks
}
