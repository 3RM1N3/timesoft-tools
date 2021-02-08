package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
)

// moveFilesInCatalogs 能够将给出文件夹下名称为“目录”的文件夹下的文件转移至上层文件夹
func moveFilesInCatalogs(dirPath string) {
	logChan <- "开始运行..."
	logChan <- "finding " + dirPath
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		logChan <- fmt.Sprintf("读取目录错误：%v", err)
		return
	}
	for _, file := range files {
		fullPath := path.Join(dirPath, file.Name())
		fullPath = path.Clean(fullPath)
		if file.IsDir() && file.Name() == "目录" {
			filesInMulu, err := ioutil.ReadDir(fullPath)
			if err != nil {
				logChan <- fmt.Sprintf("读取目录错误：%v", err)
			}
			for _, file := range filesInMulu {
				if !file.IsDir() {
					aimPath := path.Dir(fullPath)
					fileName := file.Name()
					aimPath = path.Join(aimPath, fileName)
					err = os.Rename(path.Join(fullPath, fileName), aimPath)
					if err != nil {
						logChan <- fmt.Sprintf("移动错误：%v", err)
					}
				}
			}
		} else if file.IsDir() {
			moveFilesInCatalogs(fullPath)
		}
	}
}

// getFileNum 能够返回给出目录下所有前缀为image后缀为.jpg的文件数量
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

// getIDList 能够读取给出表格文件的第一个sheet的第三列，
// 并从第三列的第三个单元格开始将整列输出成为一个切片
func getIDList(fileName string) ([]string, error) {

	f, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	cols, err := f.GetCols(f.GetSheetList()[0])
	if err != nil {
		return nil, err
	}
	idList := cols[2][2:]
	return idList, nil
}

