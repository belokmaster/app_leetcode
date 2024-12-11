package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
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

	// Создать срез для хранения адресов ячеек, которые удовлетворяют условиям
	var cellAddresses []string

	// Пройтись по всем строкам и проверить значения в первом и четвертом столбцах
	for i, row := range rows {
		if len(row) >= 4 {
			// Парсить дату из первого столбца
			date, err := time.Parse("02-01-06", row[0]) // Измените формат даты, если он отличается
			if err != nil {
				continue
			}

			// Проверить, прошло ли больше двух недель с даты
			if now.Sub(date).Hours() > 14*24 {
				// Проверить значение в четвертом столбце
				if row[3] != "0" {
					cell := fmt.Sprintf("A%d", i+1)
					cellAddresses = append(cellAddresses, cell) // сохранить адрес ячейки
				}
			}
		}
	}

	// Проверить, есть ли значения в срезе
	if len(cellAddresses) == 0 {
		fmt.Println("Нет ячеек, удовлетворяющих условиям")
		return
	}

	// Инициализировать источник случайных чисел
	rand.Seed(time.Now().UnixNano())

	// Получить случайный индекс и значение из cellAddresses
	randomIndex := rand.Intn(len(cellAddresses))
	randomCell := cellAddresses[randomIndex]

	fmt.Println("Случайная ячейка, удовлетворяющая условиям:", randomCell)

	// Спросить у пользователя, решена ли задача
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Решили ли вы задачу? (1 - да, 0 - нет):")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	// Проверить ввод пользователя
	if input == "1" {
		// Обновить дату в соответствующей ячейке на текущую дату
		today := now.Format("02-01-06")
		rowNumber, err := strconv.Atoi(randomCell[1:])
		if err != nil {
			log.Fatal(err)
		}
		cell := fmt.Sprintf("A%d", rowNumber)
		err = f.SetCellValue(sheetName, cell, today)
		if err != nil {
			log.Fatal(err)
		}

		// Сохранить изменения в файл
		if err := f.SaveAs("example.xlsx"); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Дата в ячейке", randomCell, "обновлена на сегодняшнюю:", today)
	} else {
		fmt.Println("Дата в ячейке", randomCell, "осталась без изменений.")
	}
}
