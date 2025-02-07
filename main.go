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

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println()
		outputStartMessage()

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		switch input {
		case "q":
			fmt.Println("Выход из программы.")
			return
		case "1":
			neededTasks := getNeededTasks(f, sheetName, time.Now(), 1)
			randomTask := pickRandomTask(neededTasks)
			ProcessUserInput(f, sheetName, randomTask, neededTasks)
		case "2":
			neededTasks := getNeededTasks(f, sheetName, time.Now(), 0)
			if len(neededTasks) == 0 {
				fmt.Println("Нет ячеек, удовлетворяющих условиям. Все задачи решены самостоятельно.")
				continue
			}
			randomTask := pickRandomTask(neededTasks)
			ProcessUserInput(f, sheetName, randomTask, neededTasks)
		case "3":
			ProcessNewTaskInput(f, sheetName)
		case "4":
			ProcessOldTaskChangeInput(f, sheetName)
		default:
			fmt.Println("Некорректный ввод, попробуйте снова.")
		}
	}
}