// imageXXXXRename 能够将文件夹下以image开头并以.jpg结尾的文件根据表格文件重命名为对应身份证号
func imageXXXXRename(imagesPath, excelFile string) {
	logChan <- "开始运行..."
	fileNum, err := getFileNum(imagesPath)
	if err != nil {
		logChan <- fmt.Sprint("打开文件夹失败：", err)
	}
	idList, err := getIDList(excelFile)
	if err != nil {
		logChan <- fmt.Sprint("打开表格文件失败：", err)
		return
	}

	if len(idList) != fileNum {
		logChan <- fmt.Sprintf("文件数目与身份证数量不等！文件数目：%d，身份证数量：%d", fileNum, len(idList))
		return
	}

	for i := 1; i < fileNum+1; i++ {
		file := fmt.Sprintf("image%04d.jpg", i)
		fullName := path.Join(imagesPath, file)
		err = os.Rename(fullName, path.Join(imagesPath, idList[i-1]+".jpg"))
		if err != nil {
			logChan <- fmt.Sprint("重命名失败：", fullName)
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

// getMap 能够读取表格中的工号、姓名和身份证号并制作成映射表
func getMap(fileName string) (map[string]string, error) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		errText := fmt.Sprintf("打开 %s 失败！", fileName)
		return nil, errors.New(errText)
	}
	cols, err := f.GetCols(f.GetSheetList()[0])
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

// folderRename 能够将给出的项目文件夹下的特定格式的文件夹名通过excel重命名成对应的身份证号
func folderRename(projectPath, excelFile string) {
	logChan <- "开始运行..."

	dict, err := getMap(excelFile) // 读取表格
	if err != nil {
		logChan <- fmt.Sprint("读取表格错误：", err)
		return
	}

	dirArr, err := readDirs(projectPath) // 读取目录
	if err != nil {
		logChan <- fmt.Sprint("读取项目文件夹失败：", err)
		return
	}
	for _, a := range dirArr {
		dirName := path.Base(a)
		dirNameClean := strings.TrimSpace(dirName)
		dirNameClean = strings.ReplaceAll(dirName, " ", "")
		id, exist := dict[dirNameClean]
		if !exist {
			logChan <- fmt.Sprintf("文件夹 %s 在表格中不存在！如果该文件夹非客户提交的姓名，请忽略此条消息。", a)
			continue
		}
		a = path.Clean(a)
		err := os.Rename(a, path.Join(path.Dir(a), id))
		if err != nil {
			logChan <- fmt.Sprint("重命名失败！：", err)
		}
	}
}

// scanAndRename 能够递归查找文件夹并重命名
func scanAndRename(dirPath string, menu *[]string) error {
	dirname := path.Base(dirPath)
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		logChan <- fmt.Sprint("读取文件夹失败：", dirPath)
		return err
	}
	for _, entry := range entries {
		fullPath := path.Join(dirPath, entry.Name())
		if entry.IsDir() && entry.Name() == "目录" {
			*menu = append(*menu, fullPath)
			continue
		}
		if entry.IsDir() {
			if err := scanAndRename(fullPath, menu); err != nil {
				return err
			}
		} else {
			fileName := entry.Name()
			ext := path.Ext(fileName)
			fileName = strings.TrimSuffix(fileName, ext)
			newName := fmt.Sprintf("%s-%s(1%s", dirname, fileName, ext)
			if err := os.Rename(fullPath, path.Join(dirPath, newName)); err != nil {
				logChan <- fmt.Sprint("重命名失败！", fullPath)
			}
		}
	}
	return nil
}

// fileRename 能够将给出的文件夹内文件的名字命名成上级文件夹名+一定规则
func fileRename(projectPath string) {
	logChan <- "开始运行..."
	menu := &[]string{}
	if err := scanAndRename(projectPath, menu); err != nil {
		logChan <- fmt.Sprint("错误：", err)
	}

	// 单独处理目录内文件
	for _, menuPath := range *menu {
		menuCounter := 3
		files, err := ioutil.ReadDir(menuPath)
		if err != nil {
			logChan <- fmt.Sprint("目录读取失败：", menuPath)
			continue
		}
		for _, file := range files {
			menuCounter++
			filePath, fileName := menuPath, file.Name()
			ext := path.Ext(fileName)
			fullPath := path.Join(filePath, fileName)
			id := path.Dir(filePath)
			_, id = path.Split(id)
			newName := fmt.Sprintf("%s-0(%d%s", id, menuCounter, ext)
			newFull := path.Join(filePath, newName)

			if err := os.Rename(fullPath, newFull); err != nil {
				logChan <- fmt.Sprint("重命名失败：", fullPath)
			}
		}
	}
}

// getDirList 递归获取全部文件夹，返回列表
func getDirList(dirPath string) ([]string, error) {
	dirList := []string{}
	dirs, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, dir := range dirs {
		if dir.IsDir() {
			fullPath := path.Join(dirPath, dir.Name())
			sonList, err := getDirList(fullPath)
			if err != nil {
				return nil, err
			}
			dirList = append(dirList, sonList...)
			dirList = append(dirList, fullPath)
		}
	}
	return dirList, nil
}

// getDirMap 能够生成文件夹名和路径名的映射类型
func getDirMap(dirPath string) (map[string]string, error) {
	dirList, err := getDirList(dirPath)
	if err != nil {
		return nil, err
	}
	dirMap := map[string]string{}
	for _, dir := range dirList {
		name := path.Base(dir)
		dirMap[name] = dir
	}
	return dirMap, nil
}

// getBindings 读取封皮文件夹，返回其内.jpg文件列表
func getBindings(dirPath string) ([]string, error) {
	bindings := []string{}
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		return nil, err
	}
	for _, file := range files {
		fileName := file.Name()
		if !file.IsDir() && path.Ext(fileName) == ".jpg" {
			bindings = append(bindings, path.Join(dirPath, fileName))
		}
	}
	return bindings, nil
}

// moveImage 能够将封皮和档案袋转移至目录文件夹
func moveImage(imgDir, projectDir, typeOfFile string) {
	logChan <- "开始处理..."
	jpgList, err := getBindings(imgDir) // 读取jpg列表
	if err != nil {
		logChan <- fmt.Sprint("读取封皮文件夹失败：", err)
	}
	dirMap, err := getDirMap(projectDir) // 读取目标文件夹
	if err != nil {
		logChan <- fmt.Sprint("读取项目文件夹失败：", err)
		return
	}

	for _, jpg := range jpgList {
		_, jpgName := path.Split(jpg)
		ext := path.Ext(jpgName)                   // 获取扩展名
		jpgName = strings.TrimSuffix(jpgName, ext) // 获取文件名（身份证号）
		aimDir, ok := dirMap[jpgName]
		if !ok { //
			logChan <- fmt.Sprint("目标路径不存在！", jpg)
			continue
		}
		aimDir = path.Join(aimDir, "目录")
		if _, err := os.Stat(aimDir); os.IsNotExist(err) {
			logChan <- fmt.Sprint(aimDir, " 不存在！取消移动文件", jpg)
			continue
		}
		newPath := ""
		if typeOfFile == "x" {
			previousImgNum := 0
			previousImgs, err := ioutil.ReadDir(aimDir)
			if err != nil {
				logChan <- fmt.Sprint("读取目标文件夹内文件数量失败", aimDir)
				return
			}
			for _, previousImg := range previousImgs {
				if !previousImg.IsDir() && strings.HasSuffix(previousImg.Name(), ".jpg") {
					previousImgNum++
				}
			}
			newPath = path.Join(aimDir, fmt.Sprintf("%s-0(%d%s", jpgName, previousImgNum+1, ext))
		} else {
			newPath = path.Join(aimDir, jpgName+"-0("+typeOfFile+ext)
		}
		if _, err := os.Stat(newPath); !os.IsNotExist(err) {
			logChan <- fmt.Sprint(newPath, " 已存在！取消移动文件")
			continue
		}
		err = os.Rename(jpg, newPath)
		if err != nil {
			logChan <- fmt.Sprint("移动文件失败：", err)
		}
	}
}

