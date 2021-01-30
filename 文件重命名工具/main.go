package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var menu []string

func main() {
	basedir := ""
	fmt.Print("请复制待处理文件夹路径粘贴并按回车键开始运行：")
	fmt.Scanln(&basedir)
	scanAndRename(basedir)

	// 单独处理 目录

	for _, menuPath := range menu {
		menuCounter := 2
		files, _ := ioutil.ReadDir(menuPath)
		for _, file := range files {
			menuCounter++
			filePath, fileName := menuPath, file.Name()
			ext := path.Ext(fileName)
			fullPath := path.Join(filePath, fileName)
			id := path.Dir(filePath) // 此处可能需要再次嵌套path.Dir()
			_, id = path.Split(id)
			newName := fmt.Sprintf("%s-0(%d%s", id, menuCounter, ext)
			newFull := path.Join(filePath, newName)

			os.Rename(fullPath, newFull)
		}

	}
}

func scanAndRename(dirPath string) error {

	dirname := path.Base(dirPath)
	//fileNum := 0

	// start to rename
	entries, _ := ioutil.ReadDir(dirPath)
	for _, entry := range entries {
		fullPath := path.Join(dirPath, entry.Name())
		if entry.IsDir() && entry.Name() == "目录" {
			menu = append(menu, fullPath)
			continue
		}
		if entry.IsDir() {
			scanAndRename(fullPath)
		} else {
			fileName := entry.Name()
			ext := path.Ext(fileName)
			fileName = strings.TrimSuffix(fileName, ext)
			newName := fmt.Sprintf("%s-%s(1%s", dirname, fileName, ext)
			os.Rename(fullPath, path.Join(dirPath, newName))
		}
	}

	return nil
}
