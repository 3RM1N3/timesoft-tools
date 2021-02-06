package main

import (
	"errors"
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
	usage := `欢迎使用时源科技文件夹重命名工具！

    本程序可将输入的项目目录根据 档案号+姓名 在 总表.xlsx 的第一个Sheet中匹配数据并批量将文件夹重命名为对应身份证号。

    *注意！
    1. 总表.xlsx 的第一个Sheet中应有两行标题
    2. 工号/档案号 应在第三列
    3. 姓名 应在第四列
    4. 身份证号 应在第七列
    5. *该程序会原地操作文件，故运行后将不可逆！请提前进行文件备份！*

    Usage:  将 总表.xlsx 放入程序所在文件夹，并根据提示输入项目文件夹路径即可开始。

        如：
            输入项目文件夹的绝对路径或其与程序所在文件夹的相对路径并按回车：D:\project\中国银行


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
	dict, err := getMap("总表.xlsx") // 读取总表
	if err != nil {
		log.Println(err)
		return
	}

	projectPath := ""
	fmt.Print("输入项目文件夹的绝对路径或其与程序所在文件夹的相对路径并按回车：")
	fmt.Scanln(&projectPath)
	dirArr, err := readDirs(projectPath) // 读取目录
	if err != nil {
		log.Println("读取项目文件夹失败：", err)
		return
	}
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
		err := os.Rename(a, path.Join(path.Dir(a), id))
		if err != nil {
			log.Println("重命名失败！：", err)
		}
	}
}

//  readDirs 能够获取当前目录及子目录下全部的文件夹名
func readDirs(path string) ([]string, error) {
	dirs, err := ioutil.ReadDir(path)
	if err != nil {
		return nil, err
	}
	returnedArr := []string{}
	for _, dir := range dirs {
		if dir.IsDir() {
			fullPath := path + "/" + dir.Name()
			sonDir, err := readDirs(fullPath)
			if err != nil {
				return nil, err
			}
			returnedArr = append(returnedArr, sonDir...)
			returnedArr = append(returnedArr, fullPath)
		}
	}
	return returnedArr, err
}

// getMap 能够读取总表中的工号、姓名和身份证号并制作成映射表
func getMap(fileName string) (map[string]string, error) {
	log.Printf("打开表格文件：%s\n", fileName)
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		errText := fmt.Sprintf("打开 %s 失败！", fileName)
		return nil, errors.New(errText)
	}
	cols, err := f.GetCols(f.GetSheetList()[0])
	log.Printf("读取文件中第一个Sheet：%s\n", fileName)
	if err != nil {
		errText := fmt.Sprintf("读取列失败：%s", fileName)
		return nil, errors.New(errText)
	}
	ghArr := cols[2][1:]
	nameArr := cols[3][1:]
	idArr := cols[6][1:]
	returnedMap := map[string]string{}
	for i, name := range nameArr {
		gh := ghArr[i]
		id := idArr[i]
		returnedMap[gh+name] = id
	}
	return returnedMap, nil
}
