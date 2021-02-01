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

var menu []string

func main() {
	usage := `欢迎使用时源科技文件重命名工具！

    本程序可将 输入的项目目录 下全部文件及全部子文件重命名为 文件上级文件夹名 + 原文件名 + (1 ；
    若文件夹名为“目录”，则将此文件夹下文件重命名为 “目录”上级文件夹名 + 0( + 以罗马数字3开始的序号

    *注意！
    1. 程序可递归查找全部的文件夹及子文件夹
    2. *该程序会原地操作文件，故运行后将不可逆！请提前进行文件备份！*

    Usage:  根据提示输入项目文件夹路径即可开始。

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
	basedir := ""
	fmt.Print("\n输入项目文件夹的绝对路径或其与程序所在文件夹的相对路径并按回车：")
	fmt.Scanln(&basedir)
	if err := scanAndRename(basedir); err != nil {
		log.Println("错误：", err)
	}

	// 单独处理 目录内文件
	for _, menuPath := range menu {
		menuCounter := 2
		files, err := ioutil.ReadDir(menuPath)
		if err != nil {
			log.Println("目录读取失败：", menuPath)
			continue
		}
		for _, file := range files {
			menuCounter++
			filePath, fileName := menuPath, file.Name()
			ext := path.Ext(fileName)
			fullPath := path.Join(filePath, fileName)
			id := path.Dir(filePath) // 此处可能需要再次嵌套path.Dir()
			_, id = path.Split(id)
			newName := fmt.Sprintf("%s-0(%d%s", id, menuCounter, ext)
			newFull := path.Join(filePath, newName)

			if err := os.Rename(fullPath, newFull); err != nil {
				log.Println("重命名失败：", fullPath)
			}
		}
	}
}

func scanAndRename(dirPath string) error {
	dirname := path.Base(dirPath)
	//fileNum := 0

	// start to rename
	entries, err := ioutil.ReadDir(dirPath)
	if err != nil {
		log.Println("读取文件夹失败：", dirPath)
		return err
	}
	for _, entry := range entries {
		fullPath := path.Join(dirPath, entry.Name())
		if entry.IsDir() && entry.Name() == "目录" {
			menu = append(menu, fullPath)
			continue
		}
		if entry.IsDir() {
			if err := scanAndRename(fullPath); err != nil {
				return err
			}
		} else {
			fileName := entry.Name()
			ext := path.Ext(fileName)
			fileName = strings.TrimSuffix(fileName, ext)
			newName := fmt.Sprintf("%s-%s(1%s", dirname, fileName, ext)
			if err := os.Rename(fullPath, path.Join(dirPath, newName)); err != nil {
				log.Println("重命名失败！", fullPath)
			}
		}
	}

	return nil
}
