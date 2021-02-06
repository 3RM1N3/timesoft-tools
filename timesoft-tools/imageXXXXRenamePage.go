package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func imageXXXXRenamePage(win fyne.Window, logText *fyne.Container) fyne.CanvasObject {
	/*usage := `欢迎使用时源科技封皮&档案袋重命名工具

	  本软件可将 输入的装有 imageXXXX.jpg 的文件夹 下全部 imageXXXX.jpg 文件按顺序依次重命名为所选xlsx文档第一个Sheet中的第三列（身份证号）。

	  *注意！
	  1. 案卷目录.xlsx 的第一个Sheet中应有两行标题
	  2. 程序仅查找输入的文件夹，不会递归查找子文件夹
	  3. 身份证号 应在第三列
	  4. 应保证身份证号列的身份证号数量与读取到的jpg文件数量完全一致
	  5. *该程序会原地操作文件，故运行后将不可逆！请提前进行文件备份！*


	  Author: 3RM1N3@时源科技
	  E-mail: wangyu7439@hotmail.com`*/

	statusBar := widget.NewLabel("依次选择表格文件、包含需重命名图片的文件夹后点击运行以开始")

	entryChooseFile := widget.NewEntry()
	entryChooseDir := widget.NewEntry()

	buttonChooseFile := widget.NewButton("...", func() {
		fd := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
			if err != nil || f == nil {
				return
			}
			fileName := strings.TrimPrefix(f.URI().String(), "file://")
			entryChooseFile.SetText(fileName)
			if entryChooseDir.Text == "" {
				statusBar.SetText("请选择文件夹")
				return
			}
			statusBar.SetText("点击运行以开始")
		}, win)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		fd.Show()
	})
	buttonChooseDir := widget.NewButton("...", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil || list == nil {
				return
			}
			path := strings.TrimPrefix(list.String(), "file://")
			entryChooseDir.SetText(path)
			if entryChooseFile.Text == "" {
				statusBar.SetText("请选择表格文件")
				return
			}
			statusBar.SetText("点击运行以开始")
		}, win)
	})

	outputButton := widget.NewButton("导出日志", func() {
		logOutput(win)
	})

	runButton := widget.NewButton("运行", func() { // 运行按钮
		if busy {
			dialog.ShowInformation("错误", "请等待其他程序结束后重试。", win)
			return
		}
		if entryChooseFile.Text == "" && entryChooseDir.Text == "" {
			dialog.ShowInformation("错误", "请仔细阅读说明！尚未选择文件夹与表格文件！", win)
			return
		}
		if entryChooseFile.Text == "" {
			dialog.ShowInformation("错误", "尚未选择表格文件！", win)
			return
		}
		if entryChooseDir.Text == "" {
			dialog.ShowInformation("错误", "尚未选择文件夹！", win)
			return
		}

		cnf := dialog.NewConfirm("注意", "此操作执行后将不可逆！请提前进行文件备份并再次确认您已知悉程序运行后果！确定仍要继续吗？", func(r bool) {
			if !r {
				return
			}
			busy = true
			statusBar.SetText("处理中...")
			imageXXXXRename(entryChooseDir.Text, entryChooseFile.Text) // 主要功能实现
			entryChooseDir.SetText("")
			entryChooseFile.SetText("")
			busy = false
			logChan <- "操作完成"
			statusBar.SetText("完成")
		}, win)
		cnf.SetDismissText("取消")
		cnf.SetConfirmText("确定")
		cnf.Show()
	})

	return container.NewBorder(
		container.NewVBox(
			container.NewBorder(
				nil,
				nil,
				widget.NewLabel("选择表格文件(*.xlsx)："),
				buttonChooseFile,
				entryChooseFile,
			),
			container.NewBorder(
				nil,
				nil,
				widget.NewLabel("选择文件夹："),
				buttonChooseDir,
				entryChooseDir,
			),
			container.NewBorder(
				nil,
				nil,
				nil,
				outputButton,
				runButton,
			),
		),
		statusBar,
		nil,
		nil,
		logText,
	)
}
