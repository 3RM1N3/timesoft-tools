package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/360EntSecGroup-Skylar/excelize/v2"
	"github.com/flopp/go-findfont"
)

func init() {
	fontList := findfont.List() // 读取系统字体
	for _, font := range fontList {
		if strings.Contains(font, "simhei.ttf") { // 设置字体为黑体
			os.Setenv("FYNE_FONT", font)
			break
		}
	}
}

func main() {
	sheetList := []string{} // sheet列表
	selectedSheet := ""     // 选定的sheet页
	myApp := app.New()
	window := myApp.NewWindow("分件") // 新建窗口

	tablePathEntry := widget.NewEntry() // 用于显示所选择xlsx文件的文本框
	folderEntry := widget.NewEntry()    // 用于显示所选文件夹的文本框
	sheetSelect := new(widget.Select)
	sheetSelect.OnChanged = func(s string) {
		selectedSheet = s
		log.Println(selectedSheet)
	}
	sheetSelect.Disable()

	chooseTableButton := widget.NewButton("...", func() { // 选择表格文件的按钮
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, e error) {
			if e != nil || uc == nil {
				log.Println("wrong1")
				return
			}
			defer uc.Close()

			f, err := excelize.OpenReader(uc)
			if err != nil {
				dialog.ShowInformation("错误", "非标准xlsx文档！", window)
				return
			}
			sheetList = f.GetSheetList()
			filePath := uc.URI().String()
			tablePathEntry.SetText(strings.TrimPrefix(filePath, "file://"))
			fmt.Println(sheetList)
			sheetSelect.Enable()
			sheetSelect.Options = sheetList
		}, window)
	})

	chooseFolderButton := widget.NewButton("...", func() { // 选择文件夹按钮
		dialog.ShowFolderOpen(func(lu fyne.ListableURI, e error) {
			if e != nil || lu == nil {
				return
			}
			folderEntry.SetText(strings.TrimPrefix(lu.String(), "file://"))
		}, window)
	})

	ctnr := container.NewVBox( // 容器
		widget.NewLabel("文件部分："),
		container.NewBorder(
			nil,
			nil,
			widget.NewLabel("1. 选择表格文件："),
			chooseTableButton,
			tablePathEntry,
		),
		container.NewBorder(
			nil,
			nil,
			widget.NewLabel("2. 选择使用的Sheet："),
			nil,
			sheetSelect,
		),
		widget.NewSeparator(),
		widget.NewLabel("文件夹部分："),
		container.NewBorder(
			nil,
			nil,
			widget.NewLabel("1. 选择项目文件夹："),
			chooseFolderButton,
			folderEntry,
		),
		widget.NewSeparator(),
		container.NewMax(
			widget.NewButton("运行", func() { // 运行按钮
				if tablePathEntry.Text == "" || folderEntry.Text == "" || selectedSheet == "" {
					dialog.ShowInformation("错误", "尚有字段为空！", window)
					return
				}
				do(tablePathEntry.Text, selectedSheet, folderEntry.Text) // 开始处理
			}),
		),
	)

	window.SetContent(ctnr)
	window.Resize(fyne.NewSize(600, 400))
	window.ShowAndRun()
}
