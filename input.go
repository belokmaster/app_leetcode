package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func ProcessOldTaskChangeInput(f *excelize.File, sheetName string) {
	for {
		fmt.Println("Хотите ли вы изменить данные решенной задачи? (1 - да, q - выход):")
		input := ""
		fmt.Scan(&input)

		if input == "q" {
			fmt.Println("Выход из программы.")
			return
		}

		if input == "1" {
			numTask := ""
			fmt.Println("Введите номер задачи: ")
			fmt.Scan(&numTask)

			_, err := findTaskByNumber(f, sheetName, numTask)
			if err != nil {
				fmt.Println("Данная задача не существует.")
			}

			changeTaskStatus(f, sheetName, numTask)
		} else {
			fmt.Println("Некорректный ввод. Пожалуйста, введите 1 или q для выхода.")
			continue
		}
	}
}

func ProcessNewTaskInput(f *excelize.File, sheetName string) {
	for {
		fmt.Println("Хотите ли вы добавить новую задачу? (1 - да, q - выход):")
		input := ""
		fmt.Scan(&input)

		if input == "q" {
			fmt.Println("Выход из программы.")
			return
		}

		if input == "1" {
			fmt.Println("Введите номер задачи: ")
			var numTask string
			fmt.Scan(&numTask)

			_, err := findTaskByNumber(f, sheetName, numTask)
			if err == nil {
				fmt.Println("Данная задача уже существует. Обновляем информацию.")
				changeTaskStatus(f, sheetName, numTask)
				return
			}

			var newTask Task
			newTask.TaskNum = numTask
			today := time.Now().Format("02-01-06")
			newTask.Date = today

			fmt.Println("Введите сложность задачи: ")
			fmt.Scan(&newTask.Difficulty)

			fmt.Println("Решили ли вы задачу без подсказок? (0 - да, 1 - нет): ")
			fmt.Scan(&newTask.IsSolved)

			newTask.countSolved = "1"

			rows, err := f.GetRows(sheetName)
			if err != nil {
				log.Fatalf("Ошибка получения строк: %v", err)
			}

			newTask.RowNumber = len(rows) + 1
			fmt.Printf("Задача %s успешно добавлена.\n", numTask)
			addNewRow(f, sheetName, newTask)
		} else {
			fmt.Println("Некорректный ввод. Пожалуйста, введите 1 или q для выхода.")
			continue
		}
	}
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
