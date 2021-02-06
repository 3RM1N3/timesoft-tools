package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func moveFilesInCatalogsPage(win fyne.Window, logText *fyne.Container) fyne.CanvasObject {
	/*usage := `欢迎使用时源科技目录转移工具

	  本软件可将选中的文件夹下全部名为“目录”的子文件夹内的文件转移至其上层文件夹。

	  *注意！
	  1. 程序会递归查找名称不为“目录”的全部子文件夹
	  2. *该程序会原地操作文件，故运行后将不可逆！请提前进行文件备份！*


	  Author: 3RM1N3@时源科技
	  E-mail: wangyu7439@hotmail.com`*/

	statusBar := widget.NewLabel("选择文件夹后点击运行以开始")

	entryChooseDir := widget.NewEntry()
	buttonChooseDir := widget.NewButton("...", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil || list == nil {
				return
			}
			path := strings.TrimPrefix(list.String(), "file://")
			entryChooseDir.SetText(path)
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
			moveFilesInCatalogs(entryChooseDir.Text) // 主要功能实现
			entryChooseDir.SetText("")
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
