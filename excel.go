package main

import (
	"log"
	"strconv"

	"github.com/xuri/excelize/v2"
)

/*
	Этот файл содержит функции для работы с Excel: открытия файла, закрытия, чтения данных и обновления ячеек.
*/

func openExcelFile(fileName string) (*excelize.File, string) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		log.Fatal(err)
	}

	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		log.Fatal("Файл Excel не содержит листов")
	}

	return f, sheetList[0]
}

func closeExcelFile(f *excelize.File) {
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func updateExcelCell(f *excelize.File, sheetName, cell, value string) {
	err := f.SetCellValue(sheetName, cell, value)
	if err != nil {
		log.Fatal(err)
	}
}

func updateExcelCellCountSolved(f *excelize.File, sheetName, cell string) {
	currentValue, err := f.GetCellValue(sheetName, cell)
	if err != nil {
		log.Fatal(err)
	}

	num, err := strconv.Atoi(currentValue)
	if err != nil {
		log.Fatal(err)
	}

	num++
	newValue := strconv.Itoa(num)

	err = f.SetCellValue(sheetName, cell, newValue)
	if err != nil {
		log.Fatal(err)
	}
}

func saveExcelFile(f *excelize.File, fileName string) {
	err := f.SaveAs(fileName)
	if err != nil {
		log.Fatal(err)
	}
}
