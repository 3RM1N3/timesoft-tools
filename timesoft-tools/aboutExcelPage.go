package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func aboutExcelPage(win fyne.Window, logVbox *fyne.Container) fyne.CanvasObject {
	entryChoose := widget.NewEntry()
	fileChooser := widget.NewButton("...", func() {
		fd := dialog.NewFileOpen(func(f fyne.URIReadCloser, e error) {
			if e != nil || f == nil {
				return
			}
			entryChoose.SetText(strings.TrimPrefix(f.URI().String(), "file://"))
		}, win)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".xlsx"}))
		fd.Show()
	})

	options := map[string]func(fileName string){
		"Trim Space工具": trimTool,
		"身份证号码扩展":      excelIDExtand,
		"案卷目录格式检测":     xlsxCheck,
	}

	optionList := []string{}
	for k := range options {
		optionList = append(optionList, k)
	}

	choosed := ""
	selecter := widget.NewSelect(optionList, func(s string) {
		choosed = s
	})
	outputButton := widget.NewButton("导出日志", func() {
		logOutput(win)
	})
	runButton := widget.NewButton("运行", func() {
		if busy {
			dialog.ShowInformation("错误", "请等待其他程序结束后重试。", win)
			return
		} else if entryChoose.Text == "" && choosed == "" {
			dialog.ShowInformation("错误", "请仔细阅读说明！尚未选择表格文件与操作类型。", win)
			return
		} else if entryChoose.Text == "" {
			dialog.ShowInformation("错误", "尚未选择表格文件。", win)
			return
		} else if choosed == "" {
			dialog.ShowInformation("错误", "尚未选择操作类型。", win)
			return
		}
		cnf := dialog.NewConfirm("注意", "此操作执行后将不可逆！请提前进行文件备份并再次确认您已知悉程序运行后果！确定仍要继续吗？", func(r bool) {
			if !r {
				return
			}
			busy = true
			options[choosed](entryChoose.Text) // 主要功能实现
			entryChoose.SetText("")
			busy = false
			logChan <- "操作完成"
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
				widget.NewLabel("选择文件(*.xlsx)："),
				fileChooser,
				entryChoose,
			),
			container.NewGridWithColumns(
				2,
				container.NewBorder(
					nil,
					nil,
					widget.NewLabel("选择操作类型："),
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
		nil,
		nil,
		nil,
		logVbox,
	)
}

/*container.NewCenter(widget.NewLabel(`欢迎使用时源科技 Excel 相关工具


Trim Space工具：可以将输入文件内所有表格的所有单元格前后的空白去除

身份证号码扩展：可以将输入文件内所有的15位数字扩展为18位身份证号

案卷目录格式检测：可以检查输入文件文件是否符合案卷目录规范


Usage:  选择需要处理的表格文件和操作类型后点击运行即可开始。`)),
	)
}
*/
