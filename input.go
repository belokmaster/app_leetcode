package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// ProcessUserInput обрабатывает пользовательский ввод и обновляет Excel.
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
			// Обновить дату и статус задачи
			today := time.Now().Format("02-01-06")
			updateExcelCell(f, sheetName, fmt.Sprintf("A%d", task.RowNumber), today)
			updateExcelCell(f, sheetName, fmt.Sprintf("C%d", task.RowNumber), "0")
			updateExcelCellCountSolved(f, sheetName, fmt.Sprintf("E%d", task.RowNumber))
			saveExcelFile(f, "example.xlsx")

			fmt.Println("Дата в строке", task.RowNumber, "обновлена на сегодняшнюю:", today)
		} else if input == "0" {
			fmt.Println("Дата в строке", task.RowNumber, "осталась без изменений.")
		} else {
			fmt.Println("Некорректный ввод. Пожалуйста, введите 1, 0 или q для выхода.")
			continue
		}

		// Удаляем обработанную задачу из списка
		neededTasks = removeTask(neededTasks, task)

		// Проверяем, остались ли задачи
		if len(neededTasks) == 0 {
			fmt.Println("Нет задач, удовлетворяющих условиям.")
			return
		}

		// Выбираем следующую случайную задачу
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
