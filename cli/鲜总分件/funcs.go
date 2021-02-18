package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// 读取数据源表格
func readSourcesExcel(filePath, sheetName string) (map[string][]int, error) {
	log.Println("检查表格数据：", filePath)

	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Println("表格文件打开失败：", filePath)
		return nil, err
	}
	rows, err := f.GetRows(sheetName)
	if err != nil {
		log.Printf("表格文件打开失败：%s，Sheet名称：%s\n", filePath, sheetName)
		return nil, err
	}
	sourceMap := map[string][]int{}
	haveErrors := false
	for i, row := range rows {
		if len(row) < 2 {
			log.Printf("表格数据有空值，第 1 个 Sheet，第 %d 行\n", i+1)
			haveErrors = true
			break
		}
		someThing := strings.TrimSpace(row[0])
		pageNum, err := strconv.Atoi(strings.TrimSpace(row[1]))
		if err != nil {
			log.Printf("表格数据 B 列有非整数值，第 1 个 Sheet，第 %d 行\n", i+1)
			haveErrors = true
			break
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

func do(excelPath, sheetName, folderPath, numOfBit string) {
	dataMap, err := readSourcesExcel(excelPath, sheetName) // 读取excel文件
	if err != nil {
		log.Println("表格文件打开失败：", err)
		return
	}
	fileMap, err := getFileMap(folderPath) // 读取文件夹
	if err != nil {
		log.Println("文件夹读取失败：", folderPath)
		return
	}
	fillZero := fmt.Sprintf("%%s-%%0%sd", numOfBit)
	for folderName, maxFileNum := range fileMap {
		currentDir := path.Join(folderPath, folderName)
		neededList, ok := dataMap[folderName]
		if !ok {
			continue
		}
		neededList = append(neededList, maxFileNum+1)
		for i := 0; i < len(neededList)-1; i++ {
			subFolder := fmt.Sprintf(fillZero, folderName, i+1)
			moveToDir := path.Join(currentDir, subFolder)
			os.MkdirAll(moveToDir, 0755)
			for j := neededList[i]; j < neededList[i+1]; j++ {
				fileName := fmt.Sprintf("%d.jpg", j)
				log.Printf("将 %s 从 %s 移动到 %s\n", fileName, currentDir, moveToDir)
				err := os.Rename(path.Join(currentDir, fileName), path.Join(moveToDir, fileName)) // 文件移动操作
				if err != nil {
					log.Println("移动失败：", err)
				}
			}
		}
	}
	log.Println("完成")
}

func getFileMap(folderPath string) (map[string]int, error) {
	log.Println("读取文件夹：", folderPath)
	files, err := ioutil.ReadDir(folderPath)
	if err != nil {
		log.Println("读取文件夹失败：", err)
		return nil, err
	}
	fileMap := map[string]int{}
	for _, file := range files {
		if file.IsDir() {
			fullPath := path.Join(folderPath, file.Name())
			log.Println("读取文件夹：", fullPath)
			f, err := ioutil.ReadDir(fullPath)
			if err != nil {
				log.Println("读取文件夹错误：", fullPath)
				continue
			}
			imgCounter := 0
			for _, imgFile := range f {
				if !imgFile.IsDir() && strings.HasSuffix(imgFile.Name(), ".jpg") {
					imgCounter++
				}
			}
			fileMap[file.Name()] = imgCounter
		}
	}
	return fileMap, nil
}
