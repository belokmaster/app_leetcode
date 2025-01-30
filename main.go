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
	fmt.Println("1 - для получения случайной задачи")
	fmt.Println("2 - для получения нерешенной случайной задачи")
	fmt.Println("q - для выхода из программы")

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" {
		fmt.Println("Выход из программы.")
		return
	}

	if input == "1" {
		// Получаем задачи, удовлетворяющие условиям
		neededTasks := getNeededTasks(f, sheetName, time.Now(), 1)

		if len(neededTasks) == 0 {
			fmt.Println("Нет ячеек, удовлетворяющих условиям")
			return
		}

		randomTask := pickRandomTask(neededTasks)
		ProcessUserInput(f, sheetName, randomTask, neededTasks)
	}

	if input == "2" {
		// Получаем задачи, удовлетворяющие условиям
		neededTasks := getNeededTasks(f, sheetName, time.Now(), 0)

		if len(neededTasks) == 0 {
			fmt.Println("Нет ячеек, удовлетворяющих условиям")
			return
		}

		randomTask := pickRandomTask(neededTasks)
		ProcessUserInput(f, sheetName, randomTask, neededTasks)
	}
}
