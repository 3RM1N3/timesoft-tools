package main

import (
	"os"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/flopp/go-findfont"
)

func init() {
	fontList := findfont.List()
	for _, font := range fontList {
		if strings.Contains(font, "simhei.ttf") {
			os.Setenv("FYNE_FONT", font)
			break
		}
	}
}

var logChan = make(chan string, 1024)
var logTxt string
var busy = false

func main() {
	a := app.New()
	w := a.NewWindow("Timesoft Tools v0.4.3")
	var logVbox = *container.NewVBox()

	w.SetContent(container.NewAppTabs(
		container.NewTabItem("Excel操作", aboutExcelPage(w, &logVbox)),
		container.NewTabItem("文件夹重命名", folderRenamePage(w, &logVbox)),
		container.NewTabItem("文件重命名", fileRenamePage(w, &logVbox)),
		container.NewTabItem("图片重命名", imageXXXXRenamePage(w, &logVbox)),
		container.NewTabItem("多文件简易目录归位", simpleCatalogsHomingPage(w, &logVbox)),
		container.NewTabItem("图片归位", moveImagePage(w, &logVbox)),
		container.NewTabItem("移动目录", moveFilesInCatalogsPage(w, &logVbox)),
		container.NewTabItem("使用帮助&关于", showHelpPage()),
	))

	warning := "请在执行任何操作前均保证已仔细阅读其说明，确保您已知悉程序运行后果。\n如因操作不当造成数据损坏或丢失，使用者自行承担一切责任。确定要继续吗？"
	cnf := dialog.NewConfirm("警告", warning, func(r bool) {
		if !r {
			os.Exit(0)
		}
	}, w)
	cnf.SetDismissText("取消")
	cnf.SetConfirmText("确定")

	w.Resize(fyne.NewSize(750, 550))
	w.SetFixedSize(true) // 禁止改变窗口大小

	go func() {
		counter := 0
		for {
			logVbox.Refresh()
			t := <-logChan
			logTxt += time.Now().Format("2006-01-02 15:04:05 ") + t + "\r\n"
			if counter == 9 {
				logVbox = *container.NewVBox(widget.NewLabel(t))
				counter = 1
				continue
			}
			counter++
			logVbox.Add(widget.NewLabel(t))
		}

	}()

	go func() {
		time.Sleep(1 * time.Second)
		cnf.Show()
	}()

	w.ShowAndRun()
}
