package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func addNewRow(f *excelize.File, sheetName string, newTask Task) {
	updateExcelCell(f, sheetName, fmt.Sprintf("A%d", newTask.RowNumber), newTask.Date)
	updateExcelCell(f, sheetName, fmt.Sprintf("B%d", newTask.RowNumber), newTask.TaskNum)
	updateExcelCell(f, sheetName, fmt.Sprintf("C%d", newTask.RowNumber), newTask.IsSolved)
	updateExcelCell(f, sheetName, fmt.Sprintf("D%d", newTask.RowNumber), newTask.Difficulty)
	updateExcelCell(f, sheetName, fmt.Sprintf("E%d", newTask.RowNumber), newTask.countSolved)
	saveExcelFile(f, "example.xlsx")
}

func ProcessUserInput(f *excelize.File, sheetName string, task Task, neededTasks []Task) {
	for {
		fmt.Printf("Случайная задача:\n[ ] Последняя дата решения: %s\n[ ] Номер задачи: %s\n", task.Date, task.TaskNum)
		reader := bufio.NewReader(os.Stdin)

		fmt.Println("Решили ли вы задачу? (1 - да, 0 - нет, q - выход):")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "q" {
			fmt.Println("Выход из программы.")
			return
		}

		if input == "1" {
			today := time.Now().Format("02-01-06")
			updateExcelCell(f, sheetName, fmt.Sprintf("A%d", task.RowNumber), today)     // обновляем дату на сегодняшную
			updateExcelCell(f, sheetName, fmt.Sprintf("C%d", task.RowNumber), "0")       // обнуляем счетчик решения с подсказкой
			updateExcelCellCountSolved(f, sheetName, fmt.Sprintf("E%d", task.RowNumber)) //+= 1 решений
			saveExcelFile(f, "example.xlsx")

			fmt.Println("Дата в строке", task.RowNumber, "обновлена на сегодняшнюю:", today)
		} else if input == "0" {
			fmt.Println("Дата в строке", task.RowNumber, "осталась без изменений.")
		} else {
			fmt.Println("Некорректный ввод. Пожалуйста, введите 1, 0 или q для выхода.")
			continue
		}

		neededTasks = removeTask(neededTasks, task)

		if len(neededTasks) == 0 {
			fmt.Println("Нет задач, удовлетворяющих условиям.")
			return
		}

		task = pickRandomTask(neededTasks)
	}
}

// removeTask удаляет задачу из списка задач.
func removeTask(tasks []Task, task Task) []Task {
	for i, t := range tasks {
		if t.RowNumber == task.RowNumber {
			return append(tasks[:i], tasks[i+1:]...)
		}
	}
	return tasks
}
