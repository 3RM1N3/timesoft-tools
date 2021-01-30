package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var (
	log string = ""
)

func main() {
	fmt.Println("\r\n欢迎使用时源科技 Excel 单元格文本前后空白去除工具！")
	// 获取文件名
	args := os.Args
	usage := `

	本工具可以将输入文件内所有表格的所有单元格前后的空白去除。
	
	Usage:  将需要验证的 .xlsx 文件拖拽到此文件上并松开即可开始，也可将文件名用作命令行参数。

	Arthur: 3RM1N3@时源科技

	E-mail: wangyu7439@hotmail.com

	`
	if len(args) < 2 || args[1] == "-h" {
		fmt.Print(usage)
		pressExit()
	}

	fileName := args[1]
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println("打开文件", fileName, "失败！\r\n\r\n  ERROR: 非标准xlsx文档！")
		pressExit()
	}

	// 获取文件中的sheet列表
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		fmt.Print("空文件！")
		pressExit()
	}

	// 显示警告信息
	fmt.Println("\r\n*注意！该程序会操作原文件，故运行后将不可逆！请提前进行文件备份！")
	fmt.Print("\r\n\r\n确定要继续吗？(输入y后按回车键继续，否则退出程序):")
	goon := ""
	fmt.Scanln(&goon)
	goon = strings.TrimSpace(goon)
	if goon != "y" && goon != "Y" {
		pressExit()
	}

	// 读取文件中的每个sheet
	for i, sheet := range sheetList {
		rows, err := f.GetRows(sheet)
		if err != nil {
			log += fmt.Sprintf("第%d个sheet、读取失败！\r\n", i+1)
			continue
		}

		// 逐行检查有无空白
		if len(rows) == 0 {
			continue
		}
		for y, row := range rows {
			if len(row) == 0 {
				continue
			}
			for x, col := range row {
				trimed := strings.TrimSpace(col)
				if trimed == col {
					continue
				}
				axis := getAxis(x, y)
				f.SetCellValue(sheet, axis, trimed)
				log += fmt.Sprintf("第%d个sheet、单元格%s：“%s”-已修复\r\n", i+1, axis, col)
			}
		}
	}

	if log == "" {
		fmt.Print("\r\n验证通过！")
		pressExit()
	}
	// 保存文件
	if err := f.Save(); err != nil {
		fmt.Println("\r\n文件保存失败！请检查文件是否被占用！")
		pressExit()
	}

	log += "\r\n\r\n3RM1N3@时源科技 感谢您的使用！\r\n"
	errFile := fileName[:(len(fileName)-5)] + ".log"

	// 验证log文件是否已经存在
	for {
		if _, err := os.Stat(errFile + ".txt"); os.IsNotExist(err) {
			break
		}
		errFile += "(1)"
	}
	errFile += ".txt"
	if err := writeLog(errFile, log); err != nil {
		fmt.Print("日志文件写入失败！\r\n\r\n")
		fmt.Println(log)
		pressExit()
	}
	errFile = "\"" + errFile + "\""
	cmd := exec.Command("powershell", "/c", "notepad", errFile)
	err = cmd.Run()
	if err != nil {
		fmt.Println("无法写入错误文件！请关闭程序后手动打开", errFile, "文件！")
		pressExit()
	}
}

// 通过x和y生成excel坐标文本
func getAxis(x, y int) string {
	return fmt.Sprintf("%s%d", getX(x), y+1)
}

// 生成横坐标字母
func getX(x int) string {
	abc := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	if x < 26 {
		return string(abc[x])
	}
	mod := (x + 1) % 26
	return getX(int((x+1)/26)-1) + string(abc[mod-1])
}

// isDigit 用于验证字符串是否为纯数字
func isDigit(s string) bool {
	for _, a := range s {
		i := int(a)
		if i < 48 || i > 57 {
			return false
		}
	}
	return true
}

// pressExit 实现了按回车键退出的功能
func pressExit() {
	fmt.Println()
	fmt.Println("\r\n按回车键退出...")
	fmt.Scanln()
	os.Exit(0)
}

// writeLog 用于写入文本文件
func writeLog(filePath, s string) error {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer file.Close()
	file.WriteString(s)
	return nil
}
