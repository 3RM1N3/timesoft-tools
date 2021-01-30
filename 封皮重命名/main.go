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
	imagesPath := ""
	fmt.Print("输入装有imageXXXX.jpg文件的目录：")
	fmt.Scanln(&imagesPath)
	fileNum, _ := getFileNum(imagesPath)
	idList, _ := getIDList("案卷目录.xlsx")
	fmt.Println(fileNum, idList)

	if len(idList) != fileNum {
		log.Fatal("文件数目与身份证数量不等！")
	}

	for i := 1; i < fileNum+1; i++ {
		file := fmt.Sprintf("image%04d.jpg", i)
		os.Rename(path.Join(imagesPath, file), path.Join(imagesPath, idList[i-1]+".jpg"))
	}
	fmt.Println("重命名成功！")
}

func getFileNum(dirPath string) (int, error) {
	files, err := ioutil.ReadDir(dirPath)
	fileNum := 0
	if err != nil {
		log.Fatal("读取目录失败！")
		return fileNum, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.Contains(file.Name(), "image") && strings.HasSuffix(file.Name(), ".jpg") {
			fileNum++
		}
	}
	return fileNum, nil
}

func getIDList(fileName string) ([]string, error) {
	idList := []string{}
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		log.Fatal("打开表格失败！")
		return idList, err
	}
	sheet := f.GetSheetList()[0]
	cols, _ := f.GetCols(sheet)
	idList = cols[2][2:]
	return idList, nil
}
