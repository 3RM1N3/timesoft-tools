package main

import (
	"log"
	"os"
	"strings"
	"time"

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
	usage := `
使用帮助：

    本软件可根据表格内数据实现基本的分件功能。

*注意事项：

1. 在*.xlsx文件所选中的Sheet中，不应包含标题，
   且档号应在第一列，页号应在第二列；
2. 所有操作均在原地操作，请提前进行文件备份。


Author: 3RM1N3@时源科技
`
	runable := true
	sheetList := []string{} // sheet列表
	selectedSheet := ""     // 选定的sheet页
	myApp := app.New()
	window := myApp.NewWindow("分件") // 新建窗口

	statusBar := widget.NewLabel("选择表格文件、数据所在Sheet与项目文件夹后点击运行开始")

	tablePathEntry := widget.NewEntry() // 用于显示所选择xlsx文件的文本框
	folderEntry := widget.NewEntry()    // 用于显示所选文件夹的文本框
	sheetSelect := new(widget.Select)
	sheetSelect.OnChanged = func(s string) {
		selectedSheet = s
		log.Println("选择Sheet：", selectedSheet)
		statusBar.SetText("请选择项目文件夹")
	}
	sheetSelect.Disable()

	usageLabel := widget.NewLabel(usage)

	chooseTableButton := widget.NewButton("...", func() { // 选择表格文件的按钮
		dialog.ShowFileOpen(func(uc fyne.URIReadCloser, e error) {
			if e != nil || uc == nil {
				log.Println(e)
				return
			}
			statusBar.SetText("请选择需要使用的Sheet页")
			defer uc.Close()

			f, err := excelize.OpenReader(uc)
			if err != nil {
				dialog.ShowInformation("错误", "非标准xlsx文档！", window)
				return
			}
			sheetList = f.GetSheetList()
			filePath := uc.URI().String()
			tablePathEntry.SetText(strings.TrimPrefix(filePath, "file://"))
			sheetSelect.Enable()
			sheetSelect.Options = sheetList
		}, window)
	})

	chooseFolderButton := widget.NewButton("...", func() { // 选择文件夹按钮
		dialog.ShowFolderOpen(func(lu fyne.ListableURI, e error) {
			if e != nil || lu == nil {
				return
			}
			if tablePathEntry.Text == "" {
				statusBar.SetText("请选择表格文件")
			} else if selectedSheet == "" {
				statusBar.SetText("请选择需要使用的Sheet页")
			} else {
				statusBar.SetText("单击运行以开始")
			}
			folderEntry.SetText(strings.TrimPrefix(lu.String(), "file://"))
		}, window)
	})

	runButton := widget.NewButton("运行", func() { // 运行按钮
		if tablePathEntry.Text == "" || folderEntry.Text == "" || selectedSheet == "" {
			dialog.ShowInformation("错误", "尚有字段为空！", window)
			return
		}
		cnf := dialog.NewConfirm("注意", "此操作执行后将不可逆！\n请提前进行文件备份并确认您已知悉程序运行后果！确定仍要继续吗？", func(r bool) {
			if !r {
				return
			}
			runable = false
			statusBar.SetText("处理中...请勿进行其他操作")
			do(tablePathEntry.Text, selectedSheet, folderEntry.Text) // 开始处理
			tablePathEntry.SetText("")
			folderEntry.SetText("")
			selectedSheet = ""
			sheetSelect.Options = []string{}
			runable = true
			statusBar.SetText("完毕")
			dialog.ShowInformation("提示", "处理完毕", window)
		}, window)
		cnf.SetDismissText("取消")
		cnf.SetConfirmText("确定")
		cnf.Show()
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
		runButton,
		usageLabel,
	)
	go func() {
		logFile, err := os.OpenFile("分件日志.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0766)
		if err != nil {
			dialog.ShowInformation("提示", "日志文件打开失败，请检查文件是否被占用；\n接下来的操作会产生效果但不会在日志中体现。", window)
			return
		}
		log.SetOutput(logFile) // 将文件设置为log输出的文件
		//log.SetFlags(log.LstdFlags | log.Lshortfile | log.LUTC)
	}()
	go func() {
		for {
			time.Sleep(200 * time.Millisecond)
			if !runable {
				runButton.Disable()
				continue
			}
			runButton.Enable()
		}
	}()

	window.SetContent(
		container.NewBorder(
			nil,
			statusBar,
			nil,
			nil,
			ctnr,
		))
	window.Resize(fyne.NewSize(600, 500))
	window.ShowAndRun()
}
