package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

func main() {
	usage := `欢迎使用时源科技封皮重命名工具！

    本程序可将 输入的装有 imageXXXX.jpg文件 的目录 下全部 imageXXXX.jpg 文件按顺序依次重命名为 案卷目录.xlsx 文档第一个Sheet中的第三列（身份证号）。

    *注意！
    1. 案卷目录.xlsx 的第一个Sheet中应有两行标题
    2. 程序仅查找输入的文件夹，不会递归查找子文件夹
	3. 身份证号 应在第三列
	4. 应保证身份证号列的身份证号数量与读取到的jpg文件数量完全一致
    5. *该程序会原地操作文件，故运行后将不可逆！请提前进行文件备份！*

    Usage:  将 案卷目录.xlsx 与程序放入同一目录下，并根据提示输入装有 imageXXXX.jpg文件 的文件夹路径即可开始。

        如：
            输入装有 imageXXXX.jpg文件 的目录 的绝对路径或其与程序所在文件夹的相对路径并按回车：D:\project\中国银行


    Author: 3RM1N3@时源科技

    E-mail: wangyu7439@hotmail.com


    按回车键继续...`
	fmt.Print(usage)
	fmt.Scanln()
	fmt.Println("*请再次确认！此操作运行后将不可逆！请提前进行文件备份！确定仍要继续吗？")
	fmt.Println()
	for i := 5; i > 0; i-- {
		fmt.Printf("\033[30D%d 秒后按回车键继续...", i)
		time.Sleep(1 * time.Second)
	}
	fmt.Println("\033[30D\033[K按回车键继续...")
	fmt.Scanln()

	defer func() {
		fmt.Printf("\n操作完成，按回车键退出...")
		fmt.Scanln()
	}()

	// 正文开始
	imagesPath := ""
	fmt.Print("输入装有imageXXXX.jpg文件的目录：")
	fmt.Scanln(&imagesPath)
	fileNum, err := getFileNum(imagesPath)
	if err != nil {
		log.Println("打开文件夹失败：", err)
	}
	idList, err := getIDList("案卷目录.xlsx")
	if err != nil {
		log.Println("打开“案卷目录.xlsx”失败！", err)
		return
	}
	fmt.Println(fileNum, idList)

	if len(idList) != fileNum {
		log.Println("文件数目与身份证数量不等！")
		return
	}

	for i := 1; i < fileNum+1; i++ {
		file := fmt.Sprintf("image%04d.jpg", i)
		fullName := path.Join(imagesPath, file)
		err = os.Rename(fullName, path.Join(imagesPath, idList[i-1]+".jpg"))
		if err != nil {
			log.Println("重命名失败：", fullName)
		}
	}
	fmt.Println("重命名成功！")
}

func getFileNum(dirPath string) (int, error) {
	files, err := ioutil.ReadDir(dirPath)
	fileNum := 0
	if err != nil {
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

	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	sheet := f.GetSheetList()[0]
	cols, err := f.GetCols(sheet)
	if err != nil {
		return nil, err
	}
	idList := cols[2][2:]
	return idList, nil
}
