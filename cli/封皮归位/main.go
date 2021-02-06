package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"
)

func main() {
	usage := `欢迎使用时源科技封皮归位工具！

    本程序可将输入的装有 身份证号.jpg的封皮文件 的目录下全部 身份证号.jpg 文件移动到输入的项目目录中对应的身份证号文件夹的“目录”内。

    *注意！
    1. 对于包含jpg文件的文件夹，程序不会递归查找子文件夹
    2. 对于项目文件夹，程序将递归查找子文件夹
    3. *该程序会原地操作文件，故运行后将不可逆！请提前进行文件备份！*

    Usage:  根据提示输入两个文件夹的路径即可开始。

		如：
		    输入包含所有身份证号.jpg封皮文件的文件夹的绝对路径或其与程序所在文件夹的相对路径并按回车：D:\project\封皮
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
	bindingsDir := ""
	fmt.Print("输入包含所有身份证号.jpg封皮文件的文件夹的绝对路径或其与程序所在文件夹的相对路径并按回车：")
	fmt.Scanln(&bindingsDir)
	jpgList, err := getBindings(bindingsDir) // 读取jpg列表
	if err != nil {
		log.Println("读取封皮文件夹失败：", err)
	}
	projectDir := ""
	fmt.Print("输入项目文件夹的绝对路径或其与程序所在文件夹的相对路径并按回车：")
	fmt.Scanln(&projectDir)
	dirMap, err := getDirMap(projectDir) // 读取目标文件夹
	if err != nil {
		log.Println("读取项目文件夹失败：", err)
		return
	}

	for _, jpg := range jpgList {
		_, jpgName := path.Split(jpg)
		ext := path.Ext(jpgName)                   // 获取扩展名
		jpgName = strings.TrimSuffix(jpgName, ext) // 获取文件名（身份证号）
		aimDir, ok := dirMap[jpgName]
		if !ok { //
			log.Println("目标路径不存在！", jpg)
			continue
		}
		aimDir = path.Join(aimDir, "目录")
		if _, err := os.Stat(aimDir); os.IsNotExist(err) {
			log.Println(aimDir, "不存在！取消移动文件", jpg)
			continue
		}
		newPath := path.Join(aimDir, jpgName+"-0(2"+ext)
		if _, err := os.Stat(newPath); !os.IsNotExist(err) {
			log.Println(newPath, "已存在！取消移动文件")
			continue
		}
		err = os.Rename(jpg, newPath)
		if err != nil {
			log.Println("移动文件失败：", err)
		}
	}
}

// 获取全部文件夹，返回列表
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

// 读取封皮文件夹，返回其内.jpg文件列表
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
