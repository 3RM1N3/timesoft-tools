package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func moveImagePage(win fyne.Window, logText *fyne.Container) fyne.CanvasObject {
	/*usage := `欢迎使用时源科技封皮&档案袋归位工具

	  本程序可将输入的装有 身份证号.jpg 的文件夹下全部 身份证号.jpg 文件移动到输入的项目文件夹中对应的身份证号文件夹的“目录”内。

	  *注意！
	  1. 对于包含jpg文件的文件夹，程序不会递归查找子文件夹
	  2. 对于项目文件夹，程序将递归查找子文件夹
	  3. *该程序会原地操作文件，故运行后将不可逆！请提前进行文件备份！*


	  Author: 3RM1N3@时源科技
	  E-mail: wangyu7439@hotmail.com`*/

	statusBar := widget.NewLabel("依次选择图片文件夹、项目文件夹和档案类型后点击运行以开始")

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

	typeOfFile := ""
	selecter := widget.NewSelect([]string{"封皮", "档案袋"}, func(s string) {
		if s == "封皮" {
			typeOfFile = "2"
			return
		}
		typeOfFile = "1"
	})

	outputButton := widget.NewButton("导出日志", func() {
		logOutput(win)
	})

	runButton := widget.NewButton("运行", func() { // 运行按钮
		if busy {
			dialog.ShowInformation("错误", "请等待其他程序结束后重试。", win)
			return
		}
		if entryChooseImgDir.Text == "" && entryChooseDir.Text == "" && typeOfFile == "" {
			dialog.ShowInformation("错误", "请仔细阅读说明！尚未选择图片文件夹、项目文件夹与档案类型！", win)
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
		if typeOfFile == "" {
			dialog.ShowInformation("错误", "尚未选择档案类型", win)
			return
		}
		cnf := dialog.NewConfirm("注意", "此操作执行后将不可逆！请提前进行文件备份并再次确认您已知悉程序运行后果！确定仍要继续吗？", func(r bool) {
			if !r {
				return
			}
			busy = true
			statusBar.SetText("处理中...")
			moveImage(entryChooseImgDir.Text, entryChooseDir.Text, typeOfFile) // 主要功能实现
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
				widget.NewLabel("选择图片文件夹："),
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
			container.NewGridWithColumns(
				2,
				container.NewBorder(
					nil,
					nil,
					widget.NewLabel("选择档案类型："),
					nil,
					selecter,
				),
				container.NewBorder(
					nil,
					nil,
					nil,
					outputButton,
					runButton,
				),
			),
		),
		statusBar,
		nil,
		nil,
		logText,
	)
}
