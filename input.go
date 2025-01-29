package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func processUserInput(f *excelize.File, sheetName string, task Task) {
	fmt.Printf("Случайная задача:\nПоследняя дата решения: %s\nНомер задачи: %s\n", task.Date, task.TaskNum)

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Решили ли вы задачу? (1 - да, 0 - нет):")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "1" {
		// Обновить дату и статус задачи
		today := time.Now().Format("02-01-06")
		updateExcelCell(f, sheetName, fmt.Sprintf("A%d", task.RowNumber), today)
		updateExcelCell(f, sheetName, fmt.Sprintf("D%d", task.RowNumber), "0")
		saveExcelFile(f, "example.xlsx")

		fmt.Println("Дата в строке", task.RowNumber, "обновлена на сегодняшнюю:", today)
	} else {
		fmt.Println("Дата в строке", task.RowNumber, "осталась без изменений.")
	}
}
