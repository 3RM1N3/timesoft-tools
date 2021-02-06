package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

var (
	log string = ""
)

func main() {
	// 获取文件名
	args := os.Args
	if len(args) < 2 {
		usage := `
欢迎使用时源科技档案格式检测！


	Usage:  将需要处理的 .xlsx 文件拖拽到此文件上并松开即可开始，也可将文件名用作命令行参数。

	Arthur: 3RM1N3@时源科技

	E-mail: wangyu7439@hotmail.com

`
		fmt.Print(usage)
		pressExit()
	}

	fileName := args[1]
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		fmt.Println("打开文件", fileName, "失败！\r\n\r\n  ERROR: 非标准xlsx文档！")
		pressExit()
	}
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		fmt.Print("空文件！")
		pressExit()
	}
	for i, sheet := range sheetList {
		rows, err := f.GetRows(sheet)
		if err != nil {
			log += fmt.Sprintf("第%d个sheet、读取失败！\r\n", i+1)
			continue
		}
		checkAll(i+1, rows)
	}
	if log == "" {
		fmt.Print("\r\n验证通过！")
		pressExit()
	}
	log += "\r\n\r\n3RM1N3@时源科技 感谢您的使用！\r\n"
	errFile := fileName[:(len(fileName)-5)] + ".errors"
	//验证文件是否已经存在
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

// verifyId 能够验证身份证号码前17位与最后一位是否相符
func verifyID(id string) bool {
	a := [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	b := [11]string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}
	sum := 0
	for i, v := range id[:17] {
		n, err := strconv.Atoi(string(v))
		if err != nil {
			return false
		}
		sum += a[i] * n
	}
	if strings.ToUpper(string(id[17])) != b[sum%11] {
		return false
	}
	return true
}

// checkAll 能够检查一个二维数组形式的sheet数据是否符合一定标准
func checkAll(index int, rows [][]string) {
	preText := fmt.Sprintf("第%d个sheet、", index)
	if len(rows) < 3 {
		log += preText + "数据过少！请保证在有两行标题的情况下至少有一行内容\r\n"
		return
	}
	if strings.TrimSpace(rows[0][0]) != "国有企业退休人员人事档案移交案卷目录" {
		log += preText + "第1行，标题不正确！请检查标题是否为“国有企业退休人员人事档案移交案卷目录”\r\n"
		return
	}
	title := rows[1]
	titleError := ""
	if len(title) != 14 {
		log += preText + "文档列数不正确！请检查此sheet是否按顺序包含以下14列：\r\n序号、姓名、身份证号码、性别、出生日期、退休手续办理时间、生存状态、死亡日期、所属单位、户籍所在地、页数、档案局档号、保管期限、备注\r\n"
		return
	}
	switch {
	case title[1] != "姓名":
		titleError += preText + "姓名不在第二列！\r\n"
		fallthrough
	case title[2] != "身份证号码":
		titleError += preText + "身份证号码不在第三列！\r\n"
		fallthrough
	case title[3] != "性别":
		titleError += preText + "性别不在第四列！\r\n"
		fallthrough
	case title[4] != "出生日期":
		titleError += preText + "出生日期不在第五列！\r\n"
		fallthrough
	case title[5] != "退休手续办理时间":
		titleError += preText + "退休手续办理时间不在第二列！\r\n"
	}
	if titleError != "" {
		log += titleError
		return
	}
	// 定义错误身份证号码列表
	wrongIDs := [][]string{}
	// 逐行读取当前sheet
	for rowNum, row := range rows[2:] {
		preText = fmt.Sprintf("第%d个sheet、第%d行，", index, rowNum+3)
		line := [14]string{}
		// 补全整行14个值
		if len(row) > 14 {
			for i, c := range row[:14] {
				line[i] = c
			}
		} else {
			for i, c := range row {
				line[i] = c
			}
		}

		// 验证每个数据前后是否包含空格字符
		for i, l := range line {
			if l != "" && strings.TrimSpace(l) != l {
				log += preText + fmt.Sprintf("第%d列，单元格内文本前后包含空格！\r\n", i+1)
				// 比较前去掉前后空格
				line[i] = strings.TrimSpace(l)
			}
		}
		id := line[2]
		// check id
		if id == "" {
			log += preText + "证件号码为空！\r\n"
			continue
		}
		if len(id) != 18 {
			log += preText + "证件号码位数不正确！\r\n"
			continue
		}
		if !verifyID(id) {
			wrongIDs = append(wrongIDs, []string{id, fmt.Sprint(rowNum + 3)})
			continue
		}
		preText += fmt.Sprintf("证件号码：%s，", id)
		sex := line[3]
		birthday := line[4]
		retirement := line[5]
		otherInfo := map[string]string{
			"姓名":    line[1],
			"所属单位":  line[8],
			"户籍所在地": line[9],
		}
		//check sex
		if sex == "" {
			log += preText + "性别为空！\r\n"
		} else {
			trueSex := "男"
			trueSexNum, _ := strconv.Atoi(string(id[16]))
			if trueSexNum%2 == 0 {
				trueSex = "女"
			}
			if trueSex != sex {
				log += preText + "性别不正确！\r\n"
			}
		}
		// check birthday
		if birthday == "" {
			log += preText + "出生日期为空！\r\n"
		} else if len(birthday) != 8 {
			log += preText + "出生日期位数不正确！\r\n"
		} else if birthday != string(id[6:14]) {
			log += preText + "出生日期与身份证件不符！\r\n"
		}
		// check retirement
		if retirement == "" {
			log += preText + "退休手续办理时间为空！\r\n"
		} else if len(retirement) != 8 {
			log += preText + "退休手续办理时间位数不正确！\r\n"
		} else if !isDigit(retirement) {
			log += preText + "退休手续办理时间包含非数字内容！\r\n"
		}
		// check others
		for k, v := range otherInfo {
			if v == "" {
				log += preText + k + "为空！\r\n"
			}
		}
	}
	if len(wrongIDs) != 0 {
		preText = fmt.Sprintf("第%d个sheet、", index)
		for _, wID := range wrongIDs {
			log += preText + fmt.Sprintf("第%s行，", wID[1]) + wID[0] + "身份证号码格式不正确，请确认\r\n"
		}
	}
}
