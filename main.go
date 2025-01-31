package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	f, sheetName := openExcelFile("example.xlsx")
	defer closeExcelFile(f)

	fmt.Println("Введите значение: ")
	fmt.Println("1 - для получения случайной задачи;")
	fmt.Println("2 - для получения нерешенной случайной задачи;")
	fmt.Println("3 - для добавления новой задачи;")
	fmt.Println("4 - для самостоятельного изменения статуса задачи;")
	fmt.Println("q - для выхода из программы.")

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
		ProcessNewTaskInput(f, sheetName)
	}

	if input == "4" {
		ProcessOldTaskChangeInput(f, sheetName)
	}
}
