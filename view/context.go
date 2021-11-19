package view

import (
	"todo/data"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type UiCtx struct {
	App        fyne.App
	MainWindow fyne.Window
}

func (c *UiCtx) SaveData(model *data.Model) {
	lbl := widget.NewLabel("Saving...")
	inf := widget.NewProgressBarInfinite()
	content := container.NewVBox(lbl, inf)
	popup := widget.NewModalPopUp(content, c.MainWindow.Canvas())
	popup.Show()
	defer popup.Hide()
	if err := data.SaveTaskData(model); err != nil {
		dialog.ShowError(err, c.MainWindow)
	}
}
