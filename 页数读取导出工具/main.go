package main

import (
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var fileMap = map[string]int{}

func main() {
	fmt.Println("\r\n欢迎使用时源科技文件数量统计工具！")
	usage := `

	该工具适用于统计文件夹及子文件夹内文件个数。
	
	Usage:  将该程序放入需要统计的文件夹内，双击运行即可。

	Arthur: 3RM1N3@时源科技

	E-mail: wangyu7439@hotmail.com

	`
	fmt.Println(usage)
	fmt.Println("按回车键继续...")
	fmt.Scanln()
	fmt.Println("开始读取文件...")
	t := time.Now().Unix()
	getFiles(".")
	fmt.Println("\r\n读取完毕！写入 Excel 文档...")
	f := excelize.NewFile() // 新建文件
	f.SetCellValue("Sheet1", "A1", "档案号")
	f.SetCellValue("Sheet1", "B1", "姓名")
	f.SetCellValue("Sheet1", "C1", "页数")
	i := 1
	for k, v := range fileMap {
		i++
		k = strings.TrimSpace(k)
		var id, name string
		for j, a := range k {
			i := int(a)
			if i == 32 {
				continue
			}
			if i > 127 {
				name = strings.TrimSpace(string(k[j:]))
				break
			} else {
				id += string(a)
			}
		}
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", i), id)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", i), name)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", i), v)
	}
	saveFile := "pageNum"
	for {
		err := f.SaveAs(saveFile + ".xlsx")
		if err == nil {
			break
		}
		saveFile += "(1)"
	}
	fmt.Printf("\r\n程序结束！获取 %d 个页数，用时 %v 秒，导出文件已保存为 %s 。\r\n", i, time.Now().Unix()-t, saveFile+".xlsx")
	fmt.Print("\r\n按回车键退出...")
	fmt.Scanln()
}

func getFiles(folder string) error {
	files, err := ioutil.ReadDir(folder) //specify the current dir
	if err != nil {
		return err
	}
	count := 0
	for _, file := range files {
		if file.IsDir() && !strings.Contains(file.Name(), "目录") {
			getFiles(folder + "/" + file.Name())
		} else if !file.IsDir() {
			count++
		}
	}
	dirName := path.Base(folder)
	if dirName != "." && count != 0 {
		_, exist := fileMap[dirName]
		if exist {
			dirName += "（" + folder + "）"
		}
		fileMap[dirName] = count
	}
	return nil
}
