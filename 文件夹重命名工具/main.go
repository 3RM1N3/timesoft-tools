package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func main() {
	dict := getMap("总表.xlsx")
	fmt.Println(dict)

	dirArr := readDirs(".")
	for _, a := range dirArr {
		dirName := path.Base(a)
		dirNameClean := strings.TrimSpace(dirName)
		dirNameClean = strings.ReplaceAll(dirName, " ", "")
		id, exist := dict[dirNameClean]
		if !exist {
			log.Printf("文件夹 %s 在表格中不存在！如果该文件夹非客户提交的姓名，请忽略此条消息。", a)
			continue
		}
		a = path.Clean(a)
		log.Println(a, "&&&", path.Join(path.Dir(a), id), "&&&", id)
		err := os.Rename(a, path.Join(path.Dir(a), id))
		log.Println(err)
	}

}

//  readDirs 能够获取当前目录及子目录下全部的文件夹名
func readDirs(path string) []string {
	dirs, _ := ioutil.ReadDir(path)
	returnedArr := []string{}
	for _, dir := range dirs {
		if dir.IsDir() {
			fullPath := path + "/" + dir.Name()
			returnedArr = append(returnedArr, readDirs(fullPath)...)
			returnedArr = append(returnedArr, fullPath)
		}
	}
	return returnedArr
}

// getMap 能够读取总表中的工号、姓名和身份证号并制作成映射表
func getMap(fileName string) map[string]string {
	f, _ := excelize.OpenFile(fileName)
	cols, _ := f.GetCols(f.GetSheetList()[0])
	ghArr := cols[2][1:]
	nameArr := cols[3][1:]
	idArr := cols[5][1:]
	returnedMap := map[string]string{}
	for i, name := range nameArr {
		gh := ghArr[i]
		id := idArr[i]
		returnedMap[gh+name] = id
	}
	return returnedMap
}
