package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func folderRenamePage(win fyne.Window, logText *fyne.Container) fyne.CanvasObject {
	/*usage := `欢迎使用时源科技文件夹重命名工具

	  本软件可将输入的项目文件夹根据 档案号+姓名 在所选xlsx文件的第一个Sheet中匹配数据并批量将文件夹重命名为对应身份证号。

	  *注意！
	  1. 所选xlsx文件的第一个Sheet中应有两行标题
	  2. 工号/档案号 应在第三列
	  3. 姓名 应在第四列
	  4. 身份证号 应在第七列
	  5. *该程序会原地操作文件，故运行后将不可逆！请提前进行文件备份！*


	  Author: 3RM1N3@时源科技
	  E-mail: wangyu7439@hotmail.com`*/

	statusBar := widget.NewLabel("依次选择表格文件、项目文件夹后点击运行以开始")

	entryChooseFile := widget.NewEntry()
	entryChooseDir := widget.NewEntry()

	buttonChooseFile := widget.NewButton("...", func() {
		fd := dialog.NewFileOpen(func(f fyne.URIReadCloser, err error) {
			if err != nil || f == nil {
				return
			}
			fileName := f.URI().String()
			if strings.HasPrefix(fileName, "file://") {
				fileName = strings.TrimPrefix(fileName, "file://")
			} else {
				fileName = strings.TrimPrefix(fileName, "file:\\\\")
			}
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
			folderRename(entryChooseDir.Text, entryChooseFile.Text) // 主要功能实现
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
				widget.NewLabel("选择项目文件夹："),
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