// getAxis 通过x和y生成excel坐标文本
func getAxis(x, y int) string {
	return fmt.Sprintf("%s%d", getX(x), y+1)
}

// getX 生成横坐标字母
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

// trimTool 能够将每一页的每个单元格前后空白去除
func trimTool(fileName string) {

	f, err := excelize.OpenFile(fileName)
	if err != nil {
		logChan <- fmt.Sprint("打开文件 ", fileName, " 失败：非标准xlsx文档！")
		return
	}

	// 获取文件中的sheet列表
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		logChan <- fmt.Sprint("空文件！")
		return
	}

	// 读取文件中的每个sheet
	for i, sheet := range sheetList {
		rows, err := f.GetRows(sheet)
		if err != nil {
			logChan <- fmt.Sprintf("第%d个sheet、读取失败！", i+1)
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
				logChan <- fmt.Sprintf("第%d个sheet、单元格%s：“%s”-已修复", i+1, axis, col)
			}
		}
	}
	// 保存文件
	if err := f.Save(); err != nil {
		logChan <- fmt.Sprint("表格文件保存失败！请检查文件是否被占用！")
	}
}

// idExtand 将输入15位身份证号扩展成18位
func idExtend(fif string) (string, error) {
	if len(fif) != 15 {
		return "", errors.New("Wrong Length Error: " + fif)
	}
	if fif[6:8] == "00" {
		return "", errors.New("Wrong Years Error: " + fif)
	}

	a := [17]int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
	b := [11]string{"1", "0", "X", "9", "8", "7", "6", "5", "4", "3", "2"}

	fif = fif[:6] + "19" + fif[6:]
	sum := 0
	for i, v := range fif {
		n, _ := strconv.Atoi(string(v))
		sum += a[i] * n
	}
	return fif + b[sum%11], nil
}

// excelIDExtand 能够将输入的excel文件中的全部15位数字转换成18位身份证号
func excelIDExtand(fileName string) {
	f, err := excelize.OpenFile(fileName)
	if err != nil {
		logChan <- fmt.Sprint("打开文件 ", fileName, " 失败: 非标准xlsx文档！")
		return
	}

	// 获取文件中的sheet列表
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		logChan <- fmt.Sprint("空文件！")
		return
	}

	// 读取文件中的每个sheet
	for i, sheet := range sheetList {
		rows, err := f.GetRows(sheet)
		if err != nil {
			logChan <- fmt.Sprintf("第%d个sheet、读取失败！", i+1)
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
				col = strings.TrimSpace(col)
				if len(col) != 15 || !isDigit(col) {
					continue
				}
				extended, err := idExtend(col)
				axis := getAxis(x, y)
				if err != nil {
					logChan <- fmt.Sprintf("第%d个sheet、单元格%s：“%s”有误：%v", i+1, axis, col, err)
					continue
				}
				f.SetCellValue(sheet, axis, extended)
				logChan <- fmt.Sprintf("第%d个sheet、单元格%s：“%s”扩展为%s", i+1, axis, col, extended)
			}
		}
	}
	err = f.Save()
	if err != nil {
		logChan <- fmt.Sprint("文件保存失败！", fileName)
	}
}

// xlsxCheck 能够检测案卷目录是否符合标准
func xlsxCheck(fileName string) {

	f, err := excelize.OpenFile(fileName)
	if err != nil {
		logChan <- fmt.Sprint("打开文件", fileName, "失败: 非标准xlsx文档！")
		return
	}
	sheetList := f.GetSheetList()
	if len(sheetList) == 0 {
		logChan <- fmt.Sprint("空文件！")
		return
	}
	for i, sheet := range sheetList {
		rows, err := f.GetRows(sheet)
		if err != nil {
			logChan <- fmt.Sprintf("第%d个sheet、读取失败！", i+1)
			continue
		}
		checkAll(i+1, rows)
	}
}

