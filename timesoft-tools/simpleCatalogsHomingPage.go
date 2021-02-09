package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func simpleCatalogsHomingPage(win fyne.Window, logText *fyne.Container) fyne.CanvasObject {

	statusBar := widget.NewLabel("依次选择简易目录文件夹和项目文件夹后点击运行以开始")

	entryChooseImgDir := widget.NewEntry()
	entryChooseDir := widget.NewEntry()

	buttonChooseImgDir := widget.NewButton("...", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil || list == nil {
				return
			}
			path := strings.TrimPrefix(list.String(), "file://")
			entryChooseImgDir.SetText(path)
			if entryChooseImgDir.Text == "" {
				statusBar.SetText("请选择项目文件夹")
				return
			}
			statusBar.SetText("点击运行以开始")
		}, win)
	})
	buttonChooseDir := widget.NewButton("...", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil || list == nil {
				return
			}
			path := strings.TrimPrefix(list.String(), "file://")
			entryChooseDir.SetText(path)
			if entryChooseImgDir.Text == "" {
				statusBar.SetText("请选择图片文件夹")
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
		if entryChooseImgDir.Text == "" && entryChooseDir.Text == "" {
			dialog.ShowInformation("错误", "请仔细阅读说明！尚未选择简易目录文件夹与项目文件夹！", win)
			return
		}
		if entryChooseImgDir.Text == "" {
			dialog.ShowInformation("错误", "尚未选择图片文件夹！", win)
			return
		}
		if entryChooseDir.Text == "" {
			dialog.ShowInformation("错误", "尚未选择项目文件夹！", win)
			return
		}
		cnf := dialog.NewConfirm("注意", "此操作执行后将不可逆！请提前进行文件备份并再次确认您已知悉程序运行后果！确定仍要继续吗？", func(r bool) {
			if !r {
				return
			}
			busy = true
			statusBar.SetText("处理中...")
			simpleCatalogsHoming(entryChooseImgDir.Text, entryChooseDir.Text) // 主要功能实现
			entryChooseDir.SetText("")
			entryChooseImgDir.SetText("")
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
				widget.NewLabel("选择多文件简易目录文件夹："),
				buttonChooseImgDir,
				entryChooseImgDir,
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
