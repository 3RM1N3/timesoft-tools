package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func fileRenamePage(win fyne.Window, logVbox *fyne.Container) fyne.CanvasObject {
	/*usage := `欢迎使用时源科技文件重命名工具

	  本程序可将 输入的项目文件夹 下全部文件及全部子文件重命名为 文件上级文件夹名 + 原文件名 + (1 ；
	  若文件夹名为“目录”，则将此文件夹下文件重命名为 “目录”上级文件夹名 + 0( + 以罗马数字3开始的序号

	  *注意！
	  1. 程序可递归查找全部的文件夹及子文件夹
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
			fileRename(entryChooseDir.Text) // 主要功能实现
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
		logVbox,
	)
}
