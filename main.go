package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

func main() {
	// Открыть файл Excel
	f, err := excelize.OpenFile("example.xlsx")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		// Закрыть файл Excel
		if err := f.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	// Получить все листы в файле
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		log.Fatal("Файл Excel не содержит листов")
	}

	// Использовать первый лист в файле
	sheetName := sheetList[0]

	// Получить все строки в первом столбце
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Fatal(err)
	}

	// Текущая дата
	now := time.Now()

	// Список задач, удовлетворяющих условиям
	var neededTasks []struct {
		RowNumber int
		Date      string
		TaskNum   string
		IsSolved  string
	}

	// Парсинг строк из Excel
	for i, row := range rows {
		if i == 0 {
			continue
		}

		if len(row) >= 4 {
			// Парсить дату из первого столбца
			date, err := time.Parse("02-01-06", row[0])
			if err != nil {
				continue
			}

			// Проверить, прошло ли больше двух недель с даты
			if now.Sub(date).Hours() > 14*24 {
				// Проверить значение в четвертом столбце
				if row[3] != "0" {
					rowNumber := i + 1 // Номер строки в Excel
					neededTasks = append(neededTasks, struct {
						RowNumber int
						Date      string
						TaskNum   string
						IsSolved  string
					}{
						RowNumber: rowNumber,
						Date:      row[0],
						TaskNum:   row[1],
						IsSolved:  row[3],
					})
				}
			}
		}
	}

	// Проверить, есть ли значения в мапе
	if len(neededTasks) == 0 {
		fmt.Println("Нет ячеек, удовлетворяющих условиям")
		return
	}

	// Инициализировать источник случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Получить случайный индекс и значение из cellAddresses
	randomIndex := rand.Intn(len(neededTasks))
	randomTask := neededTasks[randomIndex]

	fmt.Printf("Случайная задача:\nПоследняя дата решения: %s\nНомер задачи: %s\n", randomTask.Date, randomTask.TaskNum)

	// Спросить у пользователя, решена ли задача
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Решили ли вы задачу? (1 - да, 0 - нет):")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Проверить ввод пользователя
	if input == "1" {
		// Обновить дату в соответствующей ячейке на текущую дату
		today := now.Format("02-01-06")
		cell := fmt.Sprintf("A%d", randomTask.RowNumber) // Ячейка с датой
		err = f.SetCellValue(sheetName, cell, today)
		if err != nil {
			log.Fatal(err)
		}

		// Обновить значение в столбце "solved task" на 0
		solvedCell := fmt.Sprintf("D%d", randomTask.RowNumber) // Ячейка с решением задачи
		err = f.SetCellValue(sheetName, solvedCell, "0")
		if err != nil {
			log.Fatal(err)
		}

		// Сохранить изменения в файл
		if err := f.SaveAs("example.xlsx"); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Дата в строке", randomTask.RowNumber, "обновлена на сегодняшнюю:", today)
	} else {
		fmt.Println("Дата в строке", randomTask.RowNumber, "осталась без изменений.")
	}
}
