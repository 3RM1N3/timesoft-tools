package main

import (
	"errors"
	"log"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// 读取数据源表格
func readSourcesExcel(filePath string) (map[string][]int, error) {
	log.Println("检查表格数据：", filePath)

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Printf("表格文件打开失败：%s，第一个 Sheet\n", filePath)
		return nil, err
	}
	sheetName := f.GetSheetList()[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Printf("表格文件打开失败：%s，第一个 Sheet\n", filePath)
		return nil, err
	}
	sourceMap := map[string][]int{}
	haveErrors := false
	for i, row := range rows {
		if len(row) < 2 {
			log.Printf("表格数据有空值，第 1 个 Sheet，第 %d 行\n", i+1)
			haveErrors = true
			continue
		}
		someThing := strings.TrimSpace(row[0])
		pageNum, err := strconv.Atoi(strings.TrimSpace(row[1]))
		if err != nil {
			log.Printf("表格数据 B 列有非整数值，第 1 个 Sheet，第 %d 行\n", i+1)
			haveErrors = true
			continue
		}
		arr, _ := sourceMap[someThing]
		arr = append(arr, pageNum)
		sourceMap[someThing] = arr
	}
	if haveErrors {
		return nil, errors.New("存在错误，无法继续")
	}
	return sourceMap, nil
}

func do(excelPath, sheetName, folderPath string) {

}