// checkAll 能够检查一个二维数组形式的sheet数据是否符合一定标准
func checkAll(index int, rows [][]string) {
	preText := fmt.Sprintf("第%d个sheet、", index)
	if len(rows) < 3 {
		logChan <- fmt.Sprint(preText + "数据过少！请保证在有两行标题的情况下至少有一行内容")
		return
	}
	if strings.TrimSpace(rows[0][0]) != "国有企业退休人员人事档案移交案卷目录" {
		logChan <- fmt.Sprint(preText + "第1行，标题不正确！请检查标题是否为“国有企业退休人员人事档案移交案卷目录”")
		return
	}
	title := rows[1]
	titleError := ""
	if len(title) != 14 {
		logChan <- fmt.Sprint(preText + "文档列数不正确！请检查此sheet是否按顺序包含以下14列：\n序号、姓名、身份证号码、性别、出生日期、退休手续办理时间、生存状态、死亡日期、所属单位、户籍所在地、页数、档案局档号、保管期限、备注")
		return
	}
	switch {
	case title[1] != "姓名":
		titleError += preText + "姓名不在第二列！\n"
		fallthrough
	case title[2] != "身份证号码":
		titleError += preText + "身份证号码不在第三列！\n"
		fallthrough
	case title[3] != "性别":
		titleError += preText + "性别不在第四列！\n"
		fallthrough
	case title[4] != "出生日期":
		titleError += preText + "出生日期不在第五列！\n"
		fallthrough
	case title[5] != "退休手续办理时间":
		titleError += preText + "退休手续办理时间不在第二列！\n"
	}
	if titleError != "" {
		logChan <- fmt.Sprint(titleError)
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
				logChan <- fmt.Sprint(preText + fmt.Sprintf("第%d列，单元格内文本前后包含空格！", i+1))
				// 比较前去掉前后空格
				line[i] = strings.TrimSpace(l)
			}
		}
		id := line[2]
		// check id
		if id == "" {
			logChan <- fmt.Sprint(preText + "证件号码为空！")
			continue
		}
		if len(id) != 18 {
			logChan <- fmt.Sprint(preText + "证件号码位数不正确！")
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
			logChan <- fmt.Sprint(preText + "性别为空！")
		} else {
			trueSex := "男"
			trueSexNum, _ := strconv.Atoi(string(id[16]))
			if trueSexNum%2 == 0 {
				trueSex = "女"
			}
			if trueSex != sex {
				logChan <- fmt.Sprint(preText + "性别不正确！")
			}
		}
		// check birthday
		if birthday == "" {
			logChan <- fmt.Sprint(preText + "出生日期为空！")
		} else if len(birthday) != 8 {
			logChan <- fmt.Sprint(preText + "出生日期位数不正确！")
		} else if birthday != string(id[6:14]) {
			logChan <- fmt.Sprint(preText + "出生日期与身份证件不符！")
		}
		// check retirement
		if retirement == "" {
			logChan <- fmt.Sprint(preText + "退休手续办理时间为空！")
		} else if len(retirement) != 8 {
			logChan <- fmt.Sprint(preText + "退休手续办理时间位数不正确！")
		} else if !isDigit(retirement) {
			logChan <- fmt.Sprint(preText + "退休手续办理时间包含非数字内容！")
		}
		// check others
		for k, v := range otherInfo {
			if v == "" {
				logChan <- fmt.Sprint(preText + k + "为空！")
			}
		}
	}
	if len(wrongIDs) != 0 {
		preText = fmt.Sprintf("第%d个sheet、", index)
		for _, wID := range wrongIDs {
			logChan <- fmt.Sprint(preText + fmt.Sprintf("第%s行，", wID[1]) + wID[0] + "身份证号码格式不正确，请确认")
		}
	}
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

func logOutput(parent fyne.Window) {
	if logTxt == "" {
		dialog.ShowInformation("提示", "日志为空！无需导出。", parent)
		return
	}
	logFile, err := os.OpenFile("tstools.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
	if err != nil {
		dialog.ShowInformation("错误", "日志文件打开失败！请检查“tstools.log”文件是否被占用！", parent)
		return
	}
	defer logFile.Close()

	_, err = logFile.WriteString(logTxt + "\r\n")
	if err != nil {
		dialog.ShowInformation("错误", "日志文件写入失败！请检查“tstools.log”文件是否被占用！", parent)
		return
	}
	dialog.ShowInformation("提示", "导出成功！", parent)
	logTxt = ""
}
