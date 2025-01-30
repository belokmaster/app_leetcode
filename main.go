package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	// Открыть файл Excel
	f, sheetName := openExcelFile("example.xlsx")
	defer closeExcelFile(f)

	fmt.Println("Введите цифру (1 - для добавления новой задачи, 2 - для решения старых, q - выход):")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "q" || input == "1" {
		fmt.Println("Выход из программы.")
		return
	}

	if input == "2" {
		// Получаем задачи, удовлетворяющие условиям
		neededTasks := getNeededTasks(f, sheetName, time.Now())

		// Проверить, есть ли подходящие задачи
		if len(neededTasks) == 0 {
			fmt.Println("Нет ячеек, удовлетворяющих условиям")
			return
		}

		// выбираем случайную задачу
		randomTask := pickRandomTask(neededTasks)

		// обработка пользовательского ввода
		ProcessUserInput(f, sheetName, randomTask, neededTasks)
	}
}
