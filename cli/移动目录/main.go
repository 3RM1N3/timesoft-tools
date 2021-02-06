package main

import (
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	win := a.NewWindow("移动目录")
	ch := make(chan bool, 10)

	labelChooseDir := widget.NewLabel("Choose Dir:")
	entryChooseDir := widget.NewEntry()
	entryChooseDir.Resize(fyne.NewSize(100, 5))
	logEntry := widget.NewEntry()
	logEntry.MultiLine = true
	logEntry.DragEnd()
	runButton := widget.NewButton("Run", func() {
		ch <- false
		logLog("start move", logEntry)
		moveDir(entryChooseDir.Text, logEntry)
		logLog("success!", logEntry)
	})
	runButton.Disable()
	buttonChooseDir := widget.NewButton("...", func() {
		dialog.ShowFolderOpen(func(list fyne.ListableURI, err error) {
			if err != nil || list == nil {
				return
			}
			full := list.String()
			path := ""
			if strings.HasPrefix(full, "file://") {
				path = strings.TrimPrefix(list.String(), "file://")
			} else {
				path = strings.TrimPrefix(list.String(), "file:\\")
			}
			entryChooseDir.SetText(path)
			runButton.Enable()
		}, win)
	})

	win.SetContent(
		container.NewBorder(
			container.NewVBox(
				container.NewBorder(
					nil,
					nil,
					labelChooseDir,
					buttonChooseDir,
					entryChooseDir,
				),
				runButton,
			),
			nil,
			nil,
			nil,
			logEntry,
		),
	)
	go func(b *widget.Button, ch chan bool) {
		for {
			enable := <-ch
			if enable {
				b.Enable()
			} else {
				b.Disable()
			}
			time.Sleep(1 * time.Second)
		}
	}(runButton, ch)

	win.Resize(fyne.NewSize(600, 400))
	win.SetFixedSize(true)
	win.ShowAndRun()
}

func moveDir(dirPath string, logg *widget.Entry) {
	logLog("finding "+dirPath, logg)
	files, err := ioutil.ReadDir(dirPath)
	if err != nil {
		logTxt := "read dir error"
		logLog(logTxt, logg)
		return
	}
	for _, file := range files {
		fullPath := path.Join(dirPath, file.Name())
		fullPath = path.Clean(fullPath)
		if file.IsDir() && file.Name() == "目录" {
			filesInMulu, err := ioutil.ReadDir(fullPath)
			if err != nil {
				logLog("read catalog error", logg)
			}
			for _, file := range filesInMulu {
				if !file.IsDir() {
					aimPath := path.Dir(fullPath)
					fileName := file.Name()
					aimPath = path.Join(aimPath, fileName)
					err = os.Rename(path.Join(fullPath, fileName), aimPath)
					if err != nil {
						log.Println(err)
						logLog("move error", logg)
					}
				}
			}
		} else if file.IsDir() {
			moveDir(fullPath, logg)
		}
	}
}

func logLog(s string, logg *widget.Entry) {
	txt := logg.Text
	logg.SetText(s + "\n" + txt)
}
