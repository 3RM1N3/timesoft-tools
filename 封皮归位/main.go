package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
)

func main() {
	projectDir := ""
	fmt.Print("输入项目文件夹的绝对路径或其与程序所在文件夹的相对路径并按回车：")
	fmt.Scanln(&projectDir)
	dirMap := getDirMap(projectDir) // 读取目标文件夹
	bindingsDir := ""
	fmt.Print("输入包含所有封皮pdf文件的文件夹的绝对路径或其与程序所在文件夹的相对路径并按回车：")
	fmt.Scanln(&bindingsDir)
	pdfList := getBindings(bindingsDir) // 读取pdf列表

	for _, pdf := range pdfList {
		_, pdfName := path.Split(pdf)
		ext := path.Ext(pdfName)                   // 获取扩展名
		pdfName = strings.TrimSuffix(pdfName, ext) // 获取文件名（身份证号）
		aimDir, ok := dirMap[pdfName]
		if !ok { //
			log.Println("目标路径不存在！", pdf)
			continue
		}
		aimDir = path.Join(aimDir, "目录")
		if _, err := os.Stat(aimDir); os.IsNotExist(err) {
			log.Println(aimDir, "不存在！取消移动文件", pdf)
			continue
		}
		newPath := path.Join(aimDir, pdfName+"-0(2"+ext)
		if _, err := os.Stat(newPath); !os.IsNotExist(err) {
			log.Println(newPath, "已存在！取消移动文件")
			continue
		}
		os.Rename(pdf, newPath)
	}
}

// 获取全部文件夹，返回列表
func getDirList(dirPath string) []string {
	dirList := []string{}
	dirs, _ := ioutil.ReadDir(dirPath)
	for _, dir := range dirs {
		if dir.IsDir() {
			fullPath := path.Join(dirPath, dir.Name())
			sonList := getDirList(fullPath)
			dirList = append(dirList, sonList...)
			dirList = append(dirList, fullPath)
		}
	}
	return dirList
}

// getDirMap 能够生成文件夹名和路径名的映射类型
func getDirMap(dirPath string) map[string]string {
	dirList := getDirList(dirPath)
	dirMap := map[string]string{}
	for _, dir := range dirList {
		name := path.Base(dir)
		dirMap[name] = dir
	}
	return dirMap
}

// 读取封皮文件夹，返回其内.pdf文件列表
func getBindings(dirPath string) []string {
	bindings := []string{}
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Fatal("打开文件夹失败！", dirPath)
	}
	for _, file := range files {
		fileName := file.Name()
		if !file.IsDir() && path.Ext(fileName) == ".jpg" {
			bindings = append(bindings, path.Join(dirPath, fileName))
		}
	}
	return bindings
}
