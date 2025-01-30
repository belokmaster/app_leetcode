package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	f, sheetName := openExcelFile("example.xlsx")
	defer closeExcelFile(f)

	fmt.Println("Введите значение: ")
	fmt.Println("1 - для получения случайной задачи")
	fmt.Println("2 - для получения нерешенной случайной задачи")
	fmt.Println("3 - для добавления новой задачи")
	fmt.Println("q - для выхода из программы")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" {
		fmt.Println("Выход из программы.")
		return
	}

	if input == "1" {
		neededTasks := getNeededTasks(f, sheetName, time.Now(), 1)
		randomTask := pickRandomTask(neededTasks)
		ProcessUserInput(f, sheetName, randomTask, neededTasks)
	}

	if input == "2" {
		neededTasks := getNeededTasks(f, sheetName, time.Now(), 0)

		if len(neededTasks) == 0 {
			fmt.Println("Нет ячеек, удовлетворяющих условиям. Все задачи решены самостоятельно.")
			return
		}

		randomTask := pickRandomTask(neededTasks)
		ProcessUserInput(f, sheetName, randomTask, neededTasks)
	}

	if input == "3" {
		var newTask Task
		today := time.Now().Format("02-01-06")
		newTask.Date = today

		fmt.Println("Введите номер задачи: ")
		fmt.Scan(&newTask.TaskNum)

		fmt.Println("Введите сложность задачи: ")
		fmt.Scan(&newTask.Difficulty)

		newTask.IsSolved = "0"
		newTask.countSolved = "1"

		rows, err := f.GetRows(sheetName)
		if err != nil {
			log.Fatalf("Ошибка получения строк: %v", err)
		}

		newTask.RowNumber = len(rows) + 1
		addNewRow(f, sheetName, newTask)
	}
}
