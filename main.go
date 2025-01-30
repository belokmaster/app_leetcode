package main

import (
	"fmt"
	"time"
)

func main() {
	// Открыть файл Excel
	f, sheetName := openExcelFile("example.xlsx")
	defer closeExcelFile(f)

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
